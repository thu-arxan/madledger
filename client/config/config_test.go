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
	require.Equal(t, []string{"localhost:12345"}, cfg.Orderer.Address)
	require.Equal(t, []string{"localhost:23456", "localhost:34567"}, cfg.Peer.Address)
	require.Len(t, cfg.KeyStore.Privs, 1)
	require.False(t, cfg.TLS.Enable)
}

func getTestConfigFilePath() string {
	gopath := os.Getenv("GOPATH")
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/client/config/.config.yaml", gopath)
	return cfgFilePath
}
