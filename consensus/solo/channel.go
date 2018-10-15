package solo

import (
	"fmt"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/consensus"
	"time"
)

type channel struct {
	id       string
	config   consensus.Config
	txs      chan *txNotify
	pool     *txPool
	notifies *notifyPool
	num      uint64
	// store all blocks, maybe gc is needed
	// todo: gc to reduce the storage
	blocks            map[uint64]*Block
	notify            *chan consensus.Block
	init              bool
	stop              chan bool
	consensuBlockChan *chan consensus.Block
}

type txNotify struct {
	tx  []byte
	err *chan error
}

func newChannel(id string, config consensus.Config, notify *chan consensus.Block) *channel {
	return &channel{
		id:                id,
		config:            config,
		num:               config.Number,
		txs:               make(chan *txNotify, config.MaxSize),
		pool:              newTxPool(),
		notifies:          newNotifyPool(),
		notify:            notify,
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
		case notify := <-c.txs:
			err := c.addTx(notify.tx)
			if err != nil {
				go func() {
					(*notify.err) <- err
				}()
			} else {
				hash := util.Hex(crypto.Hash(notify.tx))
				c.notifies.addNotify(hash, notify.err)
			}
			// then see if there is a need to create block
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
	// log.Infof("Channel %s add tx", c.id)
	var e = make(chan error)
	notify := &txNotify{
		tx:  tx,
		err: &e,
	}
	go func() {
		c.txs <- notify
	}()

	err := <-e
	if err != nil {
		return err
	}
	// log.Infof("Channel %s succeed to add tx", c.id)
	return nil
}

func (c *channel) addTx(tx []byte) error {
	// try to add into the pool
	err := c.pool.addTx(tx)
	if err != nil {
		return err
	}
	return nil
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
	c.notifies.addBlock(block)
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
