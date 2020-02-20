// +build !rocksdb

package server

import (
	"madledger/peer/db"
)

func newDB(dir string) (db.DB, error) {
	log.Infof("using LevelDB: %s", dir)
	return db.NewLevelDB(dir)
}
