package db

import (
	"madledger/common"
	"madledger/core/types"
)

// TxStatus return the status of tx
type TxStatus struct {
	Err             string
	BlockNumber     uint64
	BlockIndex      int
	Output          []byte
	ContractAddress string
}

// DB provide a interface for peer to access the global state
// Besides, it should also include the function that the evm StateDB provide
type DB interface {
	// GetAccount returns an account of an address
	GetAccount(address common.Address) (common.Account, error)
	// SetAccount updates an account or add an account
	SetAccount(account common.Account) error
	// RemoveAccount removes an account if exist
	RemoveAccount(address common.Address) error
	// GetStorage returns the key of an address if exist, else returns an error
	GetStorage(address common.Address, key common.Word256) (common.Word256, error)
	// SetStorage sets the value of a key belongs to an address
	SetStorage(address common.Address, key common.Word256, value common.Word256) error
	// However, the peer also should provide some functions to help the client to
	// know the result of the tx
	// GetStatus return the status of the tx
	GetTxStatus(channelID, txID string) (*TxStatus, error)
	GetTxStatusAsync(channelID, txID string) (*TxStatus, error)
	SetTxStatus(tx *types.Tx, status *TxStatus) error
	BelongChannel(channelID string) bool
	AddChannel(channelID string)
	DeleteChannel(channelID string)
	GetChannels() []string
	ListTxHistory(address []byte) map[string][]string
}
