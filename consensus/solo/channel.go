package solo

import (
	"fmt"
	"madledger/common/event"
	"madledger/common/util"
	"madledger/consensus"
	"madledger/core"
	"sync"
	"sync/atomic"
	"time"
)

type channel struct {
	id     string
	lock   sync.Mutex
	config consensus.Config
	txs    chan bool
	pool   *txPool
	hub    *event.Hub
	num    uint64
	// todo: gc to reduce the storage
	blocks map[uint64]*Block
	init   int32
	stop   chan *chan bool
}

func newChannel(id string, config consensus.Config) *channel {
	return &channel{
		id:     id,
		config: config,
		num:    config.Number,
		txs:    make(chan bool, config.MaxSize),
		pool:   newTxPool(),
		hub:    event.NewHub(),
		blocks: make(map[uint64]*Block),
		init:   0,
		stop:   make(chan *chan bool),
	}
}

func (c *channel) start() error {
	if c.initialized() {
		return fmt.Errorf("Consensus of channel %s is already start", c.id)
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
		case ch := <-c.stop:
			log.Infof("Stop channel %s consensus", c.id)
			c.setInit(0)
			*ch <- true
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
func (c *channel) AddTx(tx *core.Tx) error {
	err := c.addTx(tx)
	if err != nil {
		return err
	}

	go func() {
		c.txs <- true
	}()

	result := c.hub.Watch(tx.ID, nil)
	return result.Err
}

func (c *channel) addTx(tx *core.Tx) error {
	return c.pool.addTx(tx)
}

// Stop will block the work of channel
func (c *channel) Stop() {
	stopDone := make(chan bool, 1)
	c.stop <- &stopDone
	<-stopDone
}

func (c *channel) createBlock(txs []*core.Tx) error {
	if len(txs) == 0 {
		return nil
	}
	block := &Block{
		channelID: c.id,
		num:       c.num,
		txs:       txs,
	}
	c.blocks[block.num] = block
	c.num++
	for _, tx := range block.txs {
		c.hub.Done(tx.ID, nil)
	}

	c.hub.Done(string(block.num), nil)
	return nil
}

func (c *channel) getBlock(num uint64, async bool) (*Block, error) {
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
