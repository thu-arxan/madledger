package db

import (
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	cc "madledger/blockchain/config"
	"madledger/common/event"
	"madledger/common/util"
	"madledger/core/types"
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
	var key = []byte(types.CONFIGCHANNELID)
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
	//更新key为_config@id的记录, 具体内容示例_config@test30 ,  {"Public":true,"Dependencies":null,"Members":[],"Admins":[{"PK":"BN2PLBpBd5BrSLfTY7QEBYQT0h6lFvWlZyuAVt3/bfEz1g5QJ2lIEXP2Zk15B6E2MWpA/Q4Yxnl+XjFGObvAKTY=","Name":"admin"}]}
	err = db.connect.Put(key, data, nil)
	if err != nil {
		return err
	}
	db.hub.Done(id, nil)
	return nil
}

// AddBlock will records all txs in the block to get rid of duplicated txs
func (db *LevelDB) AddBlock(block *types.Block) error {
	for _, tx := range block.Transactions {
		key := util.BytesCombine([]byte(block.Header.ChannelID), []byte(tx.ID))
		if exist, _ := db.connect.Has(key, nil); exist {
            // 打印数据库中已经包含的键值对
			/*iter := db.connect.NewIterator(nil, nil)
			for iter.Next() {
				key   := string(iter.Key())
				value := string(iter.Value())
				log.Info("db.AddBlock： (", key,", ", value,")")
			}
			iter.Release()*/

			// 打印tx的添加时间
			/*log.Info("db.AddBlock： add ",block.Header.ChannelID," block ",block.Header.Number,", the tx ",tx.ID,
				"exists before, and its time is ",tx.Time, ".")*/
			return fmt.Errorf("The tx %s exists before", tx.ID)
		}
		db.connect.Put(key, []byte("true"), nil)
		// 打印tx的添加时间
		/*log.Info("db.AddBlock： add ",block.Header.ChannelID," block ",block.Header.Number,", put (", string(key),
			", true) into db, and its time is ",tx.Time, ".")*/
	}
	return nil
}

// HasTx return if the tx is contained
func (db *LevelDB) HasTx(tx *types.Tx) bool {
	key := util.BytesCombine([]byte(tx.Data.ChannelID), []byte(tx.ID))

	if exist, _ := db.connect.Has(key, nil); exist {
		return true
	}
	return false
}

// IsMember is the implementation of DB
func (db *LevelDB) IsMember(channelID string, member *types.Member) bool {
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
func (db *LevelDB) IsAdmin(channelID string, member *types.Member) bool {
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

// addChannel add a record into key types.CONFIGCHANNELID
func (db *LevelDB) addChannel(id string) error {
	var key = []byte(types.CONFIGCHANNELID)
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
	return []byte(fmt.Sprintf("%s@%s", types.CONFIGCHANNELID, id))
}
