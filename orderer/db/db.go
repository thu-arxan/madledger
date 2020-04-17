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

// TxStatus return the status of tx
type TxStatus struct {
	Err             string
	BlockNumber     uint64
	BlockIndex      int
	Output          []byte
	ContractAddress string
}

// WriteBatch ...
type WriteBatch interface {
	// AddBlock will records all txs in the block to get rid of duplicated txs, cbn block
	AddBlock(block *core.Block) error
	// SetConsensusBlock set the consensus block of channel
	SetConsensusBlock(id string, num uint64)
	UpdateChannel(id string, profile *cc.Profile) error
	UpdateAccounts(accounts ...common.Account) error
	// SetAccount can only be called when atomicity is at one account level
	SetAccount(account common.Account) error
	UpdateSystemAdmin(profile *cc.Profile) error
	// SetAssetAdmin only succeed at the first time it is called
	// TODO: We need to change this function
	SetAssetAdmin(pk crypto.PublicKey) error
	Put(key, value []byte)
	Sync() error
}

// DB is the interface of db, and it is the implementation of DB on orderer/.tendermint/.glue
// TODO: We need reconsider all of these apis.
type DB interface {
	ListChannel() []string
	HasChannel(id string) bool
	GetChannelProfile(id string) (*cc.Profile, error)
	HasTx(tx *core.Tx) bool
	IsMember(channelID string, member *core.Member) bool
	IsAdmin(channelID string, member *core.Member) bool
	GetConsensusBlock(id string) uint64
	// WatchChannel provide a way to spy channel change. Now it mainly used to
	// spy channel create operation.
	WatchChannel(channelID string)
	Close() error
	IsSystemAdmin(member *core.Member) bool
	// if couldBeEmpty set to true and error is ErrNotFound
	// return no error
	Get(key []byte, couldBeEmpty bool) ([]byte, error)
	// GetAssetAdminPKBytes return nil is not exist
	GetAssetAdminPKBytes() []byte
	// GetOrCreateAccount return default account if not exist
	GetOrCreateAccount(address common.Address) (common.Account, error)
	// NewWriteBatch new a write batch
	NewWriteBatch() WriteBatch
}
