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
	SetAccount(account common.Account) error
	SetStorage(address common.Address, key common.Word256, value common.Word256) error
	SetTxStatus(tx *core.Tx, status *TxStatus) error
	// Put stores (key, value) into batch, the caller is responsible to avoid duplicate key
	Put(key, value []byte)
	RemoveAccountStorage(address common.Address)
	Sync() error
}

// DB provide a interface for peer to access the global state
// Besides, it should also include the function that the evm StateDB provide
type DB interface {
	AccountExist(address common.Address) bool
	// GetAccount returns an account of an address
	GetAccount(address common.Address) (common.Account, error)
	// SetAccount updates an account or add an account
	// TODO: This function is not necessary now?
	SetAccount(account common.Account) error
	// RemoveAccount removes an account if exist
	// TODO: This function is not necessary now?
	RemoveAccount(address common.Address) error
	// GetStorage returns the key of an address if exist, else returns an error
	GetStorage(address common.Address, key common.Word256) (common.Word256, error)
	// SetStorage sets the value of a key belongs to an address
	// TODO: This function is not necessary now?
	SetStorage(address common.Address, key common.Word256, value common.Word256) error
	// However, the peer also should provide some functions to help the client to
	// know the result of the tx
	// GetStatus return the status of the tx
	GetTxStatus(channelID, txID string) (*TxStatus, error)
	GetTxStatusAsync(channelID, txID string) (*TxStatus, error)
	// TODO: This function is not necessary now?
	SetTxStatus(tx *core.Tx, status *TxStatus) error
	BelongChannel(channelID string) bool
	AddChannel(channelID string)
	// TODO: This function should in WriteBatch?
	DeleteChannel(channelID string)
	GetChannels() []string
	ListTxHistory(address []byte) map[string][]string
	NewWriteBatch() WriteBatch

	// PutBlock stores block into db
	// TODO: Maybe write batch?
	PutBlock(block *core.Block) error
	// GetBlock gets block by block.num from db
	GetBlock(num uint64) (*core.Block, error)
}
