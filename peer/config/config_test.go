package config

import (
	"madledger/common/util"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig(getTestConfigFilePath())
	if err != nil {
		t.Fatal(err)
	}
	require.True(t, cfg.Debug)
	require.False(t, cfg.TLS.Enable)
}

func TestGetIdentity(t *testing.T) {
	cfg, _ := LoadConfig(getTestConfigFilePath())
	_, err := cfg.GetIdentity()
	require.NoError(t, err)
	// set key to nil
	cfg.KeyStore.Key = ""
	_, err = cfg.GetIdentity()
	require.Error(t, err, "The key should not be nil")
}

func getTestConfigFilePath() string {
	gopath := os.Getenv("GOPATH")
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/peer/config/.config.yaml", gopath)
	return cfgFilePath
}
