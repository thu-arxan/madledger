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
