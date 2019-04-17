package config

import (
	"fmt"
	"madledger/common/util"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetServerConfig(t *testing.T) {
	cfg, err := LoadConfig(getTestConfigFilePath())
	if err != nil {
		t.Fatal(err)
	}
	serverCfg, err := cfg.GetServerConfig()
	require.NoError(t, err)

	require.Equal(t, serverCfg.Port, 23456)
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
	require.NotEqual(t, chainCfg.Path, "", "The path is chain config is empty")
}

func TestGetOrdererConfig(t *testing.T) {
	cfg, err := LoadConfig(getTestConfigFilePath())
	require.NoError(t, err)

	ordererCfg, err := cfg.GetOrdererConfig()
	require.NoError(t, err)
	if !reflect.DeepEqual(ordererCfg.Address, []string{"localhost:12345"}) {
		t.Fatal(fmt.Errorf("Address is %s", ordererCfg.Address))
	}
}

func TestGetDBConfig(t *testing.T) {
	cfg, err := LoadConfig(getTestConfigFilePath())
	require.NoError(t, err)

	dbCfg, err := cfg.GetDBConfig()
	require.NoError(t, err)
	require.Equal(t, dbCfg.Type, LEVELDB)
	require.NotEqual(t, dbCfg.LevelDB.Dir, "")

	cfg.DB.Type = "unknown"
	dbCfg, err = cfg.GetDBConfig()
	require.Error(t, err, "Unsupport db type: unknown")
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
