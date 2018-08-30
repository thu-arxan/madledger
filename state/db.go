package state

import "madledger/common"

// DB is the interface of DB
// It can be implementation by leveldb or any db
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
