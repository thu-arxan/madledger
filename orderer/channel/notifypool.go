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
	lock     *sync.Mutex
	notifies map[string]*chan error
}

func newNotifyPool() *notifyPool {
	pool := new(notifyPool)
	pool.notifies = make(map[string]*chan error)
	pool.lock = new(sync.Mutex)
	return pool
}

func (pool *notifyPool) addNotify(hash string, e *chan error) error {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	if util.Contain(pool.notifies, hash) {
		return errors.New("The tx is aleardy contained")
	}
	pool.notifies[hash] = e
	return nil
}

func (pool *notifyPool) deleteNotify(hash string) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	if util.Contain(pool.notifies, hash) {
		delete(pool.notifies, hash)
	}
}

func (pool *notifyPool) addBlock(block *types.Block) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	txs := block.Transactions
	if len(txs) != 0 {
		for _, tx := range txs {
			hash := util.Hex(tx.Hash())
			if util.Contain(pool.notifies, hash) {
				e := pool.notifies[hash]
				delete(pool.notifies, hash)
				go func() {
					if e != nil {
						(*e) <- nil
					}
				}()
			}
		}
	}
}
