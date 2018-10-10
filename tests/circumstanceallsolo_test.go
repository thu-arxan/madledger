package tests

import (
	"madledger/common"
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
	startSoloOrderer()
	startSoloPeer()
}

func TestAllSoloCreateChannel(t *testing.T) {
	client, err := getSoloClient()
	require.NoError(t, err)
	testCreateChannel(t, client)
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
