// +build rocksdb

package server

import (
	"madledger/peer/db"
)

func newDB(dir string) (db.DB, error) {
	log.Infof("using RocksDB: %s", dir)
	return db.NewRocksDB(dir)
}
