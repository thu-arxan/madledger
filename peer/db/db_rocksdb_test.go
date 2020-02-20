// +build rocksdb

package db

var (
	dbConstructFunc = []func(dir string) (DB, error){NewLevelDB, NewRocksDB}
)
