package raft

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

// Raft wrap etcd raft and other things to provide a interface for using
type Raft struct {
	lock sync.Mutex

	cfg   *EraftConfig
	eraft *ERaft
	app   *App

	status int32
}

// NewRaft is the constructor of Raft
func NewRaft(cfg *EraftConfig) (*Raft, error) {
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

	log.Info("Raft Start function")
	if err := r.app.Start(); err != nil {
		return err
	}

	if err := r.eraft.Start(); err != nil {
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
func (r *Raft) AddBlock(block *HybridBlock) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.getStatus() != Running {
		return errors.New("The raft service is not running")
	}

	if !r.IsLeader() {
		return fmt.Errorf("Please send to leader %d", r.GetLeader())
	}

	// Note: Propose succeed means the block become an entry in the log, but it will not make sure that the block is the right block
	if err := r.eraft.propose(block.Bytes()); err != nil {
		return err
	}

	return r.app.watch(block)
}

// BlockCh provide a channel to fetch blocks
func (r *Raft) BlockCh() chan *HybridBlock {
	return r.app.blockCh
}

// IsLeader return if the node is leader of the cluster
func (r *Raft) IsLeader() bool {
	return r.eraft.isLeader()
}

// GetLeader return the leader's id
func (r *Raft) GetLeader() uint64 {
	return r.eraft.getLeader()
}

// NotifyLater provide a mechanism for blockchain system to deal with the block which is too advanced
func (r *Raft) NotifyLater(block *HybridBlock) {
	r.app.notifyLater(block)
}

// FetchBlockDone is used for blockchain notify raft the block is stored on the disk
// and there is no need to store the blocks before to release the pressure of db
func (r *Raft) FetchBlockDone(num uint64) {
	r.app.fetchBlockDone(num)
}

func (r *Raft) setStatus(status int32) {
	atomic.StoreInt32(&r.status, status)
}

func (r *Raft) getStatus() int32 {
	return atomic.LoadInt32(&r.status)
}
