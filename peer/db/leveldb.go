package db

import (
	"encoding/json"
	"fmt"
	"madledger/common"
	"madledger/common/util"

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
	value, err := db.connect.Get(address.Bytes(), nil)
	if err != nil {
		return common.NewDefaultAccount(address), nil
	}
	var account common.DefaultAccount
	err = json.Unmarshal(value, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// SetAccount updates an account or add an account
func (db *LevelDB) SetAccount(account common.Account) error {
	fmt.Println("Set account:", account.GetAddress().String())
	value, err := account.Bytes()
	if err != nil {
		return err
	}
	err = db.connect.Put(account.GetAddress().Bytes(), value, nil)
	return err
}

// RemoveAccount removes an account if exist
func (db *LevelDB) RemoveAccount(address common.Address) error {
	if ok, _ := db.connect.Has(address.Bytes(), nil); !ok {
		return nil
	}

	return db.connect.Delete(address.Bytes(), nil)
}

// GetStorage returns the key of an address if exist, else returns an error
func (db *LevelDB) GetStorage(address common.Address, key common.Word256) (common.Word256, error) {
	// return common.ZeroWord256, nil
	storageKey := util.BytesCombine(address.Bytes(), key.Bytes())
	value, err := db.connect.Get(storageKey, nil)
	if err != nil {
		return common.ZeroWord256, err
	}
	return common.BytesToWord256(value)
}

// SetStorage sets the value of a key belongs to an address
func (db *LevelDB) SetStorage(address common.Address, key common.Word256, value common.Word256) error {
	storageKey := util.BytesCombine(address.Bytes(), key.Bytes())
	db.connect.Put(storageKey, value.Bytes(), nil)
	return nil
}
