package evm

import "madledger/common"

// StateDB provide a interface for evm to access the global state
type StateDB interface {
	// AccountExist returns if an account exist
	AccountExist(address common.Address) bool
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
