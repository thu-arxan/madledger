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

// GetTxStatus is the implementation of interface
func (db *LevelDB) GetTxStatus(channelID, txID string) (*TxStatus, error) {
	var key = util.BytesCombine([]byte(channelID), []byte(txID))
	value, err := db.connect.Get(key, nil)
	if err != nil {
		return nil, err
	}
	var status TxStatus
	err = json.Unmarshal(value, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// SetTxStatus is the implementation of interface
func (db *LevelDB) SetTxStatus(channelID, txID string, status *TxStatus) error {
	value, err := json.Marshal(status)
	if err != nil {
		return err
	}
	var key = util.BytesCombine([]byte(channelID), []byte(txID))
	err = db.connect.Put(key, value, nil)
	if err != nil {
		return err
	}
	return nil
}

// BelongChannel is the implementation of interface
func (db *LevelDB) BelongChannel(channelID string) bool {
	channels := db.GetChannels()
	if util.Contain(channels, channelID) {
		return true
	}
	return false
}

// AddChannel is the implementation of interface
func (db *LevelDB) AddChannel(channelID string) {
	channels := db.GetChannels()
	if !util.Contain(channels, channelID) {
		channels = append(channels, channelID)
	}
	db.setChannels(channels)
}

// DeleteChannel is the implementation of interface
func (db *LevelDB) DeleteChannel(channelID string) {
	oldChannels := db.GetChannels()
	var newChannels []string
	for i := range oldChannels {
		if channelID != oldChannels[i] {
			newChannels = append(newChannels, oldChannels[i])
		}
	}
	db.setChannels(newChannels)
}

// GetChannels is the implementation of interface
func (db *LevelDB) GetChannels() []string {
	var channels []string
	var key = []byte("channels")
	if ok, _ := db.connect.Has(key, nil); !ok {
		return channels
	}
	value, err := db.connect.Get(key, nil)
	if err != nil {
		return channels
	}
	json.Unmarshal(value, &channels)
	return channels
}

func (db *LevelDB) setChannels(channels []string) {
	var key = []byte("channels")
	value, _ := json.Marshal(channels)
	db.connect.Put(key, value, nil)
}
