// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package db

import (
	"encoding/json"
	"errors"
	"fmt"
	cc "madledger/blockchain/config"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/event"
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
	hub     *event.Hub
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
	db.hub = event.NewHub()
	return db, nil
}

// HasChannel is the implementation of DB
func (db *LevelDB) HasChannel(id string) bool {
	exist, _ := db.connect.Has(getChannelProfileKey(id), nil)
	return exist
}

// UpdateChannel is the implementation of DB
func (db *LevelDB) UpdateChannel(id string, profile *cc.Profile) error {
	var key = getChannelProfileKey(id)
	if !db.HasChannel(id) {
		// 更新key为_config的记录，简单记录所有的test通道。 _config,  ["test11","test10","test21","test20"]
		err := db.addChannel(id)
		if err != nil {
			return err
		}
	}
	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	//更新key为_config@id的记录, 具体内容示例如下：
	// _config@test30 ,  {"Public":true,"Dependencies":null,"Members":[],
	// "Admins":[{"PK":"BN2PLBpBd5BrSLfTY7QEBYQT0h6lFvWlZyuAVt3/bfEz1g5QJ2lIEXP2Zk15B6E2MWpA/Q4Yxnl+XjFGObvAKTY=","Name":"admin"}]
	// "gasPrice": 1, "ratio": 1, "maxGas": 1000000 }
	err = db.connect.Put(key, data, nil)
	if err != nil {
		return err
	}
	return nil
}

func getChannelProfileKey(id string) []byte {
	return []byte(fmt.Sprintf("%s@%s", core.CONFIGCHANNELID, id))
}

// addChannel add a record into key core.CONFIGCHANNELID
func (db *LevelDB) addChannel(id string) error {
	var key = []byte(core.CONFIGCHANNELID)
	exist, _ := db.connect.Has(key, nil)
	var ids []string
	if exist {
		data, err := db.connect.Get(key, nil)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &ids)
		if err != nil {
			return err
		}
	}
	if !util.Contain(ids, id) {
		ids = append(ids, id)
	}
	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}
	err = db.connect.Put(key, data, nil)
	if err != nil {
		return err
	}
	return nil
}

// UpdateSystemAdmin update system admin
func (db *LevelDB) UpdateSystemAdmin(profile *cc.Profile) error {
	var key = getSystemAdminKey()
	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	//更新key为_config$admin的记录, 具体内容示例如下：
	//(_config$admin, {"Public":true,"Dependencies":null,"Members":null,"Admins":
	// [{"PK":"BGXcjZ3bhemsoLP4HgBwnQ5gsc8VM91b3y8bW0b6knkWu8xCSKO2qiJXARMHcbtZtvU7Jos2A5kFCD1haJ/hLdg=","Name":"SystemAdmin"}]})
	err = db.connect.Put(key, data, nil)
	if err != nil {
		return err
	}
	return nil
}

// IsSystemAdmin return if the member is the system admin
func (db *LevelDB) IsSystemAdmin(member *core.Member) bool {
	var p cc.Profile
	var key = getSystemAdminKey()
	data, err := db.connect.Get(key, nil)
	if err != nil {
		return false
	}
	err = json.Unmarshal(data, &p)
	if err != nil {
		return false
	}
	for i := range p.Admins {
		if p.Admins[i].Equal(member) {
			return true
		}
	}
	return false
}

func getSystemAdminKey() []byte {
	return []byte(fmt.Sprintf("%s$admin", core.CONFIGCHANNELID))
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
func (db *LevelDB) GetAccount(address common.Address) (*common.Account, error) {
	var key = util.BytesCombine([]byte("account:"), address.Bytes())
	value, err := db.connect.Get(key, nil)
	if err != nil {
		return common.NewAccount(address), nil
	}
	var account common.Account
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
		return nil, errors.New("not exist")
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
	status := db.hub.Watch(txID, func() { db.lock.Unlock() }).(*TxStatus)
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
	if channels == nil {
		channels = make([]string, 0)
	}
	return channels
}

// GetTxHistory is the implementation of interface
func (db *LevelDB) GetTxHistory(address []byte) map[string][]string {
	var txs = make(map[string][]string)
	if ok, _ := db.connect.Has(address, nil); ok {
		value, _ := db.connect.Get(address, nil)
		json.Unmarshal(value, &txs)
	}

	return txs
}

// NewWriteBatch implement the interface, WriteBatch is a wrapper of leveldb.Batch
func (db *LevelDB) NewWriteBatch() WriteBatch {
	batch := new(leveldb.Batch)
	return &WriteBatchWrapper{
		batch:     batch,
		db:        db,
		histories: make(map[string]map[string][]string),
	}
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

// Close close the leveldb
func (db *LevelDB) Close() {
	if db.connect != nil {
		db.connect.Close()
	}
}

// Get get the value by key
func (db *LevelDB) Get(key []byte, couldBeEmpty bool) ([]byte, error) {
	val, err := db.connect.Get(key, nil)
	if err == leveldb.ErrNotFound && couldBeEmpty {
		err = nil
	}
	return val, err
}

// GetAssetAdminPKBytes returns public key bytes of _asset admin or nil if not exists
func (db *LevelDB) GetAssetAdminPKBytes() []byte {
	var key = getAssetAdminKey()
	admin, err := db.connect.Get(key, nil)
	if err != nil {
		return nil
	}
	return admin
}

//GetOrCreateAccount return default account if not existx in leveldb
func (db *LevelDB) GetOrCreateAccount(address common.Address) (common.Account, error) {
	db.lock.Lock()
	account := common.NewAccount(address)
	key := getAccountKey(address)
	data, err := db.connect.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = nil
		}
		db.lock.Unlock()
		return *account, err
	}
	err = json.Unmarshal(data, &account)
	db.lock.Unlock()
	return *account, err
}

// WriteBatchWrapper is a wrapper of level.Batch
type WriteBatchWrapper struct {
	batch *leveldb.Batch
	db    *LevelDB

	histories map[string]map[string][]string
	channels  []string
}

// SetAccount is the implementation of interface
func (wb *WriteBatchWrapper) SetAccount(account *common.Account) error {
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
	if util.Contain(wb.histories, string(address)) {
		txs = wb.histories[string(address)]
		if !util.Contain(txs, channelID) {
			txs[channelID] = []string{txID}
		} else {
			txs[channelID] = append(txs[channelID], txID)
		}
	} else {
		if ok, _ := wb.db.connect.Has(address, nil); !ok {
			txs[channelID] = []string{txID}
		} else {
			value, err := wb.db.connect.Get(address, nil)
			if err == nil {
				json.Unmarshal(value, &txs)
				if !util.Contain(txs, channelID) {
					txs[channelID] = []string{txID}
				} else {
					txs[channelID] = append(txs[channelID], txID)
				}
			}
		}
		wb.histories[string(address)] = txs
	}
	value, _ := json.Marshal(txs)
	//db.connect.Put(address, value, nil)
	wb.batch.Put(address, value)
}

// Put stores (key, value) into batch, the caller is responsible to avoid duplicate key
func (wb *WriteBatchWrapper) Put(key, value []byte) {
	wb.batch.Put(key, value)
}

// PutBlock stores block into db
func (wb *WriteBatchWrapper) PutBlock(block *core.Block) error {
	data := block.Bytes()
	key := fmt.Sprintf("bc_data_%d", block.GetNumber())
	wb.batch.Put([]byte(key), data)
	return nil
}

// AddChannel is the implementation of interface
func (wb *WriteBatchWrapper) AddChannel(channelID string) {
	if wb.channels == nil {
		wb.channels = wb.db.GetChannels()
	}

	if !util.Contain(wb.channels, channelID) {
		wb.channels = append(wb.channels, channelID)
	}
	wb.updateChannels()
}

// DeleteChannel is the implementation of interface
func (wb *WriteBatchWrapper) DeleteChannel(channelID string) {
	if wb.channels == nil {
		wb.channels = wb.db.GetChannels()
	}
	var channels = make([]string, 0)
	for _, channel := range wb.channels {
		if channelID != channel {
			channels = append(channels, channel)
		}
	}
	wb.channels = channels
	wb.updateChannels()
}

// Sync sync batch to database
func (wb *WriteBatchWrapper) Sync() error {
	return wb.db.connect.Write(wb.batch, nil)
}

func (wb *WriteBatchWrapper) updateChannels() {
	var key = []byte("channels")
	value, _ := json.Marshal(wb.channels)
	wb.batch.Put(key, value)
}

//UpdateAccounts update asset
func (wb *WriteBatchWrapper) UpdateAccounts(accounts ...common.Account) error {
	for _, acc := range accounts {
		key := getAccountKey(acc.GetAddress())
		data, err := json.Marshal(acc)
		if err != nil {
			return err
		}
		wb.Put(key, data)
	}
	return nil
}

//SetAssetAdmin only succeed at the first time it is called
func (wb *WriteBatchWrapper) SetAssetAdmin(pk crypto.PublicKey) error {
	var key = getAssetAdminKey()
	exists, _ := wb.db.connect.Has(key, nil)
	if exists {
		return fmt.Errorf("account admin already set")
	}
	pkBytes, err := pk.Bytes()
	if err != nil {
		return err
	}
	wb.Put(key, pkBytes)
	return nil
}

func getAccountKey(address common.Address) []byte {
	return []byte(fmt.Sprintf("%s@%s", core.ASSETCHANNELID, address.String()))
}

func getAssetAdminKey() []byte {
	return []byte("_asset_admin")
}
