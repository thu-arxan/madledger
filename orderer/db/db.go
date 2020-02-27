package db

import (
	cc "madledger/blockchain/config"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/core"
)

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
	// todo:@zhq, i see these two functions want to support _asset?
	// But we can not know tx result from these two functions.
	// So you can change these two functions to GetTxStatus and SetTxStatus as peer db.
	// What's more, you should conside if your operatation if atomic? I think it is not atomic.
	IsTxExecute(txid string) bool
	SetTxExecute(txid string) error
}
