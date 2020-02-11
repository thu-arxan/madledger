package raft

import (
	"errors"
	"madledger/consensus"
	"sync/atomic"
	"time"

	"madledger/common/util"
)

// Consensus is the implementation of interface
// Consensus keeps connections to raft services
type Consensus struct {
	cfg     *Config
	chain   *BlockChain // manage channels
	clients map[uint64]*Client
	ids     []uint64 // peer id
	leader  uint64   // leader id
}

// NewConseneus is the constructor of Consensus
func NewConseneus(cfg *Config) (*Consensus, error) {
	chain, err := NewBlockChain(cfg)
	if err != nil {
		return nil, err
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
	for id, addr := range c.cfg.ec.GetPeers() {
		client, err := NewClient(addr, c.cfg.cc.TLS)
		if err != nil {
			return err
		}
		c.clients[id] = client
		c.ids = append(c.ids, id)
	}
	// todo: is it necessary to select a leader randomly
	c.setLeader(c.ids[util.RandNum(len(c.ids))])

	if err := c.chain.Start(); err != nil {
		return err
	}

	return nil
}

// Stop is the implementation of interface
func (c *Consensus) Stop() error {
	c.chain.Stop()
	// todo
	// c.chain.raft.Stop()
	// c.chain.rpcServer.Stop()

	// close client conn
	for _, client := range c.clients {
		client.close()
	}
	return nil
}

// AddTx is the implementation of interface
func (c *Consensus) AddTx(channelID string, tx []byte) error {
	var err error
	// todo: modify log
	log.Infof("Consensus.AddTx: Add tx of channel %s.", channelID)
	// todo: we should parse the leader address other than random choose a leader
	// todo: modify the upper bounds of attempt times
	for i := 0; i < 100; i++ {
		leader := c.getLeader()
		log.Infof("Try to add tx to raft %d, this is %d times' trying.", leader, i)
		err = c.clients[leader].addTx(channelID, tx, c.cfg.id)
		if err == nil {
			log.Debugf("Node[%d] succeed to add tx to leader[%d]", c.cfg.id, c.leader)
			return nil
		}

		log.Debugf("add tx failed: %v", err)
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
			log.Infof("Unknown error: %v", err)
			c.setLeader(c.ids[util.RandNum(len(c.ids))])
			time.Sleep(100 * time.Millisecond)
		}
	}

	log.Infof("Add tx failed: %v", err)
	return err
}

// AddChannel is the implementation of interface
// Note: we can ignore this function now
func (c *Consensus) AddChannel(channelID string, cfg consensus.Config) error {
	log.Infof("Add channel %s", channelID)
	return errors.New("not implement")
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
	return leader
}
