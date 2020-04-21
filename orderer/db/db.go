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

// WriteBatch provide an interface to update datas atomic
type WriteBatch interface {
	// AddBlock will records all txs in the block to get rid of duplicated txs, cbn block
	AddBlock(block *core.Block) error
	// SetConsensusBlock set the consensus block of channel
	SetConsensusBlock(id string, num uint64)
	UpdateChannel(id string, profile *cc.Profile) error
	UpdateAccounts(accounts ...common.Account) error
	// SetAccount can only be called when atomicity is at one account level
	SetAccount(account common.Account) error
	// SetAssetAdmin only succeed at the first time it is called
	// TODO: We need to change this function
	SetAssetAdmin(pk crypto.PublicKey) error
	Put(key, value []byte)
	Sync() error
}

// DB is the interface of db, and it is the implementation of DB on orderer/.tendermint/.glue
// TODO: We need reconsider all of these apis.
type DB interface {
	// ListChannel list all channels that db knows, include system channel and user channel
	ListChannel() []string
	// HasChannel return if channel exist
	HasChannel(channelID string) bool
	// GetChannelProfile return profile of channel if exist, else return an error
	GetChannelProfile(channelID string) (*cc.Profile, error)
	// HasTx return if a tx exist
	HasTx(tx *core.Tx) bool
	// IsMember return if member belong to channel
	IsMember(channelID string, member *core.Member) bool
	// IsAdmin return if member if the admin of channel
	IsAdmin(channelID string, member *core.Member) bool
	// GetConsensusBlock return the last consensus block num that db knows
	GetConsensusBlock(channelID string) (num uint64)
	// WatchChannel provide a way to spy channel change. Now it mainly used to
	// spy channel create operation.
	WatchChannel(channelID string)
	// Close close db
	Close() error
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
