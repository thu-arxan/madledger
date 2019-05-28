package raft

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"go.etcd.io/etcd/raft/raftpb"

	"madledger/common/util"
	pb "madledger/consensus/raft/protos"
	core "madledger/core/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Raft wrap etcd raft and other things to provide a interface for using
type Raft struct {
	lock sync.Mutex

	cfg   *Config
	eraft *ERaft
	app   *App

	status    int32
	rpcServer *grpc.Server
}

// NewRaft is the constructor of Raft
func NewRaft(cfg *Config) (*Raft, error) {
	app, err := NewApp(cfg)
	if err != nil {
		return nil, err
	}
	eraft, err := NewERaft(cfg, app)
	if err != nil {
		return nil, err
	}
	return &Raft{
		app:    app,
		cfg:    cfg,
		eraft:  eraft,
		status: Stopped,
	}, nil
}

// Start start the raft service
// It will return after all service are started
// It should not use the lock because it will cause the service can not stop until the service start succeed
func (r *Raft) Start() error {
	if r.getStatus() != Stopped {
		return errors.New("The raft service failed to start because it status is not stopped")
	}
	r.setStatus(OnStarting)

	// if the raft join into an existing cluster, then call AddNode first
	if r.cfg.join {
		if err := r.joinCluster(); err != nil {
			return err
		}
	}

	if err := r.app.Start(); err != nil {
		return err
	}

	if err := r.eraft.Start(); err != nil {
		return err
	}

	if err := r.serve(); err != nil {
		return err
	}

	r.setStatus(Running)

	return nil
}

// Stop the raft service
func (r *Raft) Stop() {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.getStatus() == Stopped {
		return
	}

	// r.setStatus(Stopped)

	if r.rpcServer != nil {
		r.rpcServer.Stop()
	}

	if r.eraft != nil {
		r.eraft.Stop()
	}

	if r.app != nil {
		r.app.Stop()
	}

	r.setStatus(Stopped)
}

// AddBlock try to add a block, only leader is recommend to add block, however the etcd raft can not
// gurantee this, so we are trying our best to forbid follower or candidate adding block, but the fact is
// the first one trying to add block which num is except can succeed, but it maybe leader very possible.
// And this function should gurantee that return nil if the block is really added(works in the app)
func (r *Raft) AddBlock(block *core.Block) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.getStatus() != Running {
		return errors.New("The raft service is not running")
	}

	if !r.IsLeader() {
		return fmt.Errorf("Leader is %s", r.GetLeader())
	}

	// Note: Propose succeed means the block become an entry in the log, but it will not make sure that the block is the right block
	if err := r.eraft.propose(block.Bytes()); err != nil {
		return err
	}

	// fmt.Printf("[%d] Add block %d succeed\n", r.cfg.id, block.GetNumber())
	return r.app.watch(block)
}

// BlockCh provide a channel to fetch blocks
func (r *Raft) BlockCh() chan *core.Block {
	return r.app.blockCh
}

// IsLeader return if the node is leader of the cluster
func (r *Raft) IsLeader() bool {
	return r.eraft.isLeader()
}

// GetLeader return the leader's chain address
func (r *Raft) GetLeader() string {
	leader := r.eraft.getLeader()
	return r.cfg.getPeerAddress(leader)
}

// NotifyLater provide a mechanism for blockchain system to deal with the block which is too advanced
func (r *Raft) NotifyLater(block *core.Block) {
	r.app.notifyLater(block)
}

// FetchBlockDone is used for blockchain notify raft the block is stored on the disk
// and there is no need to store the blocks before to release the pressure of db
func (r *Raft) FetchBlockDone(num uint64) {
	r.app.fetchBlockDone(num)
}

// AddNode add a node into raft cluster
func (r *Raft) AddNode(ctx context.Context, info *pb.NodeInfo) (*pb.AddNodeResponse, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.getStatus() != Running {
		return &pb.AddNodeResponse{}, errors.New("The raft service is not running")
	}

	if !r.IsLeader() {
		return &pb.AddNodeResponse{}, fmt.Errorf("Leader is %s", r.GetLeader())
	}

	err := r.eraft.proposeConfChange(raftpb.ConfChange{
		ID:      util.RandUint64(),
		Type:    raftpb.ConfChangeAddNode,
		NodeID:  info.ID,
		Context: []byte(fmt.Sprintf("%s:%d", info.URL, info.ChainPort+2)),
	})
	if err != nil {
		return &pb.AddNodeResponse{}, err
	}
	return &pb.AddNodeResponse{}, nil
}

// RemoveNode remove a node into raft cluster
func (r *Raft) RemoveNode(ctx context.Context, info *pb.NodeInfo) (*pb.RemoveNodeResponse, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.getStatus() != Running {
		return &pb.RemoveNodeResponse{}, errors.New("The raft service is not running")
	}

	if !r.IsLeader() {
		return &pb.RemoveNodeResponse{}, fmt.Errorf("Leader is %s", r.GetLeader())
	}

	err := r.eraft.proposeConfChange(raftpb.ConfChange{
		ID:      util.RandUint64(),
		Type:    raftpb.ConfChangeRemoveNode,
		NodeID:  info.ID,
		Context: []byte(fmt.Sprintf("%s:%d", info.URL, info.ChainPort+2)),
	})
	if err != nil {
		return &pb.RemoveNodeResponse{}, err
	}

	return &pb.RemoveNodeResponse{}, nil
}

// serve will listen on address of config and provide service
func (r *Raft) serve() (err error) {
	lis, err := net.Listen("tcp", r.cfg.getLocalRaftAddress())
	if err != nil {
		return fmt.Errorf("Failed to start the raft service:%s", err)
	}

	var opts []grpc.ServerOption
	if r.cfg.tls.enable {
		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{r.cfg.tls.cert},
			RootCAs:      r.cfg.tls.pool,
			ClientCAs:    r.cfg.tls.pool,
		})
		opts = append(opts, grpc.Creds(creds))
	}
	r.rpcServer = grpc.NewServer(opts...)
	pb.RegisterRaftServer(r.rpcServer, r)

	go func() {
		err = r.rpcServer.Serve(lis)
		if err != nil {
			return
		}
	}()

	time.Sleep(100 * time.Millisecond)

	return nil
}

// joinCluster join into an existing cluster
func (r *Raft) joinCluster() error {
	var cluster = make(map[uint64]string)
	for id, addr := range r.cfg.peers {
		if id != r.cfg.id {
			addr := pb.ERaftToRaft(addr)
			if addr != "" {
				cluster[id] = addr
			}
		}
	}
	if len(cluster) == 0 {
		return errors.New("Cluster info miss")
	}
	// leader will be random choosed from the cluster
	var leader = randNode(cluster)
	var err error
	for i := 0; i < 10; i++ {
		var conn *grpc.ClientConn
		if r.cfg.tls.enable {
			creds := credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{r.cfg.tls.cert},
				RootCAs:      r.cfg.tls.pool,
			})
			conn, err = grpc.Dial(leader, grpc.WithTransportCredentials(creds), grpc.WithTimeout(2000*time.Millisecond))
		} else {
			conn, err = grpc.Dial(leader, grpc.WithInsecure(), grpc.WithTimeout(2000*time.Millisecond))
		}

		if err == nil {
			defer conn.Close()
			client := pb.NewRaftClient(conn)
			_, err = client.AddNode(context.Background(), &pb.NodeInfo{
				ID:        r.cfg.id,
				URL:       r.cfg.url,
				ChainPort: int32(r.cfg.chainPort),
			})
			// Node is added into cluster succeed if err is nil
			if err == nil {
				return nil
			}
			leaderChainAddress := getLeaderFromError(err)
			if leaderChainAddress != "" {
				leader = pb.ChainToRaft(leaderChainAddress)
				if leader != "" {
					continue
				}
			}
		}
		leader = randNode(cluster)
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("Failed to add into cluster: %s", err)
}

func (r *Raft) setStatus(status int32) {
	atomic.StoreInt32(&r.status, status)
}

func (r *Raft) getStatus() int32 {
	return atomic.LoadInt32(&r.status)
}
