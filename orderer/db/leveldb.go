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
	"madledger/common"

	cc "madledger/blockchain/config"
	"madledger/common/crypto"
	"madledger/common/event"
	"madledger/common/util"
	"madledger/core"

	"github.com/syndtr/goleveldb/leveldb"
)

/*
* Here will describe how the db store key/value.
* 1. All channels: []byte("channelList") -> []string{}
* 2. Channel Profile: []byte(channelID+"@profile") -> Profile
* 3. Account: []byte(account.Address()) -> Account
* 4. Consensus Block: []byte("cbn"+channelID) -> uint64
* 5. Tx: []byte(tx.ChannelID+tx.ID) -> []byte{1}
 */

// LevelDB is the implementation of DB on orderer/data/leveldb
type LevelDB struct {
	// the dir of data
	dir     string
	connect *leveldb.DB
	hub     *event.Hub
}

// NewLevelDB is the constructor of LevelDB
func NewLevelDB(dir string) (DB, error) {
	var err error

	db := new(LevelDB)
	db.dir = dir
	db.connect, err = leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, err
	}
	db.hub = event.NewHub()

	return db, nil
}

// ListChannel is the implementation of DB
func (db *LevelDB) ListChannel() []string {
	var key = getChannelListKey()
	var channels []string
	data, err := db.connect.Get(key, nil)
	if err != nil {
		return channels
	}
	json.Unmarshal(data, &channels)
	return channels
}

// HasChannel return if channel exist
func (db *LevelDB) HasChannel(channelID string) bool {
	if channelID == "" {
		return false
	}
	exist, _ := db.connect.Has(getChannelProfileKey(channelID), nil)
	return exist
}

// GetChannelProfile return profile of channel
func (db *LevelDB) GetChannelProfile(channelID string) (*cc.Profile, error) {
	if channelID == "" {
		return nil, errors.New("channel id should not be empty")
	}
	var key = getChannelProfileKey(channelID)
	data, err := db.connect.Get(key, nil)
	if err != nil {
		return nil, err
	}
	var profile cc.Profile
	err = json.Unmarshal(data, &profile)
	return &profile, err
}

// GetConsensusBlock return the consensus block number of channel
func (db *LevelDB) GetConsensusBlock(channelID string) uint64 {
	key := getConsensusBlockKey(channelID)
	data, err := db.connect.Get(key, nil)
	if err != nil || len(data) == 0 {
		return 1
	}
	num, err := util.BytesToUint64(data)
	if err != nil || num < 1 {
		return 1
	}
	return num
}

// HasTx return if the tx is contained
func (db *LevelDB) HasTx(tx *core.Tx) bool {
	key := util.BytesCombine([]byte(tx.Data.ChannelID), []byte(tx.ID))

	if exist, _ := db.connect.Has(key, nil); exist {
		return true
	}
	return false
}

// IsMember is the implementation of DB
func (db *LevelDB) IsMember(channelID string, member *core.Member) bool {
	var p cc.Profile
	var key = getChannelProfileKey(channelID)
	if db.HasChannel(channelID) {
		data, err := db.connect.Get(key, nil)
		if err != nil {
			return false
		}
		err = json.Unmarshal(data, &p)
		if err != nil {
			return false
		}
		if p.Public {
			return true
		}
		for i := range p.Admins {
			if p.Admins[i].Equal(member) {
				return true
			}
		}
		for i := range p.Members {
			if p.Members[i].Equal(member) {
				return true
			}
		}
	}
	return false
}

// IsAdmin is the implementation of DB
func (db *LevelDB) IsAdmin(channelID string, member *core.Member) bool {
	var p cc.Profile
	var key = getChannelProfileKey(channelID)
	if db.HasChannel(channelID) {
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
	}
	return false
}

// WatchChannel is the implementation of DB
func (db *LevelDB) WatchChannel(channelID string) {
	db.hub.Watch(channelID, nil)
}

// Close is the implementation of DB
func (db *LevelDB) Close() error {
	return db.connect.Close()
}

func getChannelProfileKey(id string) []byte {
	return []byte(fmt.Sprintf("%s@profile", id))
}

func getSystemAdminKey() []byte {
	return []byte(fmt.Sprintf("%s$admin", core.CONFIGCHANNELID))
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
	account := common.NewAccount(address)
	key := getAccountKey(address)
	data, err := db.connect.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = nil
		}
		return *account, err
	}
	err = json.Unmarshal(data, &account)
	return *account, err
}

// Get get the value by key
func (db *LevelDB) Get(key []byte, couldBeEmpty bool) ([]byte, error) {
	val, err := db.connect.Get(key, nil)
	if err == leveldb.ErrNotFound && couldBeEmpty {
		err = nil
	}
	return val, err
}

// WriteBatchWrapper is a wrapper of level.Batch
type WriteBatchWrapper struct {
	batch *leveldb.Batch
	db    *LevelDB
	kvs   map[string][]byte
}

// NewWriteBatch implement the interface, WriteBatch is a wrapper of leveldb.Batch
func (db *LevelDB) NewWriteBatch() WriteBatch {
	batch := new(leveldb.Batch)
	return &WriteBatchWrapper{
		batch: batch,
		db:    db,
		kvs:   make(map[string][]byte),
	}
}

// AddBlock will records all txs in the block to get rid of duplicated txs
func (wb *WriteBatchWrapper) AddBlock(block *core.Block) error {
	for _, tx := range block.Transactions {
		key := util.BytesCombine([]byte(block.Header.ChannelID), []byte(tx.ID))
		if exist, _ := wb.db.connect.Has(key, nil); exist {
			return fmt.Errorf("The tx %s exists before", tx.ID)
		}
		wb.batch.Put(key, []byte{1})
	}
	return nil
}

// SetConsensusBlock record consensus block
// Note: ignore if num <= 1
func (wb *WriteBatchWrapper) SetConsensusBlock(id string, num uint64) {
	if num > 1 {
		key := getConsensusBlockKey(id)
		wb.batch.Put(key, util.Uint64ToBytes(num))
	}
}

// UpdateChannel is the implementation of DB
func (wb *WriteBatchWrapper) UpdateChannel(id string, profile *cc.Profile) error {
	var key = getChannelProfileKey(id)
	if !wb.db.HasChannel(id) {
		// 更新key为_config的记录，简单记录所有的test通道。 _config,  ["test11","test10","test21","test20"]
		err := wb.addChannel(id)
		if err != nil {
			return err
		}
	}
	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	wb.batch.Put(key, data)
	wb.db.hub.Done(id, nil)
	return nil
}

// Put put updated value in writebatch
func (wb *WriteBatchWrapper) Put(key, value []byte) {
	wb.batch.Put(key, value)
}

// Sync sync batch to database
func (wb *WriteBatchWrapper) Sync() error {
	for k, v := range wb.kvs {
		wb.batch.Put([]byte(k), v)
	}
	return wb.db.connect.Write(wb.batch, nil)
}

// UpdateAccounts update asset
func (wb *WriteBatchWrapper) UpdateAccounts(accounts ...common.Account) error {
	for _, acc := range accounts {
		if err := wb.SetAccount(acc); err != nil {
			return err
		}
	}
	return nil
}

// SetAccount can only be called when atomicity is at one account level
func (wb *WriteBatchWrapper) SetAccount(account common.Account) error {
	key := getAccountKey(account.GetAddress())
	data, err := json.Marshal(account)
	if err != nil {
		return err
	}
	wb.Put(key, data)
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

// addChannel add a record into key core.CONFIGCHANNELID
func (wb *WriteBatchWrapper) addChannel(id string) error {
	var key = getChannelListKey()
	var ids []string
	if util.Contain(wb.kvs, string(key)) {
		json.Unmarshal(wb.kvs[string(key)], &ids)
	} else {
		exist, _ := wb.db.connect.Has(key, nil)
		if exist {
			data, err := wb.db.connect.Get(key, nil)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &ids)
			if err != nil {
				return err
			}
		}
	}
	if !util.Contain(ids, id) {
		ids = append(ids, id)
	}
	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}
	wb.kvs[string(key)] = data
	return nil
}

func getAccountKey(address common.Address) []byte {
	return []byte(fmt.Sprintf("%s", address.String()))
}

func getAssetAdminKey() []byte {
	return []byte("_asset_admin")
}

func getConsensusBlockKey(id string) []byte {
	return []byte(fmt.Sprintf("cbn:%s", id))
}

func getChannelListKey() []byte {
	return []byte("channelList")
}
