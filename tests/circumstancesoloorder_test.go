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
	"madledger/core"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
* CircumstanceSoloOrderer begins from a empty environment and defines some operations as below.
* The circumstance includes one orderer and three peers.
* 1. Create channel test.
* 2. Create a contract.
* 3. Call the contract in different ways.
* 4. During main operates, there are some necessary query to make sure everything is ok.
 */

func TestInitCircumstanceSoloOrderer(t *testing.T) {
	err := initDir(".orderer")
	require.NoError(t, err)
	err = initDir(".peer0")
	require.NoError(t, err)
	err = initDir(".peer1")
	require.NoError(t, err)
	err = initDir(".peer2")
	require.NoError(t, err)
	err = initDir(".client")
	require.NoError(t, err)
	// then start necessary orderer and peer
	err = startSoloOrderer()
	require.NoError(t, err)
	err = startPeers(3)
	require.NoError(t, err)
}

func TestSoloOrdererCreateChannel(t *testing.T) {
	client, err := getSoloClient()
	require.NoError(t, err)
	var identities []*core.Member
	for i := 0; i < 3; i++ {
		cfg := getPeerConfig(i)
		identity, err := cfg.GetIdentity()
		require.NoError(t, err)
		identities = append(identities, identity)
	}
	testCreateChannel(t, client, identities)
}

func TestSoloOrdererCreateContract(t *testing.T) {
	client, err := getSoloClient()
	require.NoError(t, err)
	testCreateContract(t, client)
}

func TestSoloOrdererCallContract(t *testing.T) {
	client, err := getSoloClient()
	require.NoError(t, err)
	testCallContract(t, client)
}

func TestSoloOrdererTxHistory(t *testing.T) {
	client, err := getSoloClient()
	require.NoError(t, err)
	testTxHistory(t, client)
}

// func TestSoloOrdererAsset(t *testing.T) {
// 	client, err := getSoloClient()
// 	require.NoError(t, err)
// 	testAsset(t, client)
// }
func TestSoloOrdererEnd(t *testing.T) {
	stopSoloOrderer()
	stopPeers(3)
	os.RemoveAll(".orderer")
	os.RemoveAll(".peer0")
	os.RemoveAll(".peer1")
	os.RemoveAll(".peer2")
	os.RemoveAll(".client")
}
