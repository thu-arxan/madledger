package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"madledger/common"
	"madledger/common/util"

	"github.com/syndtr/goleveldb/leveldb"
)

// LevelDB is the implementation of State which using leveldb
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

// GetAccount is the implementation of DB
func (db *LevelDB) GetAccount(address common.Address) (common.Account, error) {
	var key = address.Bytes()
	data, err := db.connect.Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("The account which address is %s is not exist", address)
	}
	var account common.DefaultAccount
	err = json.Unmarshal(data, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// SetAccount is the implementation of DB
func (db *LevelDB) SetAccount(account common.Account) error {
	var key = account.GetAddress().Bytes()
	data, err := json.Marshal(account)
	if err != nil {
		return err
	}
	err = db.connect.Put(key, data, nil)
	if err != nil {
		return err
	}
	return nil
}

// RemoveAccount is the implementation of DB
func (db *LevelDB) RemoveAccount(address common.Address) error {
	var key = address.Bytes()
	if ok, err := db.connect.Has(key, nil); err == nil && ok {
		err = db.connect.Delete(key, nil)
		return err
	}
	return nil
}

// GetStorage is the implementation of DB
func (db *LevelDB) GetStorage(address common.Address, key common.Word256) (common.Word256, error) {
	var storageKey = util.BytesCombine(address.Bytes(), key.Bytes())
	value, err := db.connect.Get(storageKey, nil)
	if err != nil {
		return common.ZeroWord256, errors.New("key is not exist")
	}
	word, err := common.BytesToWord256(value)
	return word, err
}

// SetStorage is the implementation of DB
func (db *LevelDB) SetStorage(address common.Address, key common.Word256, value common.Word256) error {
	var storageKey = util.BytesCombine(address.Bytes(), key.Bytes())
	err := db.connect.Put(storageKey, value.Bytes(), nil)
	return err
}
