package tendermint

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"madledger/common/util"
	"madledger/consensus"
	"madledger/core/types"
	"madledger/orderer/config"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/require"
)

// This file will start some tendermint nodes and run basic tests

var (
	// tendermint nodes
	tns [4]consensus.Consensus
	// txSize records txs send
	txSize int
)

func TestMain(m *testing.M) {
	ret := m.Run()
	os.Exit(ret)
}

func TestInitEnv(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(getTestPath()))

	require.NoError(t, copy.Copy(gopath+"/src/madledger/env/bft/.orderers", fmt.Sprintf("%s/orderers", getTestPath())))
}

func TestStart(t *testing.T) {
	var channels = make(map[string]consensus.Config)
	channels[types.GLOBALCHANNELID] = consensus.DefaultConfig()
	channels[types.CONFIGCHANNELID] = consensus.DefaultConfig()
	for i := 0; i < len(tns); i++ {
		// update absolute path into orderer.yaml
		err := absBFTOrdererConfig(i)
		require.NoError(t, err)
		cfg, err := config.LoadConfig(getConfigPath(i))
		require.NoError(t, err)

		consensus, err := NewConsensus(channels, &Config{
			Port: Port{
				P2P: cfg.Consensus.Tendermint.Port.P2P,
				RPC: cfg.Consensus.Tendermint.Port.RPC,
				App: cfg.Consensus.Tendermint.Port.APP,
			},
			Dir:        cfg.Consensus.Tendermint.Path,
			P2PAddress: cfg.Consensus.Tendermint.P2PAddress,
		})
		require.NoError(t, err)
		tns[i] = consensus
		require.NoError(t, consensus.Start())
	}
	time.Sleep(5 * time.Second)
}

func absBFTOrdererConfig(node int) error {
	cfgPath := getConfigPath(node)
	// load config
	cfg, err := loadOrdererConfig(cfgPath)
	if err != nil {
		return err
	}
	// change relative path into absolute path
	cfg.BlockChain.Path = getNodePath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Path = getNodePath(node) + "/" + cfg.DB.LevelDB.Path
	cfg.Consensus.Tendermint.Path = getNodePath(node) + "/" + cfg.Consensus.Tendermint.Path
	cfg.TLS.CA = getNodePath(node) + "/" + cfg.TLS.CA
	cfg.TLS.RawCert = getNodePath(node) + "/" + cfg.TLS.RawCert
	cfg.TLS.Key = getNodePath(node) + "/" + cfg.TLS.Key
	// rewrite orderer config
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfgPath, data, os.ModePerm)
}

func loadOrdererConfig(cfgPath string) (*config.Config, error) {
	cfgBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	var cfg config.Config
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func TestSendDuplicateTxs(t *testing.T) {
	txSize = 128
	var txs [][]byte
	var success = make(map[string]int)
	var lock sync.Mutex

	for i := 0; i < txSize; i++ {
		tx := randomTx()
		success[string(tx)] = 0
		txs = append(txs, tx)
	}

	var wg sync.WaitGroup
	for i := range txs {
		wg.Add(1)
		tx := txs[i]
		go func() {
			defer wg.Done()
			n := util.RandNum(len(tns))
			if err := tns[n].AddTx("test", tx); err == nil {
				lock.Lock()
				success[string(tx)]++
				lock.Unlock()
			}
		}()
	}
	wg.Wait()

	for i := range success {
		require.Equal(t, 1, success[i])
	}
	// then fetch blocks and compare
	compareBlocks(t)
}

// Note: We can not do restart test because the tendermint use the global variables so that we can just start or stop
// We have to start consensus service indepedent to make sure the variables is not same
// func TestSendTxWithNodeRestart(t *testing.T) {
// 	// Here we will stop a node and start it, then check if we get same blocks
// 	var wg sync.WaitGroup

// 	wg.Add(1)
// 	go func(t *testing.T) {
// 		defer wg.Done()
// 		for i := 0; i < 100; i++ {
// 			tx := randomTx()
// 			require.NoError(t, tns[0].AddTx("test", tx))
// 			time.Sleep(50 * time.Millisecond)
// 		}
// 	}(t)

// 	wg.Add(1)
// 	go func(t *testing.T) {
// 		defer wg.Done()
// 		require.NoError(t, tns[1].Stop())
// 		time.Sleep(2000 * time.Millisecond)
// 		require.NoError(t, tns[1].Start())
// 	}(t)

// 	wg.Wait()
// 	time.Sleep(2 * time.Second)
// 	// then we check if all node get same result
// 	txSize += 100
// 	compareBlocks(t)
// }

func TestStop(t *testing.T) {
	// time.Sleep(2 * time.Second)
	for i := range tns {
		require.NoError(t, tns[i].Stop())
	}
	require.NoError(t, os.RemoveAll(getTestPath()))
}

func compareBlocks(t *testing.T) {
	// then fetch blocks and compare
	var blocks = make(map[uint64]consensus.Block)
	for i := range tns {
		var txCount = make(map[string]int)
		var num uint64 = 1
		for {
			block, err := tns[i].GetBlock("test", num, true)
			require.NoError(t, err)
			if !util.Contain(blocks, block.GetNumber()) {
				blocks[block.GetNumber()] = block
			} else {
				require.Equal(t, blocks[block.GetNumber()].GetTxs(), block.GetTxs())
			}
			for _, tx := range block.GetTxs() {
				if !util.Contain(txCount, string(tx)) {
					txCount[string(tx)] = 0
				}
				txCount[string(tx)]++
			}
			if len(txCount) == txSize {
				break
			}
			num++
		}
	}
}

func getConfigPath(node int) string {
	return getNodePath(node) + "/orderer.yaml"
}

func getTestPath() string {
	gopath := os.Getenv("GOPATH")
	testPath, _ := util.MakeFileAbs("src/madledger/consensus/tendermint/.test", gopath)
	return testPath
}

func getNodePath(node int) string {
	return fmt.Sprintf("%s/orderers/%d", getTestPath(), node)
}

func randomTx() []byte {
	return []byte(util.RandomString(32))
}
