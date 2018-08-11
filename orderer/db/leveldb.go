package db

import (
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
)

// GolevelDB is the implementation of DB
type GolevelDB struct {
	// the dir of data
	dir     string
	connect *leveldb.DB
}

// NewGolevelDB is the constructor of GolevelDB
func NewGolevelDB(dir string) (DB, error) {
	db := new(GolevelDB)
	db.dir = dir
	connect, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, err
	}
	db.connect = connect
	return db, nil
}

// ListChannel is the implementation of DB
// TODO
func (db *GolevelDB) ListChannel() []string {
	var channels []string
	return channels
}

// AddChannel is the implementation of DB
// TODO
func (db *GolevelDB) AddChannel(id string) error {
	return errors.New("Not implementation yet")
}
