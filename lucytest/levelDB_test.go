package lucytest

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
	"testing"
)

func TestLevelDB(t *testing.T){
	db,err:=leveldb.OpenFile("/home/zebra/GOPATH/src/madledger/env/raft/peers/0/data/leveldb",nil)
	require.NoError(t,err)
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		fmt.Printf("key:%s, value:%s\n", string(iter.Key()), string(iter.Value()))
	}
	iter.Release()

}
