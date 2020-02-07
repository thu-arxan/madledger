package eraft

// import (
// 	"fmt"
// 	"madledger/common/util"
// 	"madledger/consensus"
// 	"os"
// 	"sync"
// 	"testing"
// 	"time"

// 	"github.com/otiai10/copy"
// 	"github.com/stretchr/testify/require"
// )

// // This file will start some raft nodes and run basic tests

// var (
// 	// raft nodes
// 	rns   [3]consensus.Consensus
// 	peers = map[uint64]string{
// 		1: "127.0.0.1:12345",
// 		2: "127.0.0.1:12347",
// 		3: "127.0.0.1:12349",
// 	}
// 	// txSize records txs send
// 	txSize int
// )

// func TestInitEnv(t *testing.T) {
// 	gopath := os.Getenv("GOPATH")
// 	require.NoError(t, os.RemoveAll(getTestPath()))

// 	require.NoError(t, copy.Copy(gopath+"/src/madledger/env/bft/.orderers", fmt.Sprintf("%s/orderers", getTestPath())))
// }

// func TestStart(t *testing.T) {
// 	for i := range rns {
// 		cfg, err := getConfig(i)
// 		require.NoError(t, err)
// 		node, err := NewConseneus(cfg)
// 		require.NoError(t, err)
// 		rns[i] = node
// 		go func() {
// 			require.NoError(t, node.Start())
// 		}()
// 	}
// 	time.Sleep(2 * time.Second)
// }

// func TestAddTx(t *testing.T) {
// 	txSize = 128
// 	var txs [][]byte
// 	var success = make(map[string]int)
// 	var lock sync.Mutex

// 	for i := 0; i < txSize; i++ {
// 		tx := randomTx()
// 		success[string(tx)] = 0
// 		txs = append(txs, tx)
// 	}

// 	var wg sync.WaitGroup
// 	for i := range txs {
// 		wg.Add(1)
// 		tx := txs[i]
// 		go func() {
// 			defer wg.Done()
// 			n := util.RandNum(len(rns))
// 			if err := rns[n].AddTx("test", tx); err == nil {
// 				lock.Lock()
// 				success[string(tx)]++
// 				lock.Unlock()
// 			}
// 		}()
// 	}
// 	wg.Wait()

// 	for i := range success {
// 		require.Equal(t, 1, success[i])
// 	}
// }

// func TestStop(t *testing.T) {
// 	for i := range rns {
// 		fmt.Printf("Stop raft %d\n", i)
// 		rns[i].Stop()
// 	}
// 	os.RemoveAll(getTestPath())
// }

// func TestReinitEnv(t *testing.T) {
// 	gopath := os.Getenv("GOPATH")
// 	require.NoError(t, os.RemoveAll(getTestPath()))

// 	require.NoError(t, copy.Copy(gopath+"/src/madledger/env/bft/.orderers", fmt.Sprintf("%s/orderers", getTestPath())))
// }

// func TestRestart(t *testing.T) {
// 	for i := range rns {
// 		cfg, err := getConfig(i)
// 		require.NoError(t, err)
// 		node, err := NewConseneus(cfg)
// 		require.NoError(t, err)
// 		rns[i] = node
// 		go func() {
// 			require.NoError(t, node.Start())
// 		}()
// 	}
// 	time.Sleep(2 * time.Second)
// }

// func TestReStop(t *testing.T) {
// 	for i := range rns {
// 		fmt.Printf("Stop raft %d\n", i)
// 		rns[i].Stop()
// 	}
// 	os.RemoveAll(getTestPath())
// }
// func getTestPath() string {
// 	gopath := os.Getenv("GOPATH")
// 	testPath, _ := util.MakeFileAbs("src/madledger/consensus/raft/.test", gopath)
// 	return testPath
// }

// func getNodePath(node int) string {
// 	return fmt.Sprintf("%s/orderers/%d", getTestPath(), node)
// }

// func randomTx() []byte {
// 	return []byte(util.RandomString(32))
// }

// func getConfig(node int) (*Config, error) {
// 	return NewConfig(getNodePath(node), "127.0.0.1", uint64(node+1), peers, false, consensus.Config{
// 		Timeout: 100,
// 		MaxSize: 10,
// 		Resume:  false,
// 		Number:  1,
// 	})
// }
