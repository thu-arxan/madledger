package tendermint

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
	"testing"
)

func TestDB_GetDB(t *testing.T) {
	path := "/home/hadoop/GOPATH/src/madledger/env/bft/orderers/0/data/leveldb"
	db, err := leveldb.OpenFile(path, nil)
	fmt.Printf("Get bft.db from %s\n", path)
	require.NoError(t, err)
	defer db.Close()
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := string(iter.Key())
		value := string(iter.Value())
		fmt.Printf("(%s, %s)\n", key, value)
	}
	iter.Release()
}
