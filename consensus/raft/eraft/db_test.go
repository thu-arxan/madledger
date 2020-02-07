package eraft

// import (
// 	"madledger/common/util"
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/require"
// )

// func TestDB(t *testing.T) {
// 	os.RemoveAll(getDBPath())

// 	db, err := NewDB(getDBPath())
// 	require.NoError(t, err)

// 	_, err = NewDB(getDBPath())
// 	require.Errorf(t, err, "resource temporarily unavailable")
// 	db.Close()
// }

// func TestBlock(t *testing.T) {
// 	db, err := NewDB(getDBPath())
// 	require.NoError(t, err)

// 	require.Equal(t, uint64(0), db.GetMinBlock())
// 	blockSize := 100
// 	blocks := make([]*HybridBlock, blockSize)
// 	for i := range blocks {
// 		blocks[i] = &HybridBlock{
// 			Num: uint64(i),
// 		}
// 	}

// 	for i := range blocks {
// 		db.PutBlock(blocks[i])
// 		db.SetMinBlock(uint64(i))
// 	}
// 	require.Equal(t, uint64(blockSize-1), db.GetMinBlock())

// 	db.Close()
// }

// func TestEnd(t *testing.T) {
// 	require.NoError(t, os.RemoveAll(getDBPath()))
// }

// func getDBPath() string {
// 	gopath := os.Getenv("GOPATH")
// 	storePath, _ := util.MakeFileAbs("src/madledger/blockchain/chain/raft/.db", gopath)
// 	return storePath
// }
