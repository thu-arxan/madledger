// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package evm

import (
	"madledger/core"
	"madledger/peer/db"

	"github.com/thu-arxan/evm"
)

// Blockchain is the implementation of blockchain
type Blockchain struct {
	db        db.DB
	blocks    map[uint64]*core.Block
	channelID string
}

// NewBlockchain is the constructor of Blockchain
func NewBlockchain(engine db.DB, channelID string) *Blockchain {
	return &Blockchain{
		db:        engine,
		blocks:    make(map[uint64]*core.Block),
		channelID: channelID,
	}
}

// GetBlockHash is the implementation of interface
func (bc *Blockchain) GetBlockHash(num uint64) []byte {
	var hash = make([]byte, 32)
	var err error
	block := bc.blocks[num]
	if block == nil {
		block, err = bc.db.GetBlock(bc.channelID, num)
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
