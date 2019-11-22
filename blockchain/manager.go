package blockchain

import (
	"errors"
	"fmt"
	"madledger/core/types"
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
	except uint64
}

// NewManager is the constructor of manager
func NewManager(id, dir string) (*Manager, error) {
	except, err := load(dir)
	if err != nil {
		return nil, err
	}
	var m = Manager{
		lock:   new(sync.Mutex),
		id:     id,
		dir:    dir,
		except: except,
	}

	return &m, nil
}

// HasGenesisBlock return if the channel has a genesis block
func (manager *Manager) HasGenesisBlock() bool {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	return manager.except != 0
}

// GetBlock return the block of num
func (manager *Manager) GetBlock(num uint64) (*types.Block, error) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	if num >= manager.except {
		return nil, errors.New("The block is not exist")
	}
	return manager.loadBlock(num)
}

// GetExcept return the except
func (manager *Manager) GetExcept() uint64 {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	return manager.except
}

// GetPrevBlock return the prev block
func (manager *Manager) GetPrevBlock() *types.Block {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	if manager.except == 0 {
		return nil
	}
	block, err := manager.loadBlock(manager.except - 1)
	if err != nil {
		fmt.Println(err)
	}
	return block
}

// AddBlock add a block into the chain
func (manager *Manager) AddBlock(block *types.Block) error {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	if block.Header.Number != manager.except {
		return fmt.Errorf("Channel %s except block %d while receive block %d", manager.id, manager.except, block.Header.Number)
	}
	log.Infof("Channel %s add block %d", manager.id, manager.except)
	var err error

	err = manager.storeBlock(block)
	if err != nil {
		return err
	}
	err = manager.updateCache(block.Header.Number)
	if err != nil {
		return err
	}
	manager.except++
	return nil
}
