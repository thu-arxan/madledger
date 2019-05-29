package raft

import (
	"errors"
	"madledger/common/crypto"
	"madledger/common/event"
	"madledger/common/util"
	"madledger/consensus"
	"sync"
	"time"
)

// BlockChain will create blockchain in raft
type BlockChain struct {
	lock sync.Mutex

	raft    *Raft
	db      *DB
	config  consensus.Config
	txs     chan bool
	pool    *txPool
	hub     *event.Hub
	num     uint64
	blocks  map[uint64]*HybridBlock
	blockCh chan *HybridBlock

	quit chan bool
	done chan bool
}

// NewBlockChain is the constructor of blockchain
func NewBlockChain(cfg consensus.Config) (*BlockChain, error) {
	// todo:here should not be nil
	raft, err := NewRaft(nil)
	if err != nil {
		return nil, err
	}

	return &BlockChain{
		config:  cfg,
		txs:     make(chan bool, cfg.MaxSize),
		quit:    make(chan bool, 1),
		done:    make(chan bool, 1),
		blockCh: raft.BlockCh(),
	}, nil
}

func (chain *BlockChain) start() error {
	ticker := time.NewTicker(time.Duration(chain.config.Timeout) * time.Millisecond)
	defer ticker.Stop()

	log.Infof("Raft blockchain start")

	for {
		select {
		case <-ticker.C:
			chain.createBlock(chain.pool.fetchTxs(chain.config.MaxSize))
		case <-chain.txs:
			if chain.pool.getPoolSize() >= chain.config.MaxSize {
				chain.createBlock(chain.pool.fetchTxs(chain.config.MaxSize))
			}
		case <-chain.quit:
			chain.done <- true
			return nil
		}
	}
}

// AddTx will try to add a tx
func (chain *BlockChain) AddTx(tx []byte) error {
	if !chain.raft.IsLeader() {
		return errors.New("Please send to leader")
	}

	err := chain.addTx(tx)
	if err != nil {
		return err
	}

	go func() {
		chain.txs <- true
	}()

	result := chain.hub.Watch(util.Hex(crypto.Hash(tx)), nil)
	return result.Err
}

func (chain *BlockChain) addTx(tx []byte) error {
	return chain.pool.addTx(tx)
}

// Stop will block the work of channel
func (chain *BlockChain) Stop() {
	chain.quit <- true
	<-chain.done
}

func (chain *BlockChain) createBlock(txs [][]byte) error {
	if len(txs) == 0 {
		return nil
	}
	block := &HybridBlock{
		Num: chain.num,
		Txs: txs,
	}
	// then call eraft
	if err := chain.raft.AddBlock(block); err != nil {
		return err
	}
	chain.blocks[block.Num] = block
	chain.num++
	// for _, tx := range block.txs {
	// 	hash := util.Hex(crypto.Hash(tx))
	// 	c.hub.Done(hash, nil)
	// }

	// c.hub.Done(string(block.num), nil)
	return nil
}

func (chain *BlockChain) getBlock(num uint64, async bool) (*Block, error) {
	// c.lock.Lock()
	// if util.Contain(c.blocks, num) {
	// 	defer c.lock.Unlock()
	// 	return c.blocks[num], nil
	// }
	// c.lock.Unlock()
	// if async {
	// 	c.hub.Watch(string(num), nil)
	// 	c.lock.Lock()
	// 	defer c.lock.Unlock()
	// 	return c.blocks[num], nil
	// }

	// return nil, fmt.Errorf("Block %s:%d is not exist", c.id, c.num)
	return nil, errors.New("Not implementation")
}
