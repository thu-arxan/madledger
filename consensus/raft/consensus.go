// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package raft

import (
	"madledger/consensus"
	"madledger/core"
	"sync/atomic"
	"time"

	"madledger/common/crypto"
	"madledger/common/util"
)

// Consensus is the implementation of interface
// Consensus keeps connections to raft services
type Consensus struct {
	cfg     *Config            // raft config
	chain   *BlockChain        // grpc service, manage channels
	clients map[uint64]*Client // grpc clients
	ids     []uint64           // orderer id
	leader  uint64             // leader id
}

// NewConsensus is the constructor of Consensus
func NewConsensus(channels map[string]consensus.Config, cfg *Config) (*Consensus, error) {
	chain, err := NewBlockChain(cfg)
	if err != nil {
		return nil, err
	}
	for id, cc := range channels {
		if err := chain.addChannels(id, cc); err != nil {
			return nil, err
		}
	}
	return &Consensus{
		cfg:   cfg,
		chain: chain,
	}, nil
}

// Start is the implementation of interface
func (c *Consensus) Start() error {
	c.clients = make(map[uint64]*Client)
	c.ids = make([]uint64, 0)
	// init grpc clients to connect with peers through grpc
	for id, addr := range c.cfg.peers {
		client, err := NewClient(addr, c.cfg.cc.TLS)
		if err != nil {
			return err
		}
		c.clients[id] = client
		c.ids = append(c.ids, id)
	}

	if err := c.chain.Start(); err != nil {
		log.Errorf("chain[%d] start failed: %v", c.cfg.id, err)
		return err
	}

	return nil
}

// Stop is the implementation of interface
func (c *Consensus) Stop() error {
	c.chain.Stop()

	// close client conn
	for _, client := range c.clients {
		client.close()
	}
	return nil
}

// AddTx is the implementation of interface
func (c *Consensus) AddTx(tx *core.Tx) error {
	var err error

	bytes, _ := tx.Bytes()
	hash := util.Hex(crypto.Hash(bytes))
	channelID := tx.Data.ChannelID

	// todo: we should parse the leader address other than random choose a leader
	// todo: modify the upper bounds of attempt times
	for i := 0; i < 10; i++ {
		leader := c.getLeader()
		log.Infof("Raft[%d] try %d times to add tx %s to node[%d]", c.cfg.id, i, hash, leader)
		err = c.clients[leader].addTx(channelID, bytes, c.cfg.id)
		if err == nil {
			log.Infof("Node[%d] succeed to add tx %s to leader[%d]", c.cfg.id, hash, c.leader)
			return nil
		}

		switch GetError(err) {
		case TxInPool:
			return err
		case RemovedNode:
			return err
		case NotLeader:
			id := GetLeader(err)
			if id == 0 {
				id = c.ids[util.RandNum(len(c.ids))]
			}
			c.setLeader(id)
		default:
			log.Infof("Node[%d] add tx %s failed: Unknown error: %v", c.cfg.id, tx.ID, err)
			c.setLeader(c.ids[util.RandNum(len(c.ids))])
		}
		time.Sleep(100 * time.Millisecond)
	}

	log.Infof("Node[%d] add tx %s failed: %v", c.cfg.id, hash, err)
	return err
}

// AddChannel is the implementation of interface
func (c *Consensus) AddChannel(channelID string, cfg consensus.Config) error {
	if err := c.chain.addChannels(channelID, cfg); err != nil {
		return err
	}
	return c.chain.startChannel(channelID)
}

// GetBlock is the implementation of interface
func (c *Consensus) GetBlock(channelID string, num uint64, async bool) (consensus.Block, error) {
	// log.Infof("Get block %d of channel %s", num, channelID)
	return c.chain.getBlock(channelID, num, async)
}

func (c *Consensus) setLeader(leader uint64) {
	atomic.StoreUint64(&c.leader, leader)
}

func (c *Consensus) getLeader() uint64 {
	leader := atomic.LoadUint64(&c.leader)

	if leader == 0 {
		leader = c.chain.raft.GetLeader()
		if leader == 0 {
			leader = c.ids[util.RandNum(len(c.ids))]
		}
		atomic.StoreUint64(&c.leader, leader)
	}
	return leader
}
