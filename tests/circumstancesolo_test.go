package tests

import (
	"encoding/hex"
	"madledger/common"
	"madledger/common/util"
	"madledger/core/types"
	oc "madledger/orderer/config"
	orderer "madledger/orderer/server"
	pc "madledger/peer/config"
	peer "madledger/peer/server"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	client "madledger/client/lib"
)

/*
* CircumstanceSolo begins from a empty environment and defines some operations as below.
* 1. Create channel test.
* 2. Create a contract.
* 3. Call the contract in different ways.
* 4. During main operates, there are some necessary query to make sure everything is ok.
 */

// Some consts of Balance
const (
	BalanceBin = "balance/Balance.bin"
	BalanceAbi = "balance/Balance.abi"
)

var (
	contractAddress common.Address
)

func TestInitCircumstanceSolo(t *testing.T) {
	err := initDir(".orderer")
	require.NoError(t, err)
	err = initDir(".peer")
	require.NoError(t, err)
	err = initDir(".client")
	require.NoError(t, err)
}

func TestCreateChannel(t *testing.T) {
	startSoloOrderer()
	startSoloPeer()
	client, err := getSoloClient()
	require.NoError(t, err)
	// first query channels
	// then query channels
	infos, err := client.ListChannel(true)
	require.NoError(t, err)
	channels := make([]string, 0)
	for _, info := range infos {
		channels = append(channels, info.Name)
	}
	require.Contains(t, channels, types.GLOBALCHANNELID)
	require.Contains(t, channels, types.CONFIGCHANNELID)
	require.NotContains(t, channels, "test")

	// then add a channel
	err = client.CreateChannel("test")
	require.NoError(t, err)
	// then query channels
	infos, err = client.ListChannel(true)
	require.NoError(t, err)
	channels = make([]string, 0)
	for _, info := range infos {
		channels = append(channels, info.Name)
	}
	require.Contains(t, channels, types.GLOBALCHANNELID)
	require.Contains(t, channels, types.CONFIGCHANNELID)
	require.Contains(t, channels, "test")

	// create channel test again
	err = client.CreateChannel("test")
	require.Error(t, err)
}

func TestCreateContract(t *testing.T) {
	client, err := getSoloClient()
	require.NoError(t, err)
	// then try to create a tx
	contractCodes, err := readCodes(BalanceBin)
	require.NoError(t, err)
	tx, err := types.NewTx("test", common.ZeroAddress, contractCodes, client.GetPrivKey())
	require.NoError(t, err)
	status, err := client.AddTx(tx)
	require.NoError(t, err)
	contractAddress = common.HexToAddress(status.ContractAddress)
}

func TestCallContract(t *testing.T) {
	client, err := getSoloClient()
	require.NoError(t, err)
	// then call the contract which is created before
	var payload []byte
	// 1. get
	payload, _ = hex.DecodeString("6d4ce63c")
	tx, _ := types.NewTx("test", contractAddress, payload, client.GetPrivKey())
	status, err := client.AddTx(tx)
	require.NoError(t, err)
	txStatus, err := getTxStatus(BalanceAbi, "get", status)
	assert.Equal(t, []string{"10"}, txStatus.Output)
	// 2. set 1314
	payload, _ = hex.DecodeString("60fe47b10000000000000000000000000000000000000000000000000000000000000522")
	tx, _ = types.NewTx("test", contractAddress, payload, client.GetPrivKey())
	status, err = client.AddTx(tx)
	require.NoError(t, err)
	txStatus, err = getTxStatus(BalanceAbi, "set", status)
	require.NoError(t, err)
	assert.Equal(t, []string{"true"}, txStatus.Output)
	// 3. get
	payload, _ = hex.DecodeString("6d4ce63c")
	tx, _ = types.NewTx("test", contractAddress, payload, client.GetPrivKey())
	status, err = client.AddTx(tx)
	require.NoError(t, err)
	txStatus, err = getTxStatus(BalanceAbi, "get", status)
	assert.Equal(t, []string{"1314"}, txStatus.Output)
}

func TestTxHistory(t *testing.T) {
	client, err := getSoloClient()
	require.NoError(t, err)
	// then get the history of the client
	address, _ := client.GetPrivKey().PubKey().Address()
	history, err := client.GetHistory(address.Bytes())
	require.NoError(t, err)
	// check channel test
	require.Contains(t, history.Txs, "test")
	require.Len(t, history.Txs["test"].Value, 4)
	// check config test
	// TODO: fix the bug that the history does not contain config channel
	// require.Contains(t, history.Txs, types.CONFIGCHANNELID)
}

func TestEnd(t *testing.T) {
	os.RemoveAll(".orderer")
	os.RemoveAll(".peer")
	os.RemoveAll(".client")
}

func startSoloOrderer() error {
	cfg, err := getSoloOrdererConfig()
	if err != nil {
		return err
	}
	server, err := orderer.NewServer(cfg)
	if err != nil {
		return err
	}
	go func() {
		server.Start()
	}()
	time.Sleep(300 * time.Millisecond)
	return nil
}

func startSoloPeer() error {
	server, err := peer.NewServer(getSoloPeerConfig())
	if err != nil {
		return err
	}
	go func() {
		server.Start()
	}()
	time.Sleep(300 * time.Millisecond)
	return nil
}

func getSoloOrdererConfig() (*oc.Config, error) {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/solo_orderer.yaml", gopath)
	cfg, err := oc.LoadConfig(cfgFilePath)
	if err != nil {
		return nil, err
	}
	chainPath, _ := util.MakeFileAbs("src/madledger/tests/.orderer/data/blocks", gopath)
	dbPath, _ := util.MakeFileAbs("src/madledger/tests/.orderer/data/leveldb", gopath)
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Dir = dbPath
	return cfg, nil
}

func getSoloPeerConfig() *pc.Config {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/solo_peer.yaml", gopath)
	cfg, _ := pc.LoadConfig(cfgFilePath)
	chainPath, _ := util.MakeFileAbs("src/madledger/tests/.peer/data/blocks", gopath)
	dbPath, _ := util.MakeFileAbs("src/madledger/tests/.peer/data/leveldb", gopath)
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Dir = dbPath
	return cfg
}

func getSoloClient() (*client.Client, error) {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/solo_client.yaml", gopath)
	c, err := client.NewClient(cfgFilePath)
	if err != nil {
		return nil, err
	}
	return c, nil
}
