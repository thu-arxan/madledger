package channel

import (
	"errors"
	"madledger/common/util"
	"madledger/core/types"
	"sync"
)

// notifyPool provide a way to use channel to notify events
type notifyPool struct {
	// use the lock to protect the map
	lock *sync.Mutex
	// txs include all tx events
	txs map[string]*chan error
	// blocks include all block events
	blocks map[uint64][]*chan bool
}

func newNotifyPool() *notifyPool {
	pool := new(notifyPool)
	pool.txs = make(map[string]*chan error)
	pool.blocks = make(map[uint64][]*chan bool)
	pool.lock = new(sync.Mutex)
	return pool
}

func (pool *notifyPool) addTxNotify(hash string, e *chan error) error {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	if util.Contain(pool.txs, hash) {
		return errors.New("The tx is aleardy contained")
	}
	pool.txs[hash] = e
	return nil
}

func (pool *notifyPool) deleteTxNotify(hash string) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	if util.Contain(pool.txs, hash) {
		delete(pool.txs, hash)
	}
}

func (pool *notifyPool) addBlockNotify(num uint64, b *chan bool) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	if !util.Contain(pool.blocks, num) {
		pool.blocks[num] = make([]*chan bool, 0)
	}
	pool.blocks[num] = append(pool.blocks[num], b)
}

func (pool *notifyPool) addBlock(block *types.Block) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	num := block.Header.Number
	if util.Contain(pool.blocks, num) {
		for i := range pool.blocks[num] {
			b := pool.blocks[num][i]
			go func() {
				if b != nil {
					(*b) <- true
				}
			}()
		}
	}

	txs := block.Transactions
	if len(txs) != 0 {
		for _, tx := range txs {
			hash := util.Hex(tx.Hash())
			// fmt.Printf("Notify notify %s\n", hash)
			if util.Contain(pool.txs, hash) {
				e := pool.txs[hash]
				delete(pool.txs, hash)
				go func() {
					if e != nil {
						(*e) <- nil
					}
				}()
			}
		}
	}
}
