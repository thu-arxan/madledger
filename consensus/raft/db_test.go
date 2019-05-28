package raft

import (
	"madledger/common/util"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T) {
	os.RemoveAll(getDBPath())

	db, err := NewDB(getDBPath())
	require.NoError(t, err)

	_, err = NewDB(getDBPath())
	require.Errorf(t, err, "resource temporarily unavailable")
	db.Close()
}

// func TestBlock(t *testing.T) {
// 	db, err := NewDB(getDBPath())
// 	require.NoError(t, err)

// 	blocks := db.GetBlocks()
// 	require.Len(t, blocks, 0)

// 	blocks = make([]*core.Block, 100)
// 	genesisBlockPrevHash := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
// 	blocks[0] = core.NewBlock(0, genesisBlockPrevHash, nil)
// 	for i := range blocks {
// 		if i != 0 {
// 			blocks[i] = core.NewBlock(uint64(i), blocks[i-1].Hash(), nil)
// 		}
// 		if i != 10 {
// 			db.AddBlock(blocks[i])
// 		}
// 	}
// 	db.AddBlock(blocks[10])

// 	blocks = db.GetBlocks()
// 	require.Len(t, blocks, 100)
// 	require.Equal(t, uint64(0), blocks[0].GetNumber())
// 	require.Equal(t, uint64(10), blocks[10].GetNumber())
// 	require.Equal(t, uint64(99), blocks[99].GetNumber())

// 	require.Equal(t, uint64(0), db.GetMinBlock())
// 	db.RemoveBlocks(20)
// 	require.Equal(t, uint64(21), db.GetMinBlock())
// 	blocks = db.GetBlocks()
// 	for i := range blocks {
// 		require.Equal(t, blocks[i].GetNumber(), uint64(i+21))
// 	}

// 	db.Close()
// }

func TestEnd(t *testing.T) {
	require.NoError(t, os.RemoveAll(getDBPath()))
}

func getDBPath() string {
	gopath := os.Getenv("GOPATH")
	storePath, _ := util.MakeFileAbs("src/transaction_service/blockchain/chain/raft/.db", gopath)
	return storePath
}
