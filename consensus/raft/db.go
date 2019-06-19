package raft

import (
	"encoding/json"
	"errors"
	"fmt"
	"madledger/common/util"
	"time"

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
	log.Infof("Put hybridBlock %d into raft.db", block.GetNumber())
	db.connect.Put(key, block.Bytes(), nil)
}

// GetHybridBlock return the hybrid block which Num is num
// Return nil, errors.New("Not exist") if not exist
func (db *DB) GetHybridBlock(num uint64) (*HybridBlock, error) {
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
	log.Infof("GetMinBlock: get minBlock %d", num)
	return num
}

// SetMinBlock set min block number
func (db *DB) SetMinBlock(num uint64) {
	var key = []byte("minBlock")
	data, _ := json.Marshal(num)
	db.connect.Put(key, data, nil)
}

// GetPrevBlockNum return the prev block num of channel
func (db *DB) GetPrevBlockNum(channelID string) uint64 {
	var key = []byte(channelID)
	var num uint64
	if exist, _ := db.connect.Has(key, nil); exist {
		data, _ := db.connect.Get(key, nil)
		json.Unmarshal(data, &num)
	}
	return num
}

// SetPrevBlockNum set the prev block num of channel
func (db *DB) SetPrevBlockNum(channelID string, num uint64) {
	var key = []byte(channelID)
	data, _ := json.Marshal(num)
	db.connect.Put(key, data, nil)
}

// AddBlock add block
func (db *DB) AddBlock(block *Block) {
	var key = []byte(fmt.Sprintf("%s:%d", block.ChannelID, block.Num))
	log.Infof("Add block into raft.db: %s, %d", block.ChannelID, block.Num)
	db.connect.Put(key, block.Bytes(), nil)
}

// GetBlock return the block of channel, return nil if not exist
func (db *DB) GetBlock(channelID string, num uint64, async bool) *Block {
	block := db.getBlock(channelID, num)
	if !async || (block != nil) {
		log.Infof("Get block %d of channel %s from raft.db", num, channelID)
		return block
	}
	for {
		time.Sleep(10 * time.Millisecond)
		block = db.getBlock(channelID, num)
		if block != nil {
			log.Infof("Get block %d of channel %s from raft.db asynchronously", num, channelID)
			return block
		}
	}
}

// getBlock return the block of channel, return nil if not exist
func (db *DB) getBlock(channelID string, num uint64) *Block {
	var key = []byte(fmt.Sprintf("%s:%d", channelID, num))
	if exist, _ := db.connect.Has(key, nil); exist {
		var block Block
		data, _ := db.connect.Get(key, nil)
		json.Unmarshal(data, &block)
		return &block
	}

	return nil
}
