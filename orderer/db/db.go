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
	cc "madledger/blockchain/config"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/core"
)

// TxStatus ...
type TxStatus struct {
	Executed bool
}

// WriteBatch ...
type WriteBatch interface {
	SetTxStatus(tx *core.Tx, status *TxStatus) error
	Sync() error
}

// DB is the interface of db, and it is the implementation of DB on orderer/.tendermint/.glue
type DB interface {
	ListChannel() []string
	HasChannel(id string) bool
	UpdateChannel(id string, profile *cc.Profile) error
	// AddBlock will records all txs in the block to get rid of duplicated txs
	AddBlock(block *core.Block) error
	HasTx(tx *core.Tx) bool
	IsMember(channelID string, member *core.Member) bool
	IsAdmin(channelID string, member *core.Member) bool
	// WatchChannel provide a way to spy channel change. Now it mainly used to
	// spy channel create operation.
	WatchChannel(channelID string)
	Close() error
	UpdateSystemAdmin(profile *cc.Profile) error
	IsSystemAdmin(member *core.Member) bool

	//IsAssetAdmin return true if pk is the public key of account channel admin
	IsAssetAdmin(pk crypto.PublicKey) bool
	//SetAssetAdmin only succeed at the first time it is called
	SetAssetAdmin(pk crypto.PublicKey) error
	//GetOrCreateAccount return default account if not exist
	GetOrCreateAccount(address common.Address) (common.Account, error)
	UpdateAccounts(accounts ...common.Account) error

	NewWriteBatch() WriteBatch
	GetTxStatus(channelID, txID string) (*TxStatus, error)
}
