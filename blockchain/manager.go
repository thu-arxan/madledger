package blockchain

import (
	"errors"
	"fmt"
	"madledger/core"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "blockchain", "package": "blockchain"})
)

// Manager manage the blockchain
type Manager struct {
	lock   *sync.Mutex
	id     string
	dir    string
	expect uint64
}

// NewManager is the constructor of manager
func NewManager(id, dir string) (*Manager, error) {
	expect, err := load(dir)
	if err != nil {
		return nil, err
	}
	var m = Manager{
		lock:   new(sync.Mutex),
		id:     id,
		dir:    dir,
		expect: expect,
	}

	return &m, nil
}

// HasGenesisBlock return if the channel has a genesis block
func (manager *Manager) HasGenesisBlock() bool {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	return manager.expect != 0
}

// GetBlock return the block of num
func (manager *Manager) GetBlock(num uint64) (*core.Block, error) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	if num >= manager.expect {
		return nil, errors.New("The block does not exist")
	}
	return manager.loadBlock(num)
}

// GetExpect return the expect
func (manager *Manager) GetExpect() uint64 {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	return manager.expect
}

// GetPrevBlock return the prev block
func (manager *Manager) GetPrevBlock() *core.Block {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	if manager.expect == 0 {
		return nil
	}
	block, err := manager.loadBlock(manager.expect - 1)
	if err != nil {
		log.Warnf("channel %s manager failed to load block %d because of %v", manager.id, manager.expect-1, err)
	}
	return block
}

// AddBlock add a block into the chain
func (manager *Manager) AddBlock(block *core.Block) error {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	if block.Header.Number != manager.expect {
		return fmt.Errorf("Channel %s expect block %d while receive block %d", manager.id, manager.expect, block.Header.Number)
	}
	var err error
	defer func() {
		log.Infof("Channel %s add block %d, err: %v", manager.id, manager.expect, err)
	}()

	err = manager.storeBlock(block)
	if err != nil {
		return err
	}
	err = manager.updateCache(block.Header.Number)
	if err != nil {
		return err
	}
	manager.expect++
	return nil
}
