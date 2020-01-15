package simulate

import (
	"fmt"
	"madledger/common"
	"madledger/common/util"
	"madledger/core"
	"madledger/peer/db"

	"github.com/syndtr/goleveldb/leveldb"
)

// StateDB is a memory db to simulate the actions of
// a normal StateDB
type StateDB struct {
	accounts map[common.Address]common.Account
	storages map[common.Address]map[common.Word256]common.Word256
}

// NewStateDB is the constructor of StateDB
func NewStateDB() *StateDB {
	return &StateDB{
		accounts: make(map[common.Address]common.Account),
		storages: make(map[common.Address]map[common.Word256]common.Word256),
	}
}

// AccountExist is the implementation of StateDB
func (s *StateDB) AccountExist(address common.Address) bool {
	if util.Contain(s.accounts, address) {
		return true
	}
	return false
}

// GetAccount is the implementaion of StateDB
func (s *StateDB) GetAccount(address common.Address) (common.Account, error) {
	if util.Contain(s.accounts, address) {
		return s.accounts[address], nil
	}
	return nil, fmt.Errorf("The address %s is not exist", address.String())
}

// GetStorage is the implementation of StateDB
// However, should it return no error if there doesn't return the address or key?
func (s *StateDB) GetStorage(address common.Address, key common.Word256) (common.Word256, error) {
	if util.Contain(s.storages, address) {
		storage := s.storages[address]
		if util.Contain(storage, key) {
			return storage[key], nil
		}
	}
	return common.ZeroWord256, nil
}

// SetAccount is the implementation of StateDB
func (s *StateDB) SetAccount(account common.Account) error {
	// if len(account.GetAddress()) == 0 {
	// 	return errors.New("The address of account can not be empty")
	// }
	s.accounts[account.GetAddress()] = account
	return nil
}

// SetStorage is the implementaion of StateDB
func (s *StateDB) SetStorage(address common.Address, key common.Word256, value common.Word256) error {
	if !util.Contain(s.storages, address) {
		s.storages[address] = make(map[common.Word256]common.Word256)
	}
	s.storages[address][key] = value
	return nil
}

// RemoveAccount removes an account if exist
func (s *StateDB) RemoveAccount(address common.Address) error {
	if !util.Contain(s.accounts, address) {
		return fmt.Errorf("The address %s is not exist", address.String())
	}
	delete(s.accounts, address)
	delete(s.storages, address)
	return nil
}

// NewWriteBatch ...
func (s *StateDB) NewWriteBatch() db.WriteBatch {
	return &WriteBatchWrapper{
		unknown: "madledger",
		stateDB: s,
	}
}

// WriteBatchWrapper ...
type WriteBatchWrapper struct {
	unknown string
	stateDB *StateDB
}

// SetAccount is the implementation of interface
func (wb *WriteBatchWrapper) SetAccount(account common.Account) error {
	wb.stateDB.accounts[account.GetAddress()] = account
	return nil

}

// RemoveAccount is the implementation of interface
func (wb *WriteBatchWrapper) RemoveAccount(address common.Address) error {
	if !util.Contain(wb.stateDB.accounts, address) {
		return fmt.Errorf("The address %s is not exist", address.String())
	}
	delete(wb.stateDB.accounts, address)
	delete(wb.stateDB.storages, address)
	return nil
}

// SetStorage is the implementation of interface
func (wb *WriteBatchWrapper) SetStorage(address common.Address, key common.Word256, value common.Word256) error {
	if !util.Contain(wb.stateDB.storages, address) {
		wb.stateDB.storages[address] = make(map[common.Word256]common.Word256)
	}
	wb.stateDB.storages[address][key] = value
	return nil
}

// SetTxStatus is the implementation of interface
func (wb *WriteBatchWrapper) SetTxStatus(tx *core.Tx, status *db.TxStatus) error {
	return nil
}

// GetBatch ...
func (wb *WriteBatchWrapper) GetBatch() *leveldb.Batch {
	return nil
}
