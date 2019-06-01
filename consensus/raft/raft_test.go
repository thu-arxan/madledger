package raft

import (
	"fmt"
	"madledger/common/util"
	"madledger/consensus"
	"os"
	"testing"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/require"
)

// This file will start some raft nodes and run basic tests

var (
	// raft nodes
	rns   [3]consensus.Consensus
	peers = map[uint64]string{
		1: "127.0.0.1:12345",
		2: "127.0.0.1:12347",
		3: "127.0.0.1:12349",
	}
	// txSize records txs send
	txSize int
)

func TestInitEnv(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(getTestPath()))

	require.NoError(t, copy.Copy(gopath+"/src/madledger/env/bft/.orderers", fmt.Sprintf("%s/orderers", getTestPath())))
}

func TestStart(t *testing.T) {
	for i := range rns {
		cfg, err := getConfig(i)
		require.NoError(t, err)
		node, err := NewConseneus(cfg)
		require.NoError(t, err)
		rns[i] = node
		require.NoError(t, node.Start())
	}
}

func TestAddTx(t *testing.T) {

}

func TestStop(t *testing.T) {

}

func getTestPath() string {
	gopath := os.Getenv("GOPATH")
	testPath, _ := util.MakeFileAbs("src/madledger/consensus/raft/.test", gopath)
	return testPath
}

func getNodePath(node int) string {
	return fmt.Sprintf("%s/orderers/%d", getTestPath(), node)
}

func randomTx() []byte {
	return []byte(util.RandomString(32))
}

func getConfig(node int) (*Config, error) {
	return NewConfig(getNodePath(node), "127.0.0.1", uint64(node+1), peers, consensus.Config{
		Timeout: 100,
		MaxSize: 10,
		Resume:  false,
		Number:  1,
	})
}
