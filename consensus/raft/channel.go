package raft

import (
	"encoding/json"
	"fmt"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/event"
	"madledger/common/util"
	"madledger/consensus"
	"madledger/consensus/raft/eraft"
	"madledger/core"

	"go.etcd.io/etcd/raft/raftpb"

	"sync"
	"sync/atomic"
	"time"
)

type channel struct {
	sync.Mutex

	id        uint64
	channelID string
	raft      *eraft.Raft
	config    consensus.Config
	txs       chan bool
	pool      *txPool
	hub       *event.Hub
	// todo: read it from db
	num     uint64 // block height
	blockCh chan *eraft.Block

	init int32
	stop chan bool
}

func newChannel(id uint64, channelID string, config consensus.Config, raft *eraft.Raft) *channel {
	return &channel{
		id:        id,
		channelID: channelID,
		raft:      raft,
		config:    config,
		num:       config.Number,
		txs:       make(chan bool, config.MaxSize),
		pool:      newTxPool(),
		hub:       event.NewHub(),
		init:      0,
		stop:      make(chan bool),
	}
}

func (c *channel) start() error {
	if c.initialized() {
		return fmt.Errorf("Consensus of channel %s is already start", c.channelID)
	}
	c.blockCh = c.raft.BlockCh(c.channelID)

	atomic.StoreUint64(&(c.num), c.raft.GetChainNum(c.channelID))
	log.Infof("Node[%d] start channel %s succeed, chainNum: %d", c.id, c.channelID, c.num)
	c.setInit(1)
	go func() {
		ticker := time.NewTicker(time.Duration(c.config.Timeout) * time.Millisecond)
		// panic(c.config.Timeout)
		log.Infof("Ticker duration is %d and block size is %d", c.config.Timeout, c.config.MaxSize)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.createBlock(c.pool.fetchTxs(c.config.MaxSize))
			case <-c.txs:
				if c.pool.getPoolSize() >= c.config.MaxSize {
					c.createBlock(c.pool.fetchTxs(c.config.MaxSize))
				}
			case block := <-c.blockCh:
				num := block.GetNumber()
				if num == c.num+1 {
					// todo: more work here
					log.Infof("channel %s fineshed block %d", c.channelID, num)
					c.blockDone(block)
				} else if num <= c.num {
					c.raft.FetchBlockDone(c.channelID, num)
				} else {
					c.raft.NotifyLater(block)
				}
			case <-c.stop:
				log.Infof("Stop channel %s consensus", c.channelID)
				c.setInit(0)
				return
			}
		}
	}()
	return nil
}

func (c *channel) setInit(init int32) {
	atomic.StoreInt32(&c.init, init)
}

func (c *channel) initialized() bool {
	return atomic.LoadInt32(&c.init) != 0
}

// AddTx will try to add a tx
func (c *channel) addTx(tx []byte) error {
	err := c.pool.addTx(tx)
	if err != nil {
		return err
	}

	go func() {
		c.txs <- true
	}()

	hash := util.Hex(crypto.Hash(tx))
	log.Infof("[%d][%s] watch tx: %s", c.id, c.channelID, hash)
	result := c.hub.Watch(hash, nil)

	return result.Err
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
		ChannelID: c.channelID,
		Num:       c.num + 1,
		Txs:       txs,
	}

	// todo: call raft to add block
	if err := c.raft.ProposeBlock(block); err != nil {
		// todo: if we failed to create block we should release all txs
		log.Infof("[%d]Failed to add block %d because %v", c.id, block.Num, err)
		return err
	}

	log.Infof("[%d] proposed channel(%s) addBlock[%d] succeed", c.id, c.channelID, block.Num)

	return c.blockDone(block)
}

// addBlock
func (c *channel) blockDone(block *eraft.Block) error {
	num := block.GetNumber()
	c.num = num
	c.raft.SetChainNum(c.channelID, num)
	c.raft.PutBlock(block)

	for _, tx := range block.Txs {
		hash := util.Hex(crypto.Hash(tx))
		log.Infof("Node[%d] channel[%s] hub done tx %s", c.raft.GetID(), c.channelID, hash)
		c.hub.Done(hash, nil)
	}

	// todo: why done here
	c.hub.Done(string(block.Num), nil)

	// todo: more to do
	// if c.channelID != "_global" {
	// 	return nil
	// }

	if !c.raft.IsLeader() {
		return nil
	}

	// todo: deal with conf change
	for _, tx := range block.Txs {
		if cfgChange := getConfChange(tx); cfgChange != nil {
			if err := c.raft.ProposeConfChange(cfgChange); err != nil {
				log.Errorf("conf change failed: %v", err)
				return err
			}
		}
	}
	return nil
}

func (c *channel) getBlock(num uint64, async bool) (*eraft.Block, error) {
	block := c.raft.GetBlock(c.channelID, num, async)
	if block != nil {
		return block, nil
	}
	return nil, fmt.Errorf("Block %s:%d is not exist", c.channelID, c.num)
}

func getConfChange(tx []byte) *raftpb.ConfChange {
	var coreTx core.Tx
	var cfgChange raftpb.ConfChange
	err := json.Unmarshal(tx, &coreTx)
	// Note: The reason return nil because Tx may be just random bytes,
	// and this is a bad implementation so we should change the way to do this
	// TODO: Reimplement it
	if err != nil {
		log.Errorf("invalid tx format: %v, tx: %s", err, string(tx))
		return nil
	}
	// get tx type according to recipient
	//log.Infof("Recipient: %s", common.BytesToAddress(coreTx.Data.Recipient).String())
	txType, err := core.GetTxType(common.BytesToAddress(coreTx.Data.Recipient).String())
	if err != nil {
		return nil
	}
	if txType != core.NODE {
		return nil
	}
	err = json.Unmarshal(coreTx.Data.Payload, &cfgChange)
	if err != nil {
		log.Errorf("failed to unmarshal cfgChange tx payload: %v, tx: %s", err, string(tx))
		return nil
	}
	return &cfgChange
}
