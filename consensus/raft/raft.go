package raft

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"madledger/consensus"
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

// AddTx is the implementation of interface
func (r *Raft) AddTx(channelID string, tx []byte) error {
	return errors.New("Not implementation yet")
}

// AddChannel is the implementation of interface
func (r *Raft) AddChannel(channelID string, cfg Config) error {
	return errors.New("Not implementation yet")
}

// GetBlock is the implementation of interface
func (r *Raft) GetBlock(channelID string, num uint64, async bool) (consensus.Block, error) {
	return nil, errors.New("Not implementation yet")
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

	// if !r.IsLeader() {
	// 	return fmt.Errorf("Leader is %s", r.GetLeader())
	// }

	// Note: Propose succeed means the block become an entry in the log, but it will not make sure that the block is the right block
	if err := r.eraft.propose(block.Bytes()); err != nil {
		return err
	}

	// fmt.Printf("[%d] Add block %d succeed\n", r.cfg.id, block.GetNumber())
	return r.app.watch(block)
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

func (r *Raft) setStatus(status int32) {
	atomic.StoreInt32(&r.status, status)
}

func (r *Raft) getStatus() int32 {
	return atomic.LoadInt32(&r.status)
}
