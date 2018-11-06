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
	// store all blocks, maybe gc is needed
	// todo: gc to reduce the storage
	blocks            map[uint64]*Block
	init              bool
	stop              chan bool
	consensuBlockChan *chan consensus.Block
}

func newChannel(id string, config consensus.Config) *channel {
	return &channel{
		id:                id,
		config:            config,
		num:               config.Number,
		txs:               make(chan bool, config.MaxSize),
		pool:              newTxPool(),
		hub:               event.NewHub(),
		blocks:            make(map[uint64]*Block),
		init:              false,
		stop:              make(chan bool),
		consensuBlockChan: nil,
	}
}

func (c *channel) start() error {
	if c.init {
		return fmt.Errorf("Consensus of channel %s is aleardy start", c.id)
	}
	c.init = true
	ticker := time.NewTicker(time.Duration(c.config.Timeout) * time.Millisecond)
	defer ticker.Stop()
	log.Infof("Channel %s start", c.id)
	for {
		select {
		case <-ticker.C:
			// log.Infof("Channel %s tick", c.id)
			c.createBlock(c.pool.fetchTxs(c.config.MaxSize))
		case <-c.txs:
			// see if there is a need to create block
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
	c.lock.Lock()
	err := c.addTx(tx)
	c.lock.Unlock()
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
	if c.consensuBlockChan != nil {
		go func(block *Block) {
			(*c.consensuBlockChan) <- block
		}(block)
	}
	return nil
}

func (c *channel) setConsensusBlockChan(ch *chan consensus.Block) error {
	c.consensuBlockChan = ch
	return nil
}
