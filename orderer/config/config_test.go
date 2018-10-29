package config

import (
	"madledger/common/util"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetServerConfig(t *testing.T) {
	cfg, err := LoadConfig(getTestConfigFilePath())
	require.NoError(t, err)

	serverCfg, err := cfg.GetServerConfig()
	require.NoError(t, err)
	require.Equal(t, serverCfg.Port, 12345)
	require.Equal(t, serverCfg.Address, "localhost")
	require.Equal(t, serverCfg.Debug, true)

	// then change the value of cfg
	// check address
	cfg.Address = ""
	_, err = cfg.GetServerConfig()
	require.EqualError(t, err, "The address can not be empty")

	// check port
	cfg.Address = "localhost"
	cfg.Port = -1
	_, err = cfg.GetServerConfig()
	require.EqualError(t, err, "The port can not be -1")

	cfg.Port = 1023
	_, err = cfg.GetServerConfig()
	require.EqualError(t, err, "The port can not be 1023")

	cfg.Port = 1024
	_, err = cfg.GetServerConfig()
	require.NoError(t, err)
}

func TestGetBlockChainConfig(t *testing.T) {
	cfg, err := LoadConfig(getTestConfigFilePath())
	require.NoError(t, err)

	chainCfg, err := cfg.GetBlockChainConfig()
	require.NoError(t, err)
	require.Equal(t, chainCfg.BatchTimeout, 1000)
	require.Equal(t, chainCfg.BatchSize, 100)
	require.NotEqual(t, chainCfg.Path, "")
	// then change the value of cfg
	// check batch timeout
	cfg.BlockChain.BatchTimeout = 0
	_, err = cfg.GetBlockChainConfig()
	require.EqualError(t, err, "The batch timeout can not be 0")
	// check batch size
	cfg.BlockChain.BatchTimeout = 1000
	cfg.BlockChain.BatchSize = -1
	_, err = cfg.GetBlockChainConfig()
	require.EqualError(t, err, "The batch size can not be -1")
}

func TestGetConsensusConfig(t *testing.T) {
	cfg, err := LoadConfig(getTestConfigFilePath())
	require.NoError(t, err)

	consensusCfg, err := cfg.GetConsensusConfig()
	require.NoError(t, err)
	require.Equal(t, consensusCfg.Type, SOLO)

	// then check bft
	cfg.Consensus.Type = "bft"
	consensusCfg, err = cfg.GetConsensusConfig()
	require.NoError(t, err)
	require.NotEqual(t, "", consensusCfg.BFT.Path)
	require.Equal(t, 26656, consensusCfg.BFT.Port.P2P)
	require.Equal(t, 26657, consensusCfg.BFT.Port.RPC)
	require.Equal(t, 26658, consensusCfg.BFT.Port.APP)
	require.Equal(t, []string{""}, consensusCfg.BFT.P2PAddress)

	cfg.Consensus.Type = "raft"
	consensusCfg, err = cfg.GetConsensusConfig()
	require.EqualError(t, err, "Raft is not supported yet")

	cfg.Consensus.Type = "unknown"
	consensusCfg, err = cfg.GetConsensusConfig()
	require.EqualError(t, err, "Unsupport consensus type: unknown")
}

func TestGetDBConfig(t *testing.T) {
	cfg, err := LoadConfig(getTestConfigFilePath())
	require.NoError(t, err)

	dbCfg, err := cfg.GetDBConfig()
	require.NoError(t, err)
	require.Equal(t, dbCfg.Type, LEVELDB)
	require.NotEqual(t, dbCfg.LevelDB.Path, "")

	cfg.DB.Type = "unknown"
	dbCfg, err = cfg.GetDBConfig()
	require.EqualError(t, err, "Unsupport db type: unknown")
}

func getTestConfigFilePath() string {
	gopath := os.Getenv("GOPATH")
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/orderer/config/.orderer.yaml", gopath)
	return cfgFilePath
}
