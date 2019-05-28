package raft

import (
	"encoding/json"
	"madledger/common/util"
	core "madledger/core/types"

	"github.com/syndtr/goleveldb/leveldb"
)

type dbBlocks struct {
	Blocks []*core.Block
}

// DB is the database of raft
type DB struct {
	dir     string
	connect *leveldb.DB
}

// NewDB is the constructor of DB
func NewDB(dir string) (*DB, error) {
	var err error

	db := new(DB)
	db.dir = dir
	db.connect, err = leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Close will close the connection
func (db *DB) Close() {
	db.connect.Close()
}

// AddBlock add a block into db
func (db *DB) AddBlock(block *core.Block) {
	var key = util.Uint64ToBytes(block.GetNumber())
	db.connect.Put(key, block.Bytes(), nil)
}

// GetMinBlock return min block number for restore
func (db *DB) GetMinBlock() uint64 {
	var key = []byte("minBlock")
	var num uint64
	if exist, _ := db.connect.Has(key, nil); exist {
		data, _ := db.connect.Get(key, nil)
		json.Unmarshal(data, &num)
	}
	return num
}

// SetMinBlock set min block number
func (db *DB) SetMinBlock(num uint64) {
	var key = []byte("minBlock")
	data, _ := json.Marshal(num)
	db.connect.Put(key, data, nil)
}
