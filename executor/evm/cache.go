package evm

import (
	"fmt"
	"madledger/common"
	"madledger/util"
)

// Cache cache on a statedb.
// It will simulate operate on a db, and sync to db if necessary.
// Note: It's not thread safety because now it will only be used in one
// thread.
type Cache struct {
	db       StateDB
	accounts map[common.Address]*accountInfo
}

type accountInfo struct {
	account common.Account
	storage map[common.Word256]common.Word256
	removed bool
	updated bool
}

// NewCache is the constructor of Cache
func NewCache(db StateDB) *Cache {
	return &Cache{
		db:       db,
		accounts: make(map[common.Address]*accountInfo),
	}
}

// GetAccount return the account of address
func (cache *Cache) GetAccount(addr common.Address) (common.Account, error) {
	accountInfo, err := cache.get(addr)
	if err != nil {
		return nil, err
	}
	return accountInfo.account, nil
}

// SetAccount set account
func (cache *Cache) SetAccount(account common.Account) error {
	accInfo, err := cache.get(account.Address())
	if err != nil {
		return err
	}
	if accInfo.removed {
		return fmt.Errorf("UpdateAccount on a removed account: %s", account.Address())
	}
	accInfo.account = account
	accInfo.updated = true
	return nil
}

// RemoveAccount remove an account
func (cache *Cache) RemoveAccount(address common.Address) error {
	accInfo, err := cache.get(address)
	if err != nil {
		return err
	}
	if accInfo.removed {
		return fmt.Errorf("RemoveAccount on a removed account: %s", address)
	}
	accInfo.removed = true
	return nil
}

// GetStorage returns the key of an address if exist, else returns an error
func (cache *Cache) GetStorage(address common.Address, key common.Word256) (common.Word256, error) {
	// fmt.Printf("GetStorage of address %s and key %b\n", address.String(), key)
	accInfo, err := cache.get(address)
	if err != nil {
		return common.ZeroWord256, err
	}

	if util.Contain(accInfo.storage, key) {
		return accInfo.storage[key], nil
	}
	value, err := cache.db.GetStorage(address, key)
	if err != nil {
		return common.ZeroWord256, err
	}
	accInfo.storage[key] = value
	return value, nil
}

// SetStorage set the storage of address
// NOTE: Set value to zero to remove. How should i understand this?
func (cache *Cache) SetStorage(address common.Address, key common.Word256, value common.Word256) error {
	// fmt.Printf("!!!Set storage %s at key %b and value is %b\n", address.String(), key, value)
	accInfo, err := cache.get(address)
	if err != nil {
		return err
	}
	if accInfo.removed {
		return fmt.Errorf("SetStorage on a removed account: %s", address.String())
	}
	accInfo.storage[key] = value
	accInfo.updated = true
	return nil
}

// Sync sync changes to db
// If the sync return an error, it may cause something wrong, so it should be
// deal with by the developer.
// Also, this function may deal with the address and key in an order, so this
// function should be rethink if necessary.
// TODO: Sync should panic rather than return an error
func (cache *Cache) Sync() error {
	var err error
	for address, account := range cache.accounts {
		if account.removed {
			if err = cache.db.RemoveAccount(address); err != nil {
				return err
			}
		} else if account.updated {
			for key, value := range account.storage {
				if err = cache.db.SetStorage(address, key, value); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// get the cache accountInfo item creating it if necessary
func (cache *Cache) get(address common.Address) (*accountInfo, error) {
	if util.Contain(cache.accounts, address) {
		return cache.accounts[address], nil
	}
	// Then try to load from db
	account, err := cache.db.GetAccount(address)
	if err != nil {
		// return nil, errors.New("The address is not exist")
		// should we return an error if there contains no account in the db?
		account = common.NewDefaultAccount(address)
	}
	// set the account
	cache.accounts[address] = &accountInfo{
		account: account,
		storage: make(map[common.Word256]common.Word256),
		removed: false,
		updated: false,
	}

	return cache.accounts[address], nil
}
