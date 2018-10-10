package config

import (
	"madledger/common/util"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	cfg *Config
)

func TestLoadConfig(t *testing.T) {
	var err error
	cfg, err = LoadConfig(getTestConfigFilePath())
	require.NoError(t, err)
}

func TestGetKeyStoreConfig(t *testing.T) {
	_, err := cfg.GetKeyStoreConfig()
	require.NoError(t, err)
}

func TestGetOrdererConfig(t *testing.T) {
	ordererConfig, err := cfg.GetOrdererConfig()
	require.NoError(t, err)
	require.Equal(t, []string{"localhost:12345"}, ordererConfig.Address)
}

func TestGetPeerConfig(t *testing.T) {
	peerConfig, err := cfg.GetPeerConfig()
	require.NoError(t, err)
	require.Equal(t, []string{"localhost:23456", "localhost:34567"}, peerConfig.Address)
}

func getTestConfigFilePath() string {
	gopath := os.Getenv("GOPATH")
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/client/config/.config.yaml", gopath)
	return cfgFilePath
}
