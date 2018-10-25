package tendermint

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
)

// DB will record some necessary thing such as height and hash of tendermint
type DB struct {
	connect *leveldb.DB
}

// NewDB is the constructor of DB
func NewDB(dir string) (*DB, error) {
	db := new(DB)
	connect, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, err
	}
	db.connect = connect
	return db, nil
}

// GetHeight get the height
func (db *DB) GetHeight() int64 {
	var key = []byte("height")
	var height int64
	if exist, _ := db.connect.Has(key, nil); exist {
		data, _ := db.connect.Get(key, nil)
		json.Unmarshal(data, &height)
	}
	return height
}

// SetHeight set the height
func (db *DB) SetHeight(height int64) {
	var key = []byte("height")
	data, _ := json.Marshal(height)
	db.connect.Put(key, data, nil)
}

// GetHash return the hash
func (db *DB) GetHash() []byte {
	var key = []byte("hash")
	if exist, _ := db.connect.Has(key, nil); exist {
		data, _ := db.connect.Get(key, nil)
		return data
	}
	return nil
}

// SetHash set the hash
func (db *DB) SetHash(hash []byte) {
	var key = []byte("hash")
	db.connect.Put(key, hash, nil)
}
