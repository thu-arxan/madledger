package evm

import (
	"github.com/thu-arxan/evm"
	"madledger/core"
	"madledger/peer/db"
)

// Blockchain is the implementation of blockchain
type Blockchain struct {
	db     db.DB
	blocks map[uint64]*core.Block
}

// NewBlockchain is the constructor of Blockchain
func NewBlockchain(engine db.DB) *Blockchain {
	return &Blockchain{
		db:     engine,
		blocks: make(map[uint64]*core.Block),
	}
}

// GetBlockHash is the implementation of interface
func (bc *Blockchain) GetBlockHash(num uint64) []byte {
	var hash = make([]byte, 32)
	var err error
	block := bc.blocks[num]
	if block == nil {
		block, err = bc.db.GetBlock(num)
		if err == nil {
			blockHash := block.Hash()
			hash = blockHash[:]
			bc.blocks[num] = block
		}
	}
	return hash
}

// CreateAddress is the implementation of interface
func (bc *Blockchain) CreateAddress(caller evm.Address, nonce uint64) evm.Address {
	return nil
}

// Create2Address is the implementation of interface
func (bc *Blockchain) Create2Address(caller evm.Address, salt, code []byte) evm.Address {
	return nil
}

// NewAccount is the implementation of interface
func (bc *Blockchain) NewAccount(address evm.Address) evm.Account {
	addr := address.(*Address)
	return NewAccount(addr)
}

// BytesToAddress is the implementation of interface
func (bc *Blockchain) BytesToAddress(bytes []byte) evm.Address {
	return BytesToAddress(bytes)
}
