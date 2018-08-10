package config

import (
	"errors"
	"fmt"
	"madledger/util"
	"os"
	"testing"
)

func TestGetServerConfig(t *testing.T) {
	cfg, err := LoadConfig(getTestConfigFilePath())
	if err != nil {
		t.Fatal(err)
	}
	serverCfg, err := cfg.GetServerConfig()
	if err != nil {
		t.Fatal(err)
	}
	if serverCfg.Port != 12345 {
		t.Fatal(fmt.Errorf("The port is %d", serverCfg.Port))
	}
	if serverCfg.Address != "localhost" {
		t.Fatal(fmt.Errorf("The address is %s", serverCfg.Address))
	}
	if serverCfg.Debug != true {
		t.Fatal(fmt.Errorf("The Debug is %t", serverCfg.Debug))
	}
	// then change the value of cfg
	// check address
	cfg.Address = ""
	_, err = cfg.GetServerConfig()
	if err.Error() != "The address can not be empty" {
		t.Fatal(err)
	}
	// check port
	cfg.Address = "localhost"
	cfg.Port = -1
	_, err = cfg.GetServerConfig()
	if err.Error() != "The port can not be -1" {
		t.Fatal(err)
	}
	cfg.Port = -1
	_, err = cfg.GetServerConfig()
	if err == nil || err.Error() != "The port can not be -1" {
		t.Fatal(err)
	}
	cfg.Port = 1023
	_, err = cfg.GetServerConfig()
	if err == nil || err.Error() != "The port can not be 1023" {
		t.Fatal(err)
	}
	cfg.Port = 1024
	_, err = cfg.GetServerConfig()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetBlockChainConfig(t *testing.T) {
	cfg, err := LoadConfig(getTestConfigFilePath())
	if err != nil {
		t.Fatal(err)
	}
	chainCfg, err := cfg.GetBlockChainConfig()
	if err != nil {
		t.Fatal(err)
	}
	if chainCfg.BatchTimeout != 1000 {
		t.Fatal(fmt.Errorf("The batch timeout is %d", chainCfg.BatchTimeout))
	}
	if chainCfg.BatchSize != 100 {
		t.Fatal(fmt.Errorf("The batch size is %d", chainCfg.BatchSize))
	}
	if chainCfg.Path == "" {
		t.Fatal(errors.New("The path is chain config is empty"))
	}
	// then change the value of cfg
	// check batch timeout
	cfg.BlockChain.BatchTimeout = 0
	_, err = cfg.GetBlockChainConfig()
	if err == nil || err.Error() != "The batch timeout can not be 0" {
		t.Fatal(err)
	}
	// check batch size
	cfg.BlockChain.BatchTimeout = 1000
	cfg.BlockChain.BatchSize = -1
	_, err = cfg.GetBlockChainConfig()
	if err == nil || err.Error() != "The batch size can not be -1" {
		t.Fatal(err)
	}
	// check debug
	cfg.BlockChain.BatchSize = 100
	cfg.Debug = false
	_, err = cfg.GetBlockChainConfig()
	if err == nil || err.Error() != "The path of blockchain is not provided" {
		t.Fatal(err)
	}
}

func TestGetConsensusConfig(t *testing.T) {
	cfg, err := LoadConfig(getTestConfigFilePath())
	if err != nil {
		t.Fatal(err)
	}
	consensusCfg, err := cfg.GetConsensusConfig()
	if err != nil {
		t.Fatal(err)
	}
	if consensusCfg.Type != SOLO {
		t.Fatal(fmt.Errorf("The type of consensus if %d", consensusCfg.Type))
	}
	cfg.Consensus.Type = "raft"
	consensusCfg, err = cfg.GetConsensusConfig()
	if err == nil || err.Error() != "Raft is not supported yet" {
		t.Fatal(errors.New("Should be 'Raft is not supported yet' error"))
	}
	cfg.Consensus.Type = "unknown"
	consensusCfg, err = cfg.GetConsensusConfig()
	if err == nil || err.Error() != "Unsupport consensus type: unknown" {
		t.Fatal(errors.New("Should be 'Unsupport consensus type: unknown' error"))
	}
}

func getTestConfigFilePath() string {
	gopath := os.Getenv("GOPATH")
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/orderer/config/.config.yaml", gopath)
	return cfgFilePath
}
