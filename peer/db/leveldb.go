package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"madledger/common"
	"madledger/common/util"
	"madledger/core"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/syndtr/goleveldb/leveldb"
)

/*
* Here defines some key rules.
* 1. Account: key = []bytes("account:") + address.Bytes()
* 2. Storage: key = address.Bytes()
 */

// LevelDB is the implementation of DB on leveldb
type LevelDB struct {
	// the dir of data
	dir     string
	connect *leveldb.DB
	lock    sync.Mutex
	hub     *Hub
}

var (
	log = logrus.WithFields(logrus.Fields{"app": "peer", "package": "db"})
)

// NewLevelDB is the constructor of LevelDB
func NewLevelDB(dir string) (DB, error) {
	db := new(LevelDB)
	db.dir = dir
	connect, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, err
	}
	db.connect = connect
	db.hub = NewHub()
	return db, nil
}

// AccountExist is the implementation of the interface
func (db *LevelDB) AccountExist(address common.Address) bool {
	var key = util.BytesCombine([]byte("account:"), address.Bytes())
	_, err := db.connect.Get(key, nil)
	if err != nil {
		return false
	}
	return true
}

// GetAccount returns an account of an address
func (db *LevelDB) GetAccount(address common.Address) (common.Account, error) {
	var key = util.BytesCombine([]byte("account:"), address.Bytes())
	value, err := db.connect.Get(key, nil)
	if err != nil {
		return common.NewDefaultAccount(address), nil
	}
	var account common.DefaultAccount
	err = json.Unmarshal(value, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
	// return UnmarshalAccount(value)
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

// GetTxStatus is the implementation of interface
func (db *LevelDB) GetTxStatus(channelID, txID string) (*TxStatus, error) {
	var key = util.BytesCombine([]byte(channelID), []byte(txID))
	// TODO: Read twice is not necessary
	if ok, _ := db.connect.Has(key, nil); !ok {
		return nil, errors.New("Not exist")
	}
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

// GetTxStatusAsync is the implementation of interface
func (db *LevelDB) GetTxStatusAsync(channelID, txID string) (*TxStatus, error) {
	db.lock.Lock()
	var key = util.BytesCombine([]byte(channelID), []byte(txID))
	// for {
	if ok, _ := db.connect.Has(key, nil); ok {
		db.lock.Unlock()
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
	status := db.hub.Watch(txID, func() { db.lock.Unlock() })
	return status, nil
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

// ListTxHistory is the implementation of interface
func (db *LevelDB) ListTxHistory(address []byte) map[string][]string {
	var txs = make(map[string][]string)
	if ok, _ := db.connect.Has(address, nil); ok {
		value, _ := db.connect.Get(address, nil)
		json.Unmarshal(value, &txs)
	}

	return txs
}

func (db *LevelDB) addHistory(address []byte, channelID, txID string) {
	var txs = make(map[string][]string)
	if ok, _ := db.connect.Has(address, nil); !ok {
		txs[channelID] = []string{txID}
		value, _ := json.Marshal(txs)
		db.connect.Put(address, value, nil)
	} else {
		value, err := db.connect.Get(address, nil)
		if err == nil {
			json.Unmarshal(value, &txs)
			if !util.Contain(txs, channelID) {
				txs[channelID] = []string{txID}
			} else {
				txs[channelID] = append(txs[channelID], txID)
			}
			value, _ := json.Marshal(txs)
			db.connect.Put(address, value, nil)
		}
	}
}

func (db *LevelDB) setChannels(channels []string) {
	var key = []byte("channels")
	value, _ := json.Marshal(channels)
	db.connect.Put(key, value, nil)
}

// NewWriteBatch implement the interface, WriteBatch is a wrapper of leveldb.Batch
func (db *LevelDB) NewWriteBatch() WriteBatch {
	batch := new(leveldb.Batch)
	return &WriteBatchWrapper{
		batch: batch,
		db:    db,
	}
}

// SyncWriteBatch sync write batch into db
// func (db *LevelDB) SyncWriteBatch(batch *leveldb.Batch) error {
// 	err := db.connect.Write(batch, nil)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// PutBlock stores block into db
func (db *LevelDB) PutBlock(block *core.Block) error {
	data := block.Bytes()
	key := fmt.Sprintf("bc_data_%d", block.GetNumber())
	return db.connect.Put([]byte(key), data, nil)
}

// GetBlock gets block by block.num from db
func (db *LevelDB) GetBlock(num uint64) (*core.Block, error) {
	key := fmt.Sprintf("bc_data_%d", num)
	data, err := db.connect.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}
	return core.UnmarshalBlock(data)
}

// WriteBatchWrapper is a wrapper of level.Batch
type WriteBatchWrapper struct {
	batch *leveldb.Batch
	db    *LevelDB
}

// SetAccount is the implementation of interface
func (wb *WriteBatchWrapper) SetAccount(account common.Account) error {
	var key = util.BytesCombine([]byte("account:"), account.GetAddress().Bytes())
	value, err := account.Bytes()
	if err != nil {
		return err
	}
	// value := MarshalAccount(account)
	wb.batch.Put(key, value)
	return nil

}

// RemoveAccount is the implementation of interface
func (wb *WriteBatchWrapper) RemoveAccount(address common.Address) error {
	var key = util.BytesCombine([]byte("account:"), address.Bytes())
	wb.batch.Delete(key)
	return nil
}

// RemoveAccountStorage delete all data associated with address
func (wb *WriteBatchWrapper) RemoveAccountStorage(address common.Address) {
	// delete all associated data
	iter := wb.db.connect.NewIterator(nil, nil)
	defer iter.Release()
	addr := address.Bytes()
	iter.Seek(addr)
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		wb.batch.Delete(key)
	}
}

// SetStorage is the implementation of interface
func (wb *WriteBatchWrapper) SetStorage(address common.Address, key common.Word256, value common.Word256) error {
	storageKey := util.BytesCombine(address.Bytes(), key.Bytes())
	wb.batch.Put(storageKey, value.Bytes())
	return nil
}

// SetTxStatus is the implementation of interface
func (wb *WriteBatchWrapper) SetTxStatus(tx *core.Tx, status *TxStatus) error {
	value, err := json.Marshal(status)
	if err != nil {
		return err
	}
	var key = util.BytesCombine([]byte(tx.Data.ChannelID), []byte(tx.ID))
	wb.batch.Put(key, value)
	sender, err := tx.GetSender()
	if err != nil {
		return err
	}
	//db.addHistory(sender.Bytes(), tx.Data.ChannelID, tx.ID)
	wb.addHistory(sender.Bytes(), tx.Data.ChannelID, tx.ID)
	wb.db.hub.Done(tx.ID, status)
	return nil
}

func (wb *WriteBatchWrapper) addHistory(address []byte, channelID, txID string) {
	var txs = make(map[string][]string)
	if ok, _ := wb.db.connect.Has(address, nil); !ok {
		txs[channelID] = []string{txID}
		value, _ := json.Marshal(txs)
		//db.connect.Put(address, value, nil)
		wb.batch.Put(address, value)
	} else {
		value, err := wb.db.connect.Get(address, nil)
		if err == nil {
			json.Unmarshal(value, &txs)
			if !util.Contain(txs, channelID) {
				txs[channelID] = []string{txID}
			} else {
				txs[channelID] = append(txs[channelID], txID)
			}
			value, _ := json.Marshal(txs)
			//db.connect.Put(address, value, nil)
			wb.batch.Put(address, value)
		}
	}
}

// Put stores (key, value) into batch, the caller is responsible to avoid duplicate key
func (wb *WriteBatchWrapper) Put(key, value []byte) {
	wb.batch.Put(key, value)
}

// Sync sync batch to database
func (wb *WriteBatchWrapper) Sync() error {
	return wb.db.connect.Write(wb.batch, nil)
}
