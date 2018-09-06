package db

import "madledger/common"

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
}
