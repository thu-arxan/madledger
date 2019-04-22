package tendermint

import (
	"fmt"
	"madledger/common/util"
	"madledger/consensus"
	"madledger/core/types"
	"madledger/orderer/config"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// This file will start some tendermint nodes and run basic tests

var (
	// tendermint nodes
	tns [4]consensus.Consensus
)

func TestMain(m *testing.M) {
	ret := m.Run()
	os.Exit(ret)
}

func TestStart(t *testing.T) {
	var channels = make(map[string]consensus.Config)
	channels[types.GLOBALCHANNELID] = consensus.DefaultConfig()
	channels[types.CONFIGCHANNELID] = consensus.DefaultConfig()
	for i := 0; i < len(tns); i++ {
		cfg, err := config.LoadConfig(getConfigPath(i))
		require.NoError(t, err)
		nodePath := getNodePath(i)
		cfg.BlockChain.Path = nodePath + "/" + cfg.BlockChain.Path
		cfg.Consensus.Tendermint.Path = nodePath + "/" + cfg.Consensus.Tendermint.Path
		cfg.DB.LevelDB.Path = nodePath + "/" + cfg.DB.LevelDB.Path

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
	time.Sleep(2 * time.Second)
}

func TestSendTx(t *testing.T) {
	var txSize = 128
	var txs [][]byte
	var success = make(map[string]int)
	var lock sync.Mutex

	for i := 0; i < txSize; i++ {
		tx := randomTx()
		success[string(tx)] = 0
		txs = append(txs, tx)
	}

	var wg sync.WaitGroup
	// each tx send 3 times
	for i := range txs {
		for m := 0; m < 3; m++ {
			wg.Add(1)
			tx := txs[i]
			go func() {
				defer wg.Done()
				if err := tns[0].AddTx("test", tx); err == nil {
					lock.Lock()
					success[string(tx)]++
					lock.Unlock()
				}
			}()
		}
	}
	wg.Wait()

	for i := range success {
		require.Equal(t, 1, success[i])
	}
}

func TestClose(t *testing.T) {

}

func getConfigPath(node int) string {
	return getNodePath(node) + "/orderer.yaml"
}

func getNodePath(node int) string {
	gopath := os.Getenv("GOPATH")
	nodePath, _ := util.MakeFileAbs(fmt.Sprintf("src/madledger/consensus/tendermint/.test/env/orderers/%d", node), gopath)
	return nodePath
}

func randomTx() []byte {
	return []byte(util.RandomString(32))
}
