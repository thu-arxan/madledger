// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package tests

import (
	"madledger/common"
	"madledger/common/util"
	"madledger/core"
	pc "madledger/peer/config"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
* CircumstanceAllSolo begins from a empty environment and defines some operations as below.
* It will test an environment which contains one orderer and one peer.
* 1. Create channel test.
* 2. Create a contract.
* 3. Call the contract in different ways.
* 4. During main operates, there are some necessary query to make sure everything is ok.
 */

var (
	allSoloContractAddress common.Address
)

func TestInitCircumstanceAllSolo(t *testing.T) {
	err := initDir(".orderer")
	require.NoError(t, err)
	err = initDir(".peer")
	require.NoError(t, err)
	err = initDir(".client")
	require.NoError(t, err)
	// then start necessary orderer and peer
	err = startSoloOrderer()
	require.NoError(t, err)
	err = startSoloPeer()
	require.NoError(t, err)
}

func TestAllSoloCreateChannel(t *testing.T) {
	client, err := getSoloClient()
	require.NoError(t, err)
	// get identity of solo peer
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/config/peer/solo_peer.yaml", gopath)
	cfg, _ := pc.LoadConfig(cfgFilePath)
	identity, err := cfg.GetIdentity()
	require.NoError(t, err)
	testCreateChannel(t, client, []*core.Member{identity})
}

func TestAllSoloCreateContract(t *testing.T) {
	client, err := getSoloClient()
	require.NoError(t, err)
	testCreateContract(t, client)
}

func TestAllSoloCallContract(t *testing.T) {
	client, err := getSoloClient()
	require.NoError(t, err)
	testCallContract(t, client)
}

func TestAllSoloTxHistory(t *testing.T) {
	client, err := getSoloClient()
	require.NoError(t, err)
	testTxHistory(t, client)
}

func TestAllSoloEnd(t *testing.T) {
	stopSoloOrderer()
	stopSoloPeer()
	os.RemoveAll(".orderer")
	os.RemoveAll(".peer")
	os.RemoveAll(".client")
}
