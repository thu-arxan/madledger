package raft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"madledger/common/crypto"
	"madledger/common/event"
	"math"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"go.etcd.io/etcd/wal/walpb"

	"go.etcd.io/etcd/raft/raftpb"

	"go.etcd.io/etcd/etcdserver/api/snap"

	stats "go.etcd.io/etcd/etcdserver/api/v2stats"
	"go.etcd.io/etcd/pkg/fileutil"
	"go.etcd.io/etcd/pkg/types"

	"go.etcd.io/etcd/etcdserver/api/rafthttp"
	er "go.etcd.io/etcd/raft"
	"go.etcd.io/etcd/wal"
	"go.uber.org/zap"
)

// ERaft is a wrapper of etcd raft
type ERaft struct {
	cfg    *Config
	app    *App
	status int32
	hub    *event.Hub

	state       *state
	node        er.Node
	snapshotter *snap.Snapshotter
	storage     *er.MemoryStorage
	wal         *wal.WAL
	transport   *rafthttp.Transport
	httpServer  *http.Server

	lock       sync.Mutex
	stopCh     chan bool
	stopDoneCh chan bool
}

// state includes some necessary values
type state struct {
	lastIndex     uint64 // index of log at start
	snapshotIndex uint64
	appliedIndex  uint64
	conf          raftpb.ConfState
	leader        uint64
}

func (s *state) update(md raftpb.SnapshotMetadata) {
	s.conf = md.ConfState
	s.snapshotIndex = md.Index
	s.appliedIndex = md.Index
}

// NewERaft is the constructor of ERaft
func NewERaft(cfg *Config, app *App) (*ERaft, error) {
	return &ERaft{
		cfg:        cfg,
		app:        app,
		state:      new(state),
		hub:        event.NewHub(),
		status:     Stopped,
		stopCh:     make(chan bool, 1),
		stopDoneCh: make(chan bool, 1),
	}, nil
}

// Start start the eraft
func (e *ERaft) Start() error {
	var err error

	if atomic.LoadInt32(&e.status) != Stopped {
		return errors.New("The etcd raft is not stopped")
	}
	// remove the leader info
	atomic.StoreUint64(&(e.state.leader), 0)

	if err := initSnap(e.cfg.snapDir); err != nil {
		return err
	}
	e.snapshotter = snap.New(zap.NewExample(), e.cfg.snapDir)
	// if there exists an old wal
	oldWAL := wal.Exist(e.cfg.walDir)
	if !oldWAL {
		if err := initWAL(e.cfg.walDir); err != nil {
			return err
		}
	}
	// try to replay wal
	e.wal, err = e.replayWAL()
	if err != nil {
		return err
	}
	// set peers
	peers := make([]er.Peer, 0)
	for id := range e.cfg.peers {
		peers = append(peers, er.Peer{
			ID: id,
		})
	}
	// default eraft config
	erCfg := &er.Config{
		ID:                        e.cfg.id,
		ElectionTick:              10,
		HeartbeatTick:             1,
		Storage:                   e.storage,
		MaxSizePerMsg:             math.MaxUint64, // unlimited
		MaxInflightMsgs:           256,
		MaxUncommittedEntriesSize: 1 << 30,
		Logger:                    etcdLog,
	}
	// if old wal exist, then restart the node
	if oldWAL {
		e.node = er.RestartNode(erCfg)
	} else {
		e.node = er.StartNode(erCfg, peers)
	}

	e.transport = &rafthttp.Transport{
		Logger: zap.NewNop(),
		// Logger:      zap.NewExample(),
		ID:          types.ID(e.cfg.id),
		ClusterID:   0x10,
		Raft:        e,
		ServerStats: stats.NewServerStats("", ""),
		LeaderStats: stats.NewLeaderStats(strconv.FormatUint(e.cfg.id, 10)),
		ErrorC:      make(chan error),
	}

	if err := e.transport.Start(); err != nil {
		return err
	}
	for id := range e.cfg.peers {
		if id != e.cfg.id {
			e.transport.AddPeer(types.ID(id), []string{fmt.Sprintf("http://%s", e.cfg.peers[id])})
		}
	}

	// make sure two go routine can works right
	var wg sync.WaitGroup
	var wgErr error
	var wgLock sync.Mutex
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := e.serve(); err != nil {
			wgLock.Lock()
			defer wgLock.Unlock()
			wgErr = err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := e.startHTTP(); err != nil {
			wgLock.Lock()
			defer wgLock.Unlock()
			wgErr = err
		}
	}()

	wg.Wait()
	if wgErr != nil {
		return wgErr
	}

	// wait until the leader there is a leader
	for {
		if atomic.LoadUint64(&(e.state.leader)) != 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	atomic.StoreInt32(&e.status, Running)

	return nil
}

// Stop close etcd raft and something necessary
func (e *ERaft) Stop() {
	e.lock.Lock()
	defer e.lock.Unlock()

	if atomic.LoadInt32(&e.status) == Stopped {
		return
	}
	time.Sleep(100 * time.Millisecond)
	atomic.StoreUint64(&(e.state.leader), 0)

	e.transport.Stop()
	e.httpServer.Close()
	e.node.Stop()

	e.stopCh <- true
	<-e.stopDoneCh
	atomic.StoreInt32(&e.status, Stopped)
}

// serve  etcd raft
func (e *ERaft) serve() error {
	snap, err := e.storage.Snapshot()
	if err != nil {
		return err
	}
	e.state.update(snap.Metadata)

	go func() {
		defer e.wal.Close()

		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				e.node.Tick()
			case rd := <-e.node.Ready():
				if rd.SoftState != nil {
					atomic.StoreUint64(&(e.state.leader), rd.SoftState.Lead)
				}

				e.wal.Save(rd.HardState, rd.Entries)
				if !er.IsEmptySnap(rd.Snapshot) {
					e.saveSnapshot(rd.Snapshot)
					e.storage.ApplySnapshot(rd.Snapshot)
					e.publishSnapshot(rd.Snapshot)
				}
				e.storage.Append(rd.Entries)
				e.transport.Send(rd.Messages)
				if err := e.publishEntries(rd.CommittedEntries); err != nil {
					switch err.Error() {
					case "Removed from the cluster":
						// stop the eraft serive, but it will not gurantee the outter service
						if atomic.LoadInt32(&e.status) != Stopped {
							atomic.StoreInt32(&e.status, Stopped)
							atomic.StoreUint64(&(e.state.leader), 0)
							e.transport.Stop()
							e.httpServer.Close()
							e.node.Stop()
							// this is safe because there is no need for the machine restart again, and if the machine stop before, it must make sure the channel is empty
							e.stopDoneCh <- true
							return
						}
					default:
						panic(err)
					}
				}
				e.snapshot()
				e.node.Advance()
			case <-e.stopCh:
				e.stopDoneCh <- true
				return
			}
		}
	}()

	time.Sleep(50 * time.Millisecond)
	return nil
}

// startHTTP start http service
func (e *ERaft) startHTTP() error {
	l, err := net.Listen("tcp", e.cfg.getERaftAddress())
	if err != nil {
		return err
	}

	e.httpServer = &http.Server{
		Handler: e.transport.Handler(),
	}

	go func() {
		err = e.httpServer.Serve(l)
		if err != nil {
			if err.Error() == "http: Server closed" {
				return
			}
			panic(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	return nil
}

// Note: The etcd raft Propose will not gurantee the Propose succeed means that the data is appended to the log, so
// ERaft should make sure the thing happens by itself. And it need retry if timeout.
func (e *ERaft) propose(data []byte) error {
	hash := string(crypto.Hash(data))

	// try to propose
	if err := e.node.Propose(context.TODO(), data); err != nil {
		return err
	}

	var watchCh = make(chan bool, 1)
	go func() {
		e.hub.Watch(hash, nil)
		watchCh <- true
	}()

	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := e.node.Propose(context.TODO(), data); err != nil {
				return err
			}
		case <-watchCh:
			return nil
		}
	}

}

func (e *ERaft) proposeConfChange(cc raftpb.ConfChange) error {
	ccBytes, err := json.Marshal(cc)
	if err != nil {
		return err
	}

	hash := string(crypto.Hash(ccBytes))

	if err := e.node.ProposeConfChange(context.TODO(), cc); err != nil {
		return err
	}

	var watchCh = make(chan bool, 1)
	go func() {
		e.hub.Watch(hash, nil)
		watchCh <- true
	}()

	// todo: if conf change twice, it may cause panic, so we must conside a better way to do this
	ticker := time.NewTicker(5000 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := e.node.ProposeConfChange(context.TODO(), cc); err != nil {
				return err
			}
		case <-watchCh:
			return nil
		}
	}
}

func (e *ERaft) saveSnapshot(snap raftpb.Snapshot) error {
	walSnap := walpb.Snapshot{
		Index: snap.Metadata.Index,
		Term:  snap.Metadata.Term,
	}
	if err := e.wal.SaveSnapshot(walSnap); err != nil {
		return err
	}
	if err := e.snapshotter.SaveSnap(snap); err != nil {
		return err
	}
	return e.wal.ReleaseLockTo(snap.Metadata.Index)
}

func (e *ERaft) publishSnapshot(snap raftpb.Snapshot) error {
	if snap.Metadata.Index <= e.state.appliedIndex {
		return fmt.Errorf("Snapshot index [%d] should > applied index [%d]", snap.Metadata.Index, e.state.appliedIndex)
	}

	if err := e.app.UnMarshal(snap.Data); err != nil {
		return err
	}

	e.state.update(snap.Metadata)

	return nil
}

func (e *ERaft) publishEntries(ents []raftpb.Entry) error {
	if len(ents) == 0 {
		return nil
	}
	// find entried to apply
	var aes []raftpb.Entry
	firstIdx := ents[0].Index
	if firstIdx > e.state.appliedIndex+1 {
		return fmt.Errorf("First index of committed entry [%d] should <= [%d]+1", firstIdx, e.state.appliedIndex)
	}
	if e.state.appliedIndex-firstIdx+1 < uint64(len(ents)) {
		aes = ents[e.state.appliedIndex-firstIdx+1:]
	}

	for _, entry := range aes {
		switch entry.Type {
		case raftpb.EntryNormal:
			if len(entry.Data) == 0 {
				break
			}
			e.hub.Done(string(crypto.Hash(entry.Data)), nil)
			e.app.Commit(entry.Data)
		case raftpb.EntryConfChange:
			var cc raftpb.ConfChange
			cc.Unmarshal(entry.Data)
			e.state.conf = *e.node.ApplyConfChange(cc)
			switch cc.Type {
			case raftpb.ConfChangeAddNode:
				ccBytes, _ := json.Marshal(cc)
				e.hub.Done(string(crypto.Hash(ccBytes)), nil)
				if len(cc.Context) > 0 {
					// context should be the url of the new node etcd raft url
					e.transport.AddPeer(types.ID(cc.NodeID), []string{fmt.Sprintf("http://%s", string(cc.Context))})
				}
			case raftpb.ConfChangeRemoveNode:
				ccBytes, _ := json.Marshal(cc)
				e.hub.Done(string(crypto.Hash(ccBytes)), nil)
				if cc.NodeID == e.cfg.id {
					return errors.New("Removed from the cluster")
				}
				e.transport.RemovePeer(types.ID(cc.NodeID))
			}
		}

		e.state.appliedIndex = entry.Index

		if entry.Index == e.state.lastIndex {
			snapshot, err := e.loadSnapshot()
			if err == nil && snapshot != nil {
				e.app.UnMarshal(snapshot.Data)
			}
		}
	}

	return nil
}

func (e *ERaft) snapshot() error {
	if e.state.appliedIndex-e.state.snapshotIndex <= e.cfg.snapshotInterval {
		return nil
	}

	data, err := e.app.Marshal()
	if err != nil {
		return err
	}
	snap, err := e.storage.CreateSnapshot(e.state.appliedIndex, &e.state.conf, data)
	if err != nil {
		return err
	}
	if err := e.saveSnapshot(snap); err != nil {
		return err
	}

	compactIndex := uint64(1)
	if e.state.appliedIndex > e.cfg.snapshotInterval {
		compactIndex = e.state.appliedIndex - e.cfg.snapshotInterval
	}
	if err := e.storage.Compact(compactIndex); err != nil {
		return err
	}

	e.state.snapshotIndex = e.state.appliedIndex

	return nil
}

// replayWAL replays WAL entries into the raft instance
func (e *ERaft) replayWAL() (*wal.WAL, error) {
	// load snapshot
	snapshot, err := e.loadSnapshot()
	if err != nil {
		return nil, err
	}
	e.storage = er.NewMemoryStorage()
	// first apply snapshot
	if snapshot != nil {
		if err := e.storage.ApplySnapshot(*snapshot); err != nil {
			return nil, err
		}
	}
	// then open wal to load events after snapshot
	w, err := e.openWAL(snapshot)
	if err != nil {
		return nil, err
	}
	_, hs, ents, err := w.ReadAll()
	if err != nil {
		return nil, err
	}
	if err := e.storage.SetHardState(hs); err != nil {
		return nil, err
	}
	if err := e.storage.Append(ents); err != nil {
		return nil, err
	}
	if len(ents) != 0 {
		e.state.lastIndex = ents[len(ents)-1].Index
	} else {
		if snapshot != nil {
			e.app.UnMarshal(snapshot.Data)
		}
	}

	return w, nil
}

func (e *ERaft) openWAL(snapshot *raftpb.Snapshot) (*wal.WAL, error) {
	walSnap := walpb.Snapshot{}
	if snapshot != nil {
		walSnap.Index, walSnap.Term = snapshot.Metadata.Index, snapshot.Metadata.Term
	}
	log.Printf("loading WAL at term %d and index %d", walSnap.Term, walSnap.Index)
	return wal.Open(zap.NewExample(), e.cfg.walDir, walSnap)
}

func (e *ERaft) loadSnapshot() (*raftpb.Snapshot, error) {
	snapshot, err := e.snapshotter.Load()
	if err != nil && err != snap.ErrNoSnapshot {
		return nil, err
	}
	return snapshot, nil
}

// isLeader return if the eraft node is leader
func (e *ERaft) isLeader() bool {
	return e.getLeader() == e.cfg.id
}

// getLeader return the leader
func (e *ERaft) getLeader() uint64 {
	return atomic.LoadUint64(&(e.state.leader))
}

// Process is the implemetation of etcd raft
func (e *ERaft) Process(ctx context.Context, m raftpb.Message) error {
	return e.node.Step(ctx, m)
}

// IsIDRemoved is the implemetation of etcd raft
func (e *ERaft) IsIDRemoved(id uint64) bool { return false }

// ReportUnreachable is the implemetation of etcd raft
func (e *ERaft) ReportUnreachable(id uint64) {}

// ReportSnapshot is the implemetation of etcd raft
func (e *ERaft) ReportSnapshot(id uint64, status er.SnapshotStatus) {}

func initWAL(dir string) error {
	if !wal.Exist(dir) {
		if err := os.Mkdir(dir, 0750); err != nil {
			return err
		}

		w, err := wal.Create(zap.NewExample(), dir, nil)
		if err != nil {
			return err
		}
		w.Close()
	}
	return nil
}

func initSnap(dir string) error {
	if !fileutil.Exist(dir) {
		if err := os.Mkdir(dir, 0750); err != nil {
			return err
		}
	}
	return nil
}
