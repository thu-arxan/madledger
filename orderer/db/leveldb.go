package db

import (
	"encoding/json"
	"fmt"
	cc "madledger/blockchain/config"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/event"
	"madledger/common/util"
	"madledger/core"
	"reflect"

	"github.com/syndtr/goleveldb/leveldb"
)

/*
*  1. Channel profile: key is []byte("_config@" + channelID), value is the json.Marshal(profile)
*  2. All channel ids: key is []byte("_config"), value is json.Marshl([]string{id1, id2, ...})
*  3. Tx: key is combine of []byte(channelID) and []byte(txID), value is []byte("true")
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
	var key = []byte(core.CONFIGCHANNELID)
	var channels []string
	data, err := db.connect.Get(key, nil)
	if err != nil {
		return channels
	}
	err = json.Unmarshal(data, &channels)
	if err != nil {
		return channels
	}
	return channels
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
	// "Admins":[{"PK":"BN2PLBpBd5BrSLfTY7QEBYQT0h6lFvWlZyuAVt3/bfEz1g5QJ2lIEXP2Zk15B6E2MWpA/Q4Yxnl+XjFGObvAKTY=","Name":"admin"}]}
	err = db.connect.Put(key, data, nil)
	if err != nil {
		return err
	}
	db.hub.Done(id, nil)
	return nil
}

// AddBlock will records all txs in the block to get rid of duplicated txs
func (db *LevelDB) AddBlock(block *core.Block) error {
	for _, tx := range block.Transactions {
		key := util.BytesCombine([]byte(block.Header.ChannelID), []byte(tx.ID))
		if exist, _ := db.connect.Has(key, nil); exist {
			return fmt.Errorf("The tx %s exists before", tx.ID)
		}
		db.connect.Put(key, []byte("true"), nil)
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

// WatchChannel is the implementation of DB
func (db *LevelDB) WatchChannel(channelID string) {
	db.hub.Watch(channelID, nil)
}

// Close is the implementation of DB
func (db *LevelDB) Close() error {
	return db.connect.Close()
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

func getChannelProfileKey(id string) []byte {
	return []byte(fmt.Sprintf("%s@%s", core.CONFIGCHANNELID, id))
}

func getSystemAdminKey() []byte {
	return []byte(fmt.Sprintf("%s$admin", core.CONFIGCHANNELID))
}

// IsAssetAdmin determines whether input pk belonged to account that has the right to issue
func (db *LevelDB) IsAccountAdmin(pk crypto.PublicKey) bool {
	var key = []byte("_account_admin")
	admin, err := db.connect.Get(key, nil)
	if err != nil {
		return false
	}
	pkBytes, err := pk.Bytes()
	if err != nil {
		return false
	}
	return reflect.DeepEqual(admin, pkBytes)
}

// SetAccountAdmin only succeed at the first time it is called
func (db *LevelDB) SetAccountAdmin(pk crypto.PublicKey) error {
	var key = []byte("_account_admin")
	exists, _ := db.connect.Has(key, nil)
	if exists {
		return fmt.Errorf("account admin already set")
	}
	pkBytes, err := pk.Bytes()
	if err != nil {
		return err
	}
	return db.connect.Put(key, pkBytes, nil)
}

// GetOrCreateAccount return default account if account does not exist in leveldb
func (db *LevelDB) GetOrCreateAccount(address common.Address) (common.Account, error) {
	key := getAccountKey(address)
	var account common.DefaultAccount
	data, err := db.connect.Get(key, nil)
	if err != nil {
		if err != leveldb.ErrNotFound {
			return nil, err
		}
		return common.NewDefaultAccount(address), nil
	}
	err = json.Unmarshal(data, &account)
	return &account, err
}

func (db *LevelDB) UpdateAccounts(accounts ...common.Account) error {
	wb := &leveldb.Batch{}
	for _, acc := range accounts {
		key := getAccountKey(acc.GetAddress())
		data, err := json.Marshal(acc)
		if err != nil {
			return err
		}
		wb.Put(key, data)
	}
	return db.connect.Write(wb, nil)
}

func (db *LevelDB) IsTxExecute(txid string) bool {
	key := []byte(txid)
	data, err := db.connect.Get(key, nil)
	if err != nil {
		return false
	}
	if string(data[:]) != "ok" {
		return false
	}
	return true
}

func (db *LevelDB) SetTxExecute(txid string) error {
	key := []byte(txid)
	return db.connect.Put(key, []byte("ok"), nil)
}

func getAccountKey(address common.Address) []byte {
	return []byte(fmt.Sprintf("%s@%s", core.ASSETCHANNELID, address.String()))
}
