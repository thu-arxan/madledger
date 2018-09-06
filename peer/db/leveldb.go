package db

import (
	"madledger/common"

	"github.com/syndtr/goleveldb/leveldb"
)

// LevelDB is the implementation of DB on leveldb
type LevelDB struct {
	// the dir of data
	dir     string
	connect *leveldb.DB
}

// NewLevelDB is the constructor of LevelDB
func NewLevelDB(dir string) (DB, error) {
	db := new(LevelDB)
	db.dir = dir
	connect, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, err
	}
	db.connect = connect
	return db, nil
}

// GetAccount returns an account of an address
func (db *LevelDB) GetAccount(address common.Address) (common.Account, error) {
	return nil, nil
}

// SetAccount updates an account or add an account
func (db *LevelDB) SetAccount(account common.Account) error {
	return nil
}

// RemoveAccount removes an account if exist
func (db *LevelDB) RemoveAccount(address common.Address) error {
	return nil
}

// GetStorage returns the key of an address if exist, else returns an error
func (db *LevelDB) GetStorage(address common.Address, key common.Word256) (common.Word256, error) {
	return common.ZeroWord256, nil
}

// SetStorage sets the value of a key belongs to an address
func (db *LevelDB) SetStorage(address common.Address, key common.Word256, value common.Word256) error {
	return nil
}
