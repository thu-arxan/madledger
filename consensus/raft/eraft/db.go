// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package eraft

import (
	"encoding/json"
	"fmt"
	"madledger/common/util"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

// DB is the database of raft
type DB struct {
	// todo: rwLock?
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

// GetMinBlock return min block number of channelID for restore
func (db *DB) GetMinBlock(channelID string) uint64 {
	var key = []byte("minBlock_" + channelID)

	data, err := db.connect.Get(key, nil)
	if err != nil {
		if err != leveldb.ErrNotFound {
			log.Errorf("get minblock of channel %s failed: %v", channelID, err)
		}
		return 0
	}

	num, err := util.BytesToUint64(data)
	if err != nil {
		log.Errorf("bytes to uint64 failed: %v", err)
		return 0
	}
	return num
}

// SetMinBlock set min block number
func (db *DB) SetMinBlock(channelID string, num uint64) {
	var key = []byte("minBlock_" + channelID)
	db.connect.Put(key, util.Uint64ToBytes(num), nil)
}

// GetChainNum return the block num(height) of channel
func (db *DB) GetChainNum(channelID string) uint64 {
	var key = []byte("chainNum_" + channelID)

	data, err := db.connect.Get(key, nil)
	if err != nil {
		if err != leveldb.ErrNotFound {
			log.Errorf("get chainNum of channel %s failed: %v", channelID, err)
		}
		return 0
	}

	num, err := util.BytesToUint64(data)
	if err != nil {
		log.Errorf("bytes to uint64 failed: %v", err)
		return 0
	}
	return num
}

// SetChainNum set the block num(height) of channel
func (db *DB) SetChainNum(channelID string, num uint64) {
	var key = []byte("chainNum_" + channelID)
	db.connect.Put(key, util.Uint64ToBytes(num), nil)
}

// AddBlock add block
func (db *DB) AddBlock(block *Block) {
	var key = []byte(fmt.Sprintf("block_%s:%d", block.ChannelID, block.Num))
	log.Infof("Add block into raft.db: %s, %d", block.ChannelID, block.Num)
	db.connect.Put(key, block.Bytes(), nil)
}

// GetBlock return the block of channel, return nil if not exist
func (db *DB) GetBlock(channelID string, num uint64, async bool) *Block {
	for {
		block := db.getBlock(channelID, num)
		if block != nil {
			return block
		}
		if !async {
			return nil
		}
		time.Sleep(30 * time.Millisecond)
	}
}

// getBlock return the block of channel, return nil if not exist
func (db *DB) getBlock(channelID string, num uint64) *Block {
	var key = []byte(fmt.Sprintf("block_%s:%d", channelID, num))
	data, err := db.connect.Get(key, nil)
	if err != nil {
		if err != leveldb.ErrNotFound && err != leveldb.ErrClosed {
			log.Errorf("get block %s:%d from db failed: %v", channelID, num, err)
		}
		return nil
	}

	var block Block
	json.Unmarshal(data, &block)
	return &block
}
