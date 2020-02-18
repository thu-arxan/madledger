package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"madledger/common"
	"madledger/common/util"
	"madledger/core"
	"os"
	"sync"

	"github.com/tecbot/gorocksdb"
)

// RocksDB is the implementation of DB on rocksdb
type RocksDB struct {
	// the dir of data
	dir string

	lock sync.Mutex
	hub  *Hub

	connect      *gorocksdb.DB
	ro           *gorocksdb.ReadOptions
	wo           *gorocksdb.WriteOptions
	accountCFHdl *gorocksdb.ColumnFamilyHandle
	storageCFHdl *gorocksdb.ColumnFamilyHandle
	historyCFHdl *gorocksdb.ColumnFamilyHandle
}

// NewRocksDB is the constructor of RocksDB
func NewRocksDB(dir string) (DB, error) {
	db := new(RocksDB)
	db.dir = dir

	if !util.FileExists(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}
	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(gorocksdb.NewDefaultBlockBasedTableOptions())
	opts.SetCreateIfMissing(true)
	opts.SetCreateIfMissingColumnFamilies(true)

	var cfNames = []string{"default", "account", "storage", "history"}
	var cfOpts = make([]*gorocksdb.Options, len(cfNames))
	for i := range cfOpts {
		cfOpts[i] = opts
	}
	connect, cfHandles, err := gorocksdb.OpenDbColumnFamilies(opts, dir, cfNames, cfOpts)
	if err != nil {
		return nil, err
	}
	db.accountCFHdl = cfHandles[1]
	db.storageCFHdl = cfHandles[2]
	db.historyCFHdl = cfHandles[3]

	ro := gorocksdb.NewDefaultReadOptions()
	wo := gorocksdb.NewDefaultWriteOptions()
	wo.SetSync(false)

	db.connect = connect
	db.ro = ro
	db.wo = wo

	db.hub = NewHub()
	return db, nil
}

// AccountExist is the implementation of the interface
func (db *RocksDB) AccountExist(address common.Address) bool {
	var key = address.Bytes()
	data, err := db.connect.GetCF(db.ro, db.accountCFHdl, key)
	if err != nil {
		return false
	}
	defer data.Free()
	if data.Size() == 0 {
		return false
	}
	return true
}

// GetAccount returns an account of an address
func (db *RocksDB) GetAccount(address common.Address) (common.Account, error) {
	var key = address.Bytes()
	data, err := db.connect.GetCF(db.ro, db.accountCFHdl, key)
	if err != nil {
		return common.NewDefaultAccount(address), nil
	}
	defer data.Free()
	if data.Size() == 0 {
		return common.NewDefaultAccount(address), nil
	}
	var account common.DefaultAccount
	err = json.Unmarshal(data.Data(), &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
	// return UnmarshalAccount(value)
}

// GetStorage returns the key of an address if exist, else returns an error
func (db *RocksDB) GetStorage(address common.Address, key common.Word256) (common.Word256, error) {
	storageKey := util.BytesCombine(address.Bytes(), key.Bytes())
	data, err := db.connect.GetCF(db.ro, db.storageCFHdl, storageKey)
	if err != nil {
		return common.ZeroWord256, err
	}
	defer data.Free()
	if data.Size() == 0 {
		return common.ZeroWord256, errors.New("not found")
	}
	return common.BytesToWord256(data.Data())
}

// GetTxStatus is the implementation of interface
func (db *RocksDB) GetTxStatus(channelID, txID string) (*TxStatus, error) {
	var key = util.BytesCombine([]byte(channelID), []byte(txID))
	data, err := db.connect.Get(db.ro, key)
	if err != nil {
		return nil, err
	}
	defer data.Free()
	if data.Size() == 0 {
		return nil, errors.New("not exist")
	}
	var status TxStatus
	err = json.Unmarshal(data.Data(), &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// GetTxStatusAsync is the implementation of interface
func (db *RocksDB) GetTxStatusAsync(channelID, txID string) (*TxStatus, error) {
	db.lock.Lock()
	var key = util.BytesCombine([]byte(channelID), []byte(txID))
	// for {
	data, err := db.connect.Get(db.ro, key)
	db.lock.Unlock()
	if err == nil {
		defer data.Free()
		if data.Size() != 0 {
			var status TxStatus
			err = json.Unmarshal(data.Data(), &status)
			if err != nil {
				return nil, err
			}
			return &status, nil
		}
		status := db.hub.Watch(txID, func() {})
		return status, nil
	}
	return nil, err
}

// BelongChannel is the implementation of interface
func (db *RocksDB) BelongChannel(channelID string) bool {
	channels := db.GetChannels()
	if util.Contain(channels, channelID) {
		return true
	}
	return false
}

// AddChannel is the implementation of interface
func (db *RocksDB) AddChannel(channelID string) {
	channels := db.GetChannels()
	if !util.Contain(channels, channelID) {
		channels = append(channels, channelID)
	}
	db.setChannels(channels)
}

// DeleteChannel is the implementation of interface
func (db *RocksDB) DeleteChannel(channelID string) {
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
func (db *RocksDB) GetChannels() []string {
	var channels []string
	var key = []byte("channels")
	data, err := db.connect.Get(db.ro, key)
	if err != nil {
		return channels
	}
	defer data.Free()
	json.Unmarshal(data.Data(), &channels)
	return channels
}

// ListTxHistory is the implementation of interface
func (db *RocksDB) ListTxHistory(address []byte) map[string][]string {
	var result = make(map[string][]string)
	iter := db.connect.NewIteratorCF(db.ro, db.historyCFHdl)
	defer iter.Close()
	prefix := address
	iter.Seek(prefix)
	for ; iter.Valid() && iter.ValidForPrefix(prefix); iter.Next() {
		key := iter.Key().Data()
		if len(key) <= len(address) {
			log.Warnf("history key length is %d, which is short than address", len(key))
		} else {
			channelID := string(key[len(address):])
			var txs []string
			json.Unmarshal(iter.Value().Data(), &txs)
			if len(txs) != 0 {
				result[channelID] = txs
			}
			iter.Value().Free()
		}
		iter.Key().Free()
	}
	return result
}

func (db *RocksDB) setChannels(channels []string) {
	var key = []byte("channels")
	value, _ := json.Marshal(channels)
	db.connect.Put(db.wo, key, value)
}

// GetBlock gets block by block.num from db
func (db *RocksDB) GetBlock(num uint64) (*core.Block, error) {
	key := fmt.Sprintf("bc_data_%d", num)
	data, err := db.connect.Get(db.ro, []byte(key))
	if err != nil {
		return nil, err
	}
	defer data.Free()
	return core.UnmarshalBlock(data.Data())
}

// Close close the rocksdb
func (db *RocksDB) Close() {
	if db.connect != nil {
		db.connect.Close()
	}
}

// NewWriteBatch implement the interface, WriteBatch is a wrapper of gorocks.WriteBatch
func (db *RocksDB) NewWriteBatch() WriteBatch {
	batch := gorocksdb.NewWriteBatch()
	return &RocksDBWriteBatchWrapper{
		batch:     batch,
		db:        db,
		histories: make(map[string][]string),
	}
}

// RocksDBWriteBatchWrapper is a wrapper of gorocksdb.WriteBatchWrapper
type RocksDBWriteBatchWrapper struct {
	batch *gorocksdb.WriteBatch
	db    *RocksDB

	histories map[string][]string
}

// SetAccount is the implementation of interface
func (wb *RocksDBWriteBatchWrapper) SetAccount(account common.Account) error {
	var key = account.GetAddress().Bytes()
	value, err := account.Bytes()
	if err != nil {
		return err
	}
	// value := MarshalAccount(account)
	wb.batch.PutCF(wb.db.accountCFHdl, key, value)
	return nil

}

// RemoveAccount is the implementation of interface
func (wb *RocksDBWriteBatchWrapper) RemoveAccount(address common.Address) error {
	var key = address.Bytes()
	wb.batch.DeleteCF(wb.db.accountCFHdl, key)
	return nil
}

// RemoveAccountStorage delete all data associated with address
func (wb *RocksDBWriteBatchWrapper) RemoveAccountStorage(address common.Address) {
	// delete all associated data
	iter := wb.db.connect.NewIteratorCF(wb.db.ro, wb.db.storageCFHdl)
	defer iter.Close()
	prefix := address.Bytes()
	iter.Seek(prefix)
	for ; iter.Valid() && iter.ValidForPrefix(prefix); iter.Next() {
		key := iter.Key().Data()
		wb.batch.DeleteCF(wb.db.storageCFHdl, key)
		iter.Key().Free()
	}
}

// SetStorage is the implementation of interface
func (wb *RocksDBWriteBatchWrapper) SetStorage(address common.Address, key common.Word256, value common.Word256) error {
	storageKey := util.BytesCombine(address.Bytes(), key.Bytes())
	wb.batch.PutCF(wb.db.storageCFHdl, storageKey, value.Bytes())
	return nil
}

// SetTxStatus is the implementation of interface
func (wb *RocksDBWriteBatchWrapper) SetTxStatus(tx *core.Tx, status *TxStatus) error {
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
	wb.addHistory(sender.Bytes(), tx.Data.ChannelID, tx.ID)
	wb.db.hub.Done(tx.ID, status)
	return nil
}

// history key: address+channelID, value->[]{tx...}
func (wb *RocksDBWriteBatchWrapper) addHistory(address []byte, channelID, txID string) {
	var dbKey = util.BytesCombine(address, []byte(channelID))
	var historyKey = string(address)
	if util.Contain(wb.histories, historyKey) {
		wb.histories[historyKey] = append(wb.histories[historyKey], txID)
	} else {
		data, err := wb.db.connect.GetCF(wb.db.ro, wb.db.historyCFHdl, dbKey)
		if err != nil {
			wb.histories[historyKey] = []string{txID}
		}
		defer data.Free()
		if data.Size() == 0 {
			wb.histories[historyKey] = []string{txID}
		} else {
			var txs []string
			json.Unmarshal(data.Data(), &txs)
			txs = append(txs, txID)
			wb.histories[historyKey] = txs
		}
	}
	bytes, _ := json.Marshal(wb.histories[historyKey])
	wb.batch.PutCF(wb.db.historyCFHdl, dbKey, bytes)
}

// Put stores (key, value) into batch, the caller is responsible to avoid duplicate key
func (wb *RocksDBWriteBatchWrapper) Put(key, value []byte) {
	wb.batch.Put(key, value)
}

// PutBlock stores block into db
func (wb *RocksDBWriteBatchWrapper) PutBlock(block *core.Block) error {
	data := block.Bytes()
	key := fmt.Sprintf("bc_data_%d", block.GetNumber())
	wb.batch.Put([]byte(key), data)
	return nil
}

// Sync sync change to db
func (wb *RocksDBWriteBatchWrapper) Sync() error {
	return wb.db.connect.Write(wb.db.wo, wb.batch)
}
