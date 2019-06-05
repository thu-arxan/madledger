package raft

import (
	"encoding/json"
	"errors"
	"madledger/common/util"

	"github.com/syndtr/goleveldb/leveldb"
)

type dbBlocks struct {
	Blocks []*HybridBlock
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

// PutBlock put a block into db
func (db *DB) PutBlock(block *HybridBlock) {
	var key = util.Uint64ToBytes(block.GetNumber())
	db.connect.Put(key, block.Bytes(), nil)
}

// GetBlock return the block which Num is num
// Return nil, errors.New("Not exist") if not exist
func (db *DB) GetBlock(num uint64) (*HybridBlock, error) {
	var key = util.Uint64ToBytes(num)
	has, err := db.connect.Has(key, nil)
	if err != nil {
		return nil, err
	}
	if has {
		value, err := db.connect.Get(key, nil)
		if err != nil {
			return nil, err
		}
		return UnmarshalHybridBlock(value), nil
	}

	return nil, errors.New("Not exist")
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
