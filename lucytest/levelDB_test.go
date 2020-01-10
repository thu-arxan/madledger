package lucytest

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestLevelDB(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	db, err := leveldb.OpenFile(fmt.Sprintf("%s/src/madledger/env/bft/orderers/0/data/leveldb", gopath), nil)
	require.NoError(t, err)
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		fmt.Printf("key:%s, value:%s\n", string(iter.Key()), string(iter.Value()))
	}
	iter.Release()

}
