// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package config

import (
	"madledger/common/util"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig(getTestConfigFilePath())
	require.NoError(t, err)

	require.Equal(t, cfg.Port, 12345)
	require.Equal(t, cfg.Address, "localhost")
	require.Equal(t, cfg.Debug, true)
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
	// set some thing wrong
	cfg.Consensus.Tendermint.ID = "fatal id"
	_, err = cfg.GetConsensusConfig()
	require.EqualError(t, err, "The ID(fatal id) of tendermint is not legal")

	cfg.Consensus.Type = "raft"
	consensusCfg, err = cfg.GetConsensusConfig()
	require.EqualError(t, err, "Raft id should not be zero")

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
