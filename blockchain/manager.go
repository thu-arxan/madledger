package blockchain

import (
	"errors"
	"fmt"
	"madledger/core"
)

// Manager manage the blockchain
// Warning: Not thread safety
type Manager struct {
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
		id:     id,
		dir:    dir,
		except: except,
	}

	return &m, nil
}

// HasGenesisBlock return if the channel has a genesis block
func (manager *Manager) HasGenesisBlock() bool {
	return manager.except != 0
}

// GetBlock return the block of num
func (manager *Manager) GetBlock(num uint64) (*core.Block, error) {
	if num >= manager.except {
		return nil, errors.New("The block is not exist")
	}
	return manager.loadBlock(num)
}

// AddBlock add a block into the chain
func (manager *Manager) AddBlock(block *core.Block) error {
	if block.Header.Number != manager.except {
		return fmt.Errorf("Except block %d while receive block %d", manager.except, block.Header.Number)
	}
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
