// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package eraft

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"go.etcd.io/etcd/pkg/types"
	"go.etcd.io/etcd/raft/raftpb"
)

// Raft wrap etcd raft and other things to provide a interface for using
type Raft struct {
	lock sync.Mutex

	cfg   *EraftConfig
	eraft *ERaft // etcd raft, used for consensus
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

// ProposeBlock try to add a block, only leader is recommend to add block, however the etcd raft can not
// gurantee this, so we are trying our best to forbid follower or candidate adding block, but the fact is
// the first one trying to add block which num is expect can succeed, but it maybe leader very possible.
// And this function should gurantee that return nil if the block is really added(works in the app)
func (r *Raft) ProposeBlock(block *Block) error {
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
func (r *Raft) BlockCh(channelID string) chan *Block {
	return r.app.blockCh(channelID)
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
func (r *Raft) NotifyLater(block *Block) {
	r.app.notifyLater(block)
}

// FetchBlockDone is used for blockchain notify raft the block is stored on the disk
// and there is no need to store the blocks before to release the pressure of db
func (r *Raft) FetchBlockDone(channelID string, num uint64) {
	if r.IsLeader() {
		r.app.fetchBlockDone(channelID, num)
	}
}

func (r *Raft) setStatus(status int32) {
	atomic.StoreInt32(&r.status, status)
}

func (r *Raft) getStatus() int32 {
	return atomic.LoadInt32(&r.status)
}

// GetID ...
func (r *Raft) GetID() uint64 {
	return r.cfg.id
}

// SetChainNum set block height of channel
func (r *Raft) SetChainNum(channelID string, num uint64) {
	r.app.setChainNum(channelID, num)
}

// GetChainNum get block height of channel
func (r *Raft) GetChainNum(channelID string) uint64 {
	return r.app.getChainNum(channelID)
}

// GetBlock return the block of channel, return nil if not exist
func (r *Raft) GetBlock(channelID string, num uint64, async bool) *Block {
	return r.app.db.GetBlock(channelID, num, async)
}

// PutBlock stores block into db
func (r *Raft) PutBlock(block *Block) {
	r.app.db.AddBlock(block)
}

// ProposeConfChange propose a config change
func (r *Raft) ProposeConfChange(change *raftpb.ConfChange) error {
	// todo: check leader
	return r.eraft.proposeConfChange(*change)
}

// Removed ...
func (r *Raft) Removed(caller uint64) bool {
	return r.eraft.removed[types.ID(caller)]
}
