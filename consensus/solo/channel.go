package solo

import (
	"fmt"
	"madledger/common/crypto"
	"madledger/common/event"
	"madledger/common/util"
	"madledger/consensus"
	"sync"
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
	init   bool
	stop   chan bool
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
		init:   false,
		stop:   make(chan bool),
	}
}

func (c *channel) start() error {
	if c.init {
		return fmt.Errorf("Consensus of channel %s is aleardy start", c.id)
	}
	c.init = true
	ticker := time.NewTicker(time.Duration(c.config.Timeout) * time.Millisecond)
	// panic(c.config.Timeout)
	log.Infof("Ticker duration is %d and block size is %d", c.config.Timeout, c.config.MaxSize)
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
			c.init = false
			return nil
		}
	}
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
	for c.init {
		time.Sleep(1 * time.Millisecond)
	}
}

func (c *channel) createBlock(txs [][]byte) error {
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
		hash := util.Hex(crypto.Hash(tx))
		c.hub.Done(hash, nil)
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
