package eraft

import (
	"madledger/common/util"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testChannel = "_db_test"
)

func TestDB(t *testing.T) {
	os.RemoveAll(getDBPath())

	db, err := NewDB(getDBPath())
	require.NoError(t, err)

	_, err = NewDB(getDBPath())
	require.Errorf(t, err, "resource temporarily unavailable")
	db.Close()
}

func TestBlock(t *testing.T) {
	db, err := NewDB(getDBPath())
	require.NoError(t, err)

	require.Equal(t, uint64(0), db.GetMinBlock(testChannel))
	blockSize := 100
	blocks := make([]*Block, blockSize)
	for i := range blocks {
		blocks[i] = &Block{
			ChannelID: testChannel,
			Num:       uint64(i),
		}
	}

	for i := range blocks {
		db.AddBlock(blocks[i])
		db.SetMinBlock(testChannel, uint64(i))
	}
	require.Equal(t, uint64(blockSize-1), db.GetMinBlock(testChannel))

	db.Close()
}

func TestEnd(t *testing.T) {
	require.NoError(t, os.RemoveAll(getDBPath()))
}

func getDBPath() string {
	gopath := os.Getenv("GOPATH")
	storePath, _ := util.MakeFileAbs("src/madledger/blockchain/chain/raft/eraft/.db", gopath)
	return storePath
}
