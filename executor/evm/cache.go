package evm

import (
	"fmt"
	"madledger/common"
	"madledger/common/util"
	"madledger/peer/db"
)

// Cache cache on a statedb.
// It will simulate operate on a db, and sync to db if necessary.
// Note: It's not thread safety because now it will only be used in one
// thread.
type Cache struct {
	db       StateDB
	accounts map[string]*accountInfo
}

type accountInfo struct {
	account common.Account
	storage map[string]common.Word256
	removed bool
	updated bool
}

// NewCache is the constructor of Cache
func NewCache(db StateDB) *Cache {
	return &Cache{
		db:       db,
		accounts: make(map[string]*accountInfo),
	}
}

// AccountExist return  if an account exist
func (cache *Cache) AccountExist(addr common.Address) bool {
	if util.Contain(cache.accounts, addressToString(addr)) {
		return true
	}
	return cache.db.AccountExist(addr)
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
	accInfo, err := cache.get(account.GetAddress())
	if err != nil {
		return err
	}
	if accInfo.removed {
		return fmt.Errorf("UpdateAccount on a removed account: %s", account.GetAddress())
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

	if util.Contain(accInfo.storage, word256ToString(key)) {
		return accInfo.storage[word256ToString(key)], nil
	}
	value, err := cache.db.GetStorage(address, key)
	if err != nil {
		return common.ZeroWord256, err
	}
	accInfo.storage[word256ToString(key)] = value
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
	accInfo.storage[word256ToString(key)] = value
	accInfo.updated = true
	return nil
}

// Sync sync changes to db
// If the sync return an error, it may cause something wrong, so it should be
// deal with by the developer.
// Also, this function may deal with the address and key in an order, so this
// function should be rethink if necessary.
// TODO: Sync should panic rather than return an error
func (cache *Cache) Sync(wb db.WriteBatch) error {
	var err error
	for address, account := range cache.accounts {
		if account.removed {
			if err = wb.RemoveAccount(stringToAddress(address)); err != nil {
				return err
			}
		} else if account.updated {
			err = wb.SetAccount(account.account)
			if err != nil {
				return err
			}
			for key, value := range account.storage {
				if err = wb.SetStorage(stringToAddress(address), stringToWord256(key), value); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// get the cache accountInfo item creating it if necessary
func (cache *Cache) get(address common.Address) (*accountInfo, error) {
	if util.Contain(cache.accounts, addressToString(address)) {
		return cache.accounts[addressToString(address)], nil
	}
	// Then try to load from db
	account, err := cache.db.GetAccount(address)
	if err != nil {
		// return nil, errors.New("The address is not exist")
		// should we return an error if there contains no account in the db?
		account = common.NewDefaultAccount(address)
	}
	// set the account
	cache.accounts[addressToString(address)] = &accountInfo{
		account: account,
		storage: make(map[string]common.Word256),
		removed: false,
		updated: false,
	}

	return cache.accounts[addressToString(address)], nil
}

func addressToString(address common.Address) string {
	return string(address.Bytes())
}

func stringToAddress(s string) common.Address {
	addr, _ := common.AddressFromBytes([]byte(s))
	return addr
}

func word256ToString(word common.Word256) string {
	return string(word.Bytes())
}

func stringToWord256(s string) common.Word256 {
	word, _ := common.BytesToWord256([]byte(s))
	return word
}
