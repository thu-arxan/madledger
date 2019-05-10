package tendermint

import (
	"encoding/json"
	"fmt"
	"madledger/common/util"

	"github.com/syndtr/goleveldb/leveldb"
)

// DB will record some necessary thing such as height and hash of tendermint
type DB struct {
	dir     string
	connect *leveldb.DB
}

// NewDB is the constructor of DB
func NewDB(dir string) (*DB, error) {
	db := new(DB)
	db.dir = dir

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

// AddBlock add a block
func (db *DB) AddBlock(block *Block) error {
	var key = []byte(fmt.Sprintf("%s:%d", block.ChannelID, block.Num))
	data, err := json.Marshal(block)
	if err != nil {
		return err
	}
	err = db.connect.Put(key, data, nil)
	if err != nil {
		return err
	}
	db.SetChannelBlockNumber(block.ChannelID, block.Num)
	return nil
}

// GetBlock get a block
func (db *DB) GetBlock(channelID string, num uint64) *Block {
	var key = []byte(fmt.Sprintf("%s:%d", channelID, num))
	if exist, _ := db.connect.Has(key, nil); exist {
		data, _ := db.connect.Get(key, nil)
		var block Block
		if err := json.Unmarshal(data, &block); err != nil {
			return nil
		}
		return &block
	}
	return nil
}

// SetHash set the hash
func (db *DB) SetHash(hash []byte) {
	var key = []byte("hash")
	db.connect.Put(key, hash, nil)
}

// Close close the connection if exist
func (db *DB) Close() error {
	if db.connect != nil {
		return db.connect.Close()
	}
	return nil
}

// SetChannelBlockNumber set the block number of channel
func (db *DB) SetChannelBlockNumber(channelID string, num uint64) {
	var key = []byte(fmt.Sprintf("number:%s", channelID))
	data := util.Uint64ToBytes(num)
	db.connect.Put(key, data, nil)
}

// GetChannelBlockNumber return the block number of channel
func (db *DB) GetChannelBlockNumber(channelID string) uint64 {
	var key = []byte(fmt.Sprintf("number:%s", channelID))
	if exist, _ := db.connect.Has(key, nil); exist {
		data, _ := db.connect.Get(key, nil)
		num, err := util.BytesToUint64(data)
		if err != nil {
			return 0
		}
		return num
	}
	return 0
}
