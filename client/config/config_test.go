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
