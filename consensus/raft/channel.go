package raft

import (
	"fmt"
	"madledger/common/crypto"
	"madledger/common/event"
	"madledger/common/util"
	"madledger/consensus"
	"madledger/consensus/raft/eraft"
	"sync"
	"sync/atomic"
	"time"
)

// stolen from solo, to be modified
type channel struct {
	id     string
	lock   sync.Mutex
	raft   *eraft.Raft
	config consensus.Config
	txs    chan bool
	pool   *txPool
	hub    *event.Hub
	num    uint64
	// todo: gc to reduce the storage
	blocks map[uint64]*eraft.Block
	init   int32
	stop   chan bool
}

func newChannel(id string, config consensus.Config, raft *eraft.Raft) *channel {
	return &channel{
		id:     id,
		raft:   raft,
		config: config,
		num:    config.Number,
		txs:    make(chan bool, config.MaxSize),
		pool:   newTxPool(),
		hub:    event.NewHub(),
		blocks: make(map[uint64]*eraft.Block),
		init:   0,
		stop:   make(chan bool),
	}
}

func (c *channel) start() error {
	if c.initialized() {
		return fmt.Errorf("Consensus of channel %s is aleardy start", c.id)
	}
	c.setInit(1)
	ticker := time.NewTicker(time.Duration(c.config.Timeout) * time.Millisecond)
	// panic(c.config.Timeout)
	// log.Infof("Ticker duration is %d and block size is %d", c.config.Timeout, c.config.MaxSize)
	defer ticker.Stop()
	log.Infof("Channel %s start", c.id)
	for {
		select {
		case <-ticker.C:
			c.createBlock(c.pool.fetchTxs(c.config.MaxSize))
		case <-c.txs:
			if c.pool.getPoolSize() >= c.config.MaxSize {
				c.createBlock(c.pool.fetchTxs(c.config.MaxSize))
			}
		case <-c.stop:
			log.Infof("Stop channel %s consensus", c.id)
			c.setInit(0)
			return nil
		}
	}
}

func (c *channel) setInit(init int32) {
	atomic.StoreInt32(&c.init, init)
}

func (c *channel) initialized() bool {
	return atomic.LoadInt32(&c.init) != 0
}

// AddTx will try to add a tx
func (c *channel) AddTx(tx []byte) error {
	err := c.addTx(tx)
	if err != nil {
		return err
	}

	go func() {
		c.txs <- true
	}()

	result := c.hub.Watch(util.Hex(crypto.Hash(tx)), nil)
	return result.Err
}

func (c *channel) addTx(tx []byte) error {
	return c.pool.addTx(tx)
}

// Stop will block the work of channel
func (c *channel) Stop() {
	c.stop <- true
	for c.initialized() {
		time.Sleep(1 * time.Millisecond)
	}
}

func (c *channel) createBlock(txs [][]byte) error {
	if len(txs) == 0 {
		return nil
	}
	block := &eraft.Block{
		ChannelID: c.id,
		Num:       c.num,
		Txs:       txs,
	}

	// todo: call raft to add block
	if err := c.raft.AddBlock(block); err != nil {
		// todo: if we failed to create block we should release all txs
		log.Infof("[%d]Failed to add block %d because %v", c.raft.GetID(), block.Num, err)
		return err
	}

	// todo: handle config change tx

	c.blocks[block.Num] = block
	c.num++
	for _, tx := range block.Txs {
		hash := util.Hex(crypto.Hash(tx))
		c.hub.Done(hash, nil)
	}

	c.hub.Done(string(block.Num), nil)
	return nil
}

func (c *channel) getBlock(num uint64, async bool) (*eraft.Block, error) {
	c.lock.Lock()
	if util.Contain(c.blocks, num) {
		defer c.lock.Unlock()
		return c.blocks[num], nil
	}
	c.lock.Unlock()
	if async {
		c.hub.Watch(string(num), nil)
		c.lock.Lock()
		defer c.lock.Unlock()
		return c.blocks[num], nil
	}

	return nil, fmt.Errorf("Block %s:%d is not exist", c.id, c.num)
}
