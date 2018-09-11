package config

import (
	"madledger/common/util"
	"os"
	"testing"
)

var (
	cfg *Config
)

func TestLoadConfig(t *testing.T) {
	var err error
	cfg, err = LoadConfig(getTestConfigFilePath())
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetKeyStoreConfig(t *testing.T) {
	_, err := cfg.GetKeyStoreConfig()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetOrdererConfig(t *testing.T) {
	ordererConfig, err := cfg.GetOrdererConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(ordererConfig.Address) != 1 || ordererConfig.Address[0] != "localhost:12345" {
		t.Fatal()
	}
}

func TestGetPeerConfig(t *testing.T) {
	peerConfig, err := cfg.GetPeerConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(peerConfig.Address) != 1 || peerConfig.Address[0] != "localhost:23456" {
		t.Fatal()
	}
}

func getTestConfigFilePath() string {
	gopath := os.Getenv("GOPATH")
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/client/config/.config.yaml", gopath)
	return cfgFilePath
}
