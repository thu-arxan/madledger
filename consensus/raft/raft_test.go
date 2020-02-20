package raft

import (
	"fmt"
	"madledger/common/util"
	"madledger/consensus"
	"madledger/core"
	"os"
	"sync"
	"testing"
	"time"

	"net/http"
	_ "net/http/pprof"

	"golang.org/x/sync/errgroup"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/require"
)

// This file will start some raft nodes and run basic tests

var (
	// raft nodes
	nodes [3]consensus.Consensus
	peers = map[uint64]string{
		1: "127.0.0.1:12345",
		2: "127.0.0.1:12347",
		3: "127.0.0.1:12349",
	}
	// txSize records txs send
	txSize int
)

func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

func TestInitEnv(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(getTestPath()))

	require.NoError(t, copy.Copy(gopath+"/src/madledger/env/bft/.orderers", fmt.Sprintf("%s/orderers", getTestPath())))
}

func TestStart(t *testing.T) {
	var wg sync.WaitGroup
	for i := range nodes {
		cfg, err := getConfig(i)
		require.NoError(t, err)
		node, err := NewConsensus(map[string]consensus.Config{
			"_global": cfg.cc,
		}, cfg)
		require.NoError(t, err)
		nodes[i] = node
		go func() {
			wg.Add(1)
			defer wg.Done()
			require.NoError(t, node.Start())
		}()
	}
	time.Sleep(200 * time.Millisecond)
	wg.Wait()
}

func TestAddTx(t *testing.T) {
	txSize = 128
	var txs []*core.Tx

	for i := 0; i < txSize; i++ {
		tx := randomTx()
		txs = append(txs, tx)
	}

	var g errgroup.Group
	for i := range txs {
		tx := txs[i]
		g.Go(func() error {
			n := util.RandNum(len(nodes))
			if err := nodes[n].AddTx(tx); err != nil {
				return err
			}
			return nil
		})

	}
	require.NoError(t, g.Wait())
}

func TestStop(t *testing.T) {
	for i := range nodes {
		fmt.Printf("Stop raft %d\n", i)
		require.NoError(t, nodes[i].Stop())
	}
	os.RemoveAll(getTestPath())
}

func TestReinitEnv(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(getTestPath()))

	require.NoError(t, copy.Copy(gopath+"/src/madledger/env/bft/.orderers", fmt.Sprintf("%s/orderers", getTestPath())))
}

func TestRestart(t *testing.T) {
	var g errgroup.Group

	for i := range nodes {
		cfg, err := getConfig(i)
		require.NoError(t, err)
		node, err := NewConsensus(nil, cfg)
		require.NoError(t, err)
		nodes[i] = node
		g.Go(func() error {
			return node.Start()
		})
	}
	require.NoError(t, g.Wait())
}

func TestReStop(t *testing.T) {
	for i := range nodes {
		fmt.Printf("Stop raft %d\n", i)
		require.NoError(t, nodes[i].Stop())
	}
	os.RemoveAll(getTestPath())
}

func getTestPath() string {
	gopath := os.Getenv("GOPATH")
	testPath, _ := util.MakeFileAbs("src/madledger/consensus/raft/.test", gopath)
	return testPath
}

func getNodePath(node int) string {
	return fmt.Sprintf("%s/orderers/%d", getTestPath(), node)
}

func randomTx() *core.Tx {
	return &core.Tx{
		ID: util.RandomString(32),
		Data: core.TxData{
			ChannelID: "_global",
		},
	}
}

func getConfig(node int) (*Config, error) {
	return NewConfig(getNodePath(node), "127.0.0.1", uint64(node+1), peers, false, consensus.Config{
		Timeout: 100,
		MaxSize: 10,
		Resume:  false,
		Number:  1,
	})
}
