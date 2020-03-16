// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package db

import (
	"madledger/common"
	"madledger/core"
)

// TxStatus return the status of tx
type TxStatus struct {
	Err             string
	BlockNumber     uint64
	BlockIndex      int
	Output          []byte
	ContractAddress string
}

// WriteBatch define a write batch interface
type WriteBatch interface {
	RemoveAccount(address common.Address) error
	SetAccount(account *common.Account) error
	SetStorage(address common.Address, key common.Word256, value common.Word256) error
	SetTxStatus(tx *core.Tx, status *TxStatus) error
	// PutBlock stores block into db
	PutBlock(block *core.Block) error
	// Put stores (key, value) into batch, the caller is responsible to avoid duplicate key
	Put(key, value []byte)
	RemoveAccountStorage(address common.Address)
	AddChannel(channelID string)
	DeleteChannel(channelID string)
	Sync() error
}

// DB provide a interface for peer to access the global state
// Besides, it should also include the function that the evm StateDB provide
type DB interface {
	AccountExist(address common.Address) bool
	// GetAccount returns an account of an address
	GetAccount(address common.Address) (*common.Account, error)
	// GetStorage returns the key of an address if exist, else returns an error
	GetStorage(address common.Address, key common.Word256) (common.Word256, error)
	// GetStatus return the status of the tx
	GetTxStatus(channelID, txID string) (*TxStatus, error)
	GetTxStatusAsync(channelID, txID string) (*TxStatus, error)
	BelongChannel(channelID string) bool
	GetChannels() []string
	GetTxHistory(address []byte) map[string][]string
	NewWriteBatch() WriteBatch
	// GetBlock gets block by block.num from db
	GetBlock(num uint64) (*core.Block, error)
	Close()
}
