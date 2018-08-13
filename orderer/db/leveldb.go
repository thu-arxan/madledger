package db

import (
	"encoding/json"
	"fmt"
	cc "madledger/blockchain/config"
	"madledger/core"

	"github.com/syndtr/goleveldb/leveldb"
)

// TODO: This need a summary or do this in the docs
/*
*
 */

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
// maybe the name should be checked
// TODO
func (db *LevelDB) UpdateChannel(id string, profile cc.Profile) error {
	var p cc.Profile
	var key = getChannelProfileKey(id)
	if db.HasChannel(id) {
		data, err := db.connect.Get(key, nil)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &p)
		if err != nil {
			return err
		}
	} else {
		err := db.addChannel(id)
		if err != nil {
			return err
		}
	}
	// todo: In the future, this maybe wrong
	p.Open = profile.Open
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	err = db.connect.Put(key, data, nil)
	if err != nil {
		return err
	}
	return nil
}

// addChannel add a record into key core.CONFIGCHANNELID
// todo: maybe a map is better, and there is need to check if channel exists aleardy
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
	ids = append(ids, id)
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
