package raft

import (
	"errors"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/core"
	"sync"
)

// txPool store all txs which is not packed
type txPool struct {
	hashes map[string]bool
	txs    [][]byte
	lock   sync.Mutex
}

// newTxPool is the constructor of txPool
func newTxPool() *txPool {
	pool := new(txPool)
	pool.hashes = make(map[string]bool)
	return pool
}

// addTx add a transaction into pool
func (pool *txPool) addTx(tx []byte) error {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	// check if the tx is duplicated
	var hash = util.Hex(crypto.Hash(tx))
	if util.Contain(pool.hashes, hash) {
		return errors.New("Transaction is aleardy in the pool")
	}

	// add tx into the record
	pool.hashes[hash] = true
	pool.txs = append(pool.txs, tx)
	return nil
}

// getPoolSize return the tx size in pool
func (pool *txPool) getPoolSize() int {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	return len(pool.txs)
}

// however, we can not gc right away
// because the db is not updated yet
func (pool *txPool) fetchTxs(maxSize int) [][]byte {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	var size = len(pool.txs)
	if size > maxSize {
		size = maxSize
	}
	result := pool.txs[:size]
	newTx := pool.txs[size:]
	pool.txs = newTx
	return result
}

// todo: gc is not implementation yet
func (pool *txPool) gc(block *core.Block) error {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	return nil
}
