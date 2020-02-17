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

	var cfNames = []string{"default", "account", "storage"}
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
// TODO: We may use column family to speed up
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

// SetAccount updates an account or add an account
func (db *RocksDB) SetAccount(account common.Account) error {
	return errors.New("no need to implement")
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
		return common.ZeroWord256, nil
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
		return nil, errors.New("Not exist")
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

// SetTxStatus is the implementation of interface
// TODO: Why should we need this?
func (db *RocksDB) SetTxStatus(tx *core.Tx, status *TxStatus) error {
	value, err := json.Marshal(status)
	if err != nil {
		return err
	}
	var key = util.BytesCombine([]byte(tx.Data.ChannelID), []byte(tx.ID))
	err = db.connect.Put(db.wo, key, value)
	if err != nil {
		return err
	}
	sender, err := tx.GetSender()
	if err != nil {
		return err
	}
	db.addHistory(sender.Bytes(), tx.Data.ChannelID, tx.ID)
	db.hub.Done(tx.ID, status)
	return nil
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
	var txs = make(map[string][]string)
	data, err := db.connect.Get(db.ro, address)
	if err != nil {
		return txs
	}
	defer data.Free()
	json.Unmarshal(data.Data(), &txs)
	return txs
}

func (db *RocksDB) addHistory(address []byte, channelID, txID string) {
	var txs = make(map[string][]string)
	data, err := db.connect.Get(db.ro, address)
	if err == nil {
		defer data.Free()
		if data.Size() == 0 {
			txs[channelID] = []string{txID}
		} else {
			json.Unmarshal(data.Data(), &txs)
			if !util.Contain(txs, channelID) {
				txs[channelID] = []string{txID}
			} else {
				if !util.Contain(txs[channelID], txID) {
					txs[channelID] = append(txs[channelID], txID)
				}
			}
		}
		value, _ := json.Marshal(txs)
		db.connect.Put(db.wo, address, value)
	}
}

func (db *RocksDB) setChannels(channels []string) {
	var key = []byte("channels")
	value, _ := json.Marshal(channels)
	db.connect.Put(db.wo, key, value)
}

// SyncWriteBatch sync write batch into db
func (db *RocksDB) SyncWriteBatch(batch *gorocksdb.WriteBatch) error {
	err := db.connect.Write(db.wo, batch)
	if err != nil {
		return err
	}
	return nil
}

// PutBlock stores block into db
func (db *RocksDB) PutBlock(block *core.Block) error {
	data := block.Bytes()
	key := fmt.Sprintf("bc_data_%d", block.GetNumber())
	return db.connect.Put(db.wo, []byte(key), data)
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

// NewWriteBatch implement the interface, WriteBatch is a wrapper of gorocks.WriteBatch
func (db *RocksDB) NewWriteBatch() WriteBatch {
	batch := gorocksdb.NewWriteBatch()
	return &RocksDBWriteBatchWrapper{
		batch: batch,
		db:    db,
	}
}

// RocksDBWriteBatchWrapper is a wrapper of gorocksdb.WriteBatchWrapper
type RocksDBWriteBatchWrapper struct {
	batch *gorocksdb.WriteBatch
	db    *RocksDB
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

func (wb *RocksDBWriteBatchWrapper) addHistory(address []byte, channelID, txID string) {
	var txs = make(map[string][]string)
	data, err := wb.db.connect.Get(wb.db.ro, address)
	if err == nil {
		defer data.Free()
		if data.Size() == 0 {
			txs[channelID] = []string{txID}
		} else {
			json.Unmarshal(data.Data(), &txs)
		}
		if !util.Contain(txs, channelID) {
			txs[channelID] = []string{txID}
		} else {
			if !util.Contain(txs[channelID], txID) {
				txs[channelID] = append(txs[channelID], txID)
			}
		}
		value, _ := json.Marshal(txs)
		wb.batch.Put(address, value)
	}
}

// Put stores (key, value) into batch, the caller is responsible to avoid duplicate key
func (wb *RocksDBWriteBatchWrapper) Put(key, value []byte) {
	wb.batch.Put(key, value)
}

// Sync sync change to db
func (wb *RocksDBWriteBatchWrapper) Sync() error {
	return wb.db.connect.Write(wb.db.wo, wb.batch)
}
