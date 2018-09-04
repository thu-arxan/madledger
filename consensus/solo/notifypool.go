package solo

import (
	"madledger/common/crypto"
	"madledger/common/util"
	"sync"
)

// it need to delete old notifies automatic
// so the time these notify be added should be recorded
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
	pool.notifies[hash] = e
	return nil
}

func (pool *notifyPool) addBlock(block *Block) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	txs := block.txs
	if len(txs) != 0 {
		for _, tx := range txs {
			hash := util.Hex(crypto.Hash(tx))
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
