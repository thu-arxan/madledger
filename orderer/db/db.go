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
	SetTxStatus(tx *core.Tx, status *TxStatus) error
	UpdateAccounts(accounts ...common.Account) error
	//SetAssetAdmin only succeed at the first time it is called
	SetAssetAdmin(pk crypto.PublicKey) error
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

	//GetAssetAdminPKBytes return nil is not exist
	GetAssetAdminPKBytes() []byte
	//GetOrCreateAccount return default account if not exist
	GetOrCreateAccount(address common.Address) (common.Account, error)
	//NewWriteBatch new a write batch
	NewWriteBatch() WriteBatch
	//GetTxStatus gets tx status
	GetTxStatus(channelID, txID string) (*TxStatus, error)
}
