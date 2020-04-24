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
	"encoding/json"
	"fmt"
	"madledger/blockchain/asset"
	cc "madledger/blockchain/config"
	client "madledger/client/lib"
	"madledger/common"
	"madledger/common/abi"
	"madledger/common/crypto"
	"madledger/core"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	contractAddress = make(map[string]common.Address)
)

func testCreateChannel(t *testing.T, client *client.Client, peers []*core.Member) {
	// first query channels
	// then query channels
	infos, err := client.ListChannel(true)
	require.NoError(t, err)
	channels := make([]string, 0)
	for _, info := range infos {
		channels = append(channels, info.Name)
	}
	require.Contains(t, channels, core.GLOBALCHANNELID)
	require.Contains(t, channels, core.CONFIGCHANNELID)
	require.Contains(t, channels, core.ASSETCHANNELID)
	require.NotContains(t, channels, "public")

	// then add a channel
	err = client.CreateChannel("public", true, nil, nil, 0, 1, 10000000, []string{"localhost:23456"})
	require.NoError(t, err)

	// then query channels
	infos, err = client.ListChannel(true)
	require.NoError(t, err)
	channels = make([]string, 0)
	for _, info := range infos {
		channels = append(channels, info.Name)
	}
	require.Contains(t, channels, core.GLOBALCHANNELID)
	require.Contains(t, channels, core.CONFIGCHANNELID)
	require.Contains(t, channels, core.ASSETCHANNELID)
	require.Contains(t, channels, "public")
	// create channel test again
	err = client.CreateChannel("public", true, nil, nil, 0, 1, 10000000, []string{"localhost:23456"})
	require.Error(t, err)
	// create private channel
	err = client.CreateChannel("private", false, nil, peers, 0, 1, 10000000, []string{"localhost:23456"})
	require.NoError(t, err)
}

func testCreateChannelByHTTP(t *testing.T, client *client.HTTPClient, peers []*core.Member) {
	// first query channels
	// then query channels
	infos, err := client.ListChannelByHTTP(true)
	require.NoError(t, err)
	channels := make([]string, 0)
	for _, info := range infos {
		channels = append(channels, info.Name)
	}
	require.Contains(t, channels, core.GLOBALCHANNELID)
	require.Contains(t, channels, core.CONFIGCHANNELID)
	require.Contains(t, channels, core.ASSETCHANNELID)
	require.NotContains(t, channels, "public")

	// then add a channel
	err = client.CreateChannelByHTTP("public", true, nil, nil, 0, 1, 10000000)
	require.NoError(t, err)
	// then query channels
	infos, err = client.ListChannelByHTTP(true)
	require.NoError(t, err)
	channels = make([]string, 0)
	for _, info := range infos {
		channels = append(channels, info.Name)
	}
	require.Contains(t, channels, core.GLOBALCHANNELID)
	require.Contains(t, channels, core.CONFIGCHANNELID)
	require.Contains(t, channels, core.ASSETCHANNELID)
	require.Contains(t, channels, "public")
	// create channel test again
	err = client.CreateChannelByHTTP("public", true, nil, nil, 0, 1, 10000000)
	require.Error(t, err)
	// create private channel
	err = client.CreateChannelByHTTP("private", false, nil, peers, 0, 1, 10000000)
	require.NoError(t, err)
}

func testCreateContract(t *testing.T, client *client.Client) {
	// First, test on channel public, which is public
	createContract(t, "public", client)
	// Then, test on channel private, which is private
	createContract(t, "private", client)
}
func testCreateContractByHTTP(t *testing.T, client *client.HTTPClient) {
	// First, test on channel public, which is public
	createContractByHTTP(t, "public", client)
	// Then, test on channel private, which is private
	createContractByHTTP(t, "private", client)
}

// createContract is the detail implemtation of testCreateContract
func createContract(t *testing.T, channelID string, client *client.Client) {
	contractCodes, err := readCodes(BalanceBin)
	require.NoError(t, err)

	tx, err := core.NewTx(channelID, common.ZeroAddress, contractCodes, 0, "", client.GetPrivKey())
	require.NoError(t, err)

	status, err := client.AddTx(tx)
	require.NoError(t, err)
	require.Empty(t, status.Err)
	contractAddress[channelID] = common.HexToAddress(status.ContractAddress)
	// then try to create the same contract again, this should be a duplicate error
	tx, err = core.NewTx(channelID, common.ZeroAddress, contractCodes, 0, "", client.GetPrivKey())
	require.NoError(t, err)

	status, err = client.AddTx(tx)
	require.NoError(t, err)
	// require.Equal(t, status.Err, "Duplicate address")
	require.Equal(t, status.Err, "")
}

// createContract is the detail implemtation of testCreateContract
func createContractByHTTP(t *testing.T, channelID string, client *client.HTTPClient) {
	contractCodes, err := readCodes(BalanceBin)
	require.NoError(t, err)

	tx, err := core.NewTx(channelID, common.ZeroAddress, contractCodes, 0, "", client.GetPrivKey())
	require.NoError(t, err)

	status, err := client.AddTxByHTTP(tx)
	require.NoError(t, err)
	require.Empty(t, status.Err)
	contractAddress[channelID] = common.HexToAddress(status.ContractAddress)
	// then try to create the same contract again, this should be a duplicate error
	tx, err = core.NewTx(channelID, common.ZeroAddress, contractCodes, 0, "", client.GetPrivKey())
	require.NoError(t, err)

	status, err = client.AddTxByHTTP(tx)
	require.NoError(t, err)
	require.Equal(t, status.Err, "")
}
func testCallContract(t *testing.T, client *client.Client) {
	callContract(t, "public", client)
	callContract(t, "private", client)
}
func testCallContractByHTTP(t *testing.T, client *client.HTTPClient) {
	callContractByHTTP(t, "public", client)
	callContractByHTTP(t, "private", client)
}
func callContract(t *testing.T, channelID string, client *client.Client) {
	var contractAddress = contractAddress[channelID]
	// then call the contract which is created before
	var payload []byte
	// 1. get
	payload, _ = abi.Pack(BalanceAbi, "get")
	tx, _ := core.NewTx(channelID, contractAddress, payload, 0, "", client.GetPrivKey())
	status, err := client.AddTx(tx)
	require.NoError(t, err)
	txStatus, err := getTxStatus(BalanceAbi, "get", status)
	require.NoError(t, err)
	assert.Equal(t, []string{"10"}, txStatus.Output)
	// 2. set 1314
	payload, _ = abi.Pack(BalanceAbi, "set", "1314")
	tx, _ = core.NewTx(channelID, contractAddress, payload, 0, "", client.GetPrivKey())
	status, err = client.AddTx(tx)
	require.NoError(t, err)
	txStatus, err = getTxStatus(BalanceAbi, "set", status)
	require.NoError(t, err)
	assert.Equal(t, []string{"true"}, txStatus.Output)
	// 3. get
	payload, _ = abi.Pack(BalanceAbi, "get")
	tx, _ = core.NewTx(channelID, contractAddress, payload, 0, "", client.GetPrivKey())
	status, err = client.AddTx(tx)
	require.NoError(t, err)
	txStatus, err = getTxStatus(BalanceAbi, "get", status)
	assert.Equal(t, []string{"1314"}, txStatus.Output)
	// 4. sub
	payload, _ = abi.Pack(BalanceAbi, "sub", []string{"794"}...)
	tx, _ = core.NewTx(channelID, contractAddress, payload, 0, "", client.GetPrivKey())
	status, err = client.AddTx(tx)
	require.NoError(t, err)
	txStatus, err = getTxStatus(BalanceAbi, "sub", status)
	assert.Equal(t, []string{"520"}, txStatus.Output)
	// 5. add
	payload, _ = abi.Pack(BalanceAbi, "add", []string{"794"}...)
	tx, _ = core.NewTx(channelID, contractAddress, payload, 0, "", client.GetPrivKey())
	status, err = client.AddTx(tx)
	require.NoError(t, err)
	txStatus, err = getTxStatus(BalanceAbi, "add", status)
	assert.Equal(t, []string{"1314"}, txStatus.Output)
	// 6. info
	payload, _ = abi.Pack(BalanceAbi, "info")
	tx, _ = core.NewTx(channelID, contractAddress, payload, 0, "", client.GetPrivKey())
	status, err = client.AddTx(tx)
	require.NoError(t, err)
	txStatus, err = getTxStatus(BalanceAbi, "info", status)
	address, err := client.GetPrivKey().PubKey().Address()
	require.NoError(t, err)
	assert.Equal(t, []string{address.String(), "1314"}, txStatus.Output)
	// then call an address which is not exist
	invalidAddress := common.HexToAddress("0x829f6d8cc2a094b5b1d9e2c4e14e38bbb0ee1400")
	tx, _ = core.NewTx(channelID, invalidAddress, []byte("invalid"), 0, "", client.GetPrivKey())
	status, err = client.AddTx(tx)
	require.NoError(t, err)
	require.Equal(t, "Invalid Address", status.Err)
}
func callContractByHTTP(t *testing.T, channelID string, client *client.HTTPClient) {
	var contractAddress = contractAddress[channelID]
	// then call the contract which is created before
	var payload []byte
	// 1. get
	payload, _ = abi.Pack(BalanceAbi, "get")
	tx, _ := core.NewTx(channelID, contractAddress, payload, 0, "", client.GetPrivKey())
	status, err := client.AddTxByHTTP(tx)
	require.NoError(t, err)
	txStatus, err := getTxStatus(BalanceAbi, "get", status)
	require.NoError(t, err)
	assert.Equal(t, []string{"10"}, txStatus.Output)
	// 2. set 1314
	payload, _ = abi.Pack(BalanceAbi, "set", "1314")
	tx, _ = core.NewTx(channelID, contractAddress, payload, 0, "", client.GetPrivKey())
	status, err = client.AddTxByHTTP(tx)
	require.NoError(t, err)
	txStatus, err = getTxStatus(BalanceAbi, "set", status)
	require.NoError(t, err)
	assert.Equal(t, []string{"true"}, txStatus.Output)
	// 3. get
	payload, _ = abi.Pack(BalanceAbi, "get")
	tx, _ = core.NewTx(channelID, contractAddress, payload, 0, "", client.GetPrivKey())
	status, err = client.AddTxByHTTP(tx)
	require.NoError(t, err)
	txStatus, err = getTxStatus(BalanceAbi, "get", status)
	assert.Equal(t, []string{"1314"}, txStatus.Output)
	// 4. sub
	payload, _ = abi.Pack(BalanceAbi, "sub", []string{"794"}...)
	tx, _ = core.NewTx(channelID, contractAddress, payload, 0, "", client.GetPrivKey())
	status, err = client.AddTxByHTTP(tx)
	require.NoError(t, err)
	txStatus, err = getTxStatus(BalanceAbi, "sub", status)
	assert.Equal(t, []string{"520"}, txStatus.Output)
	// 5. add
	payload, _ = abi.Pack(BalanceAbi, "add", []string{"794"}...)
	tx, _ = core.NewTx(channelID, contractAddress, payload, 0, "", client.GetPrivKey())
	status, err = client.AddTxByHTTP(tx)
	require.NoError(t, err)
	txStatus, err = getTxStatus(BalanceAbi, "add", status)
	assert.Equal(t, []string{"1314"}, txStatus.Output)
	// 6. info
	payload, _ = abi.Pack(BalanceAbi, "info")
	tx, _ = core.NewTx(channelID, contractAddress, payload, 0, "", client.GetPrivKey())
	status, err = client.AddTxByHTTP(tx)
	require.NoError(t, err)
	txStatus, err = getTxStatus(BalanceAbi, "info", status)
	address, err := client.GetPrivKey().PubKey().Address()
	require.NoError(t, err)
	assert.Equal(t, []string{address.String(), "1314"}, txStatus.Output)
	// then call an address which is not exist
	invalidAddress := common.HexToAddress("0x829f6d8cc2a094b5b1d9e2c4e14e38bbb0ee1400")
	tx, _ = core.NewTx(channelID, invalidAddress, []byte("invalid"), 0, "", client.GetPrivKey())
	status, err = client.AddTxByHTTP(tx)
	require.NoError(t, err)
	require.Equal(t, "Invalid Address", status.Err)
}
func testTxHistory(t *testing.T, client *client.Client) {
	// then get the history of the client
	address, _ := client.GetPrivKey().PubKey().Address()
	history, err := client.GetHistory(address.Bytes())
	require.NoError(t, err)
	// check channel public
	require.Contains(t, history.Txs, "public")
	require.Len(t, history.Txs["public"].Value, 9)
	// check channel private
	require.Contains(t, history.Txs, "private")
	require.Len(t, history.Txs["private"].Value, 9)
	// check cahnnel config
	require.Contains(t, history.Txs, core.CONFIGCHANNELID)
	require.Len(t, history.Txs[core.CONFIGCHANNELID].Value, 2)
}
func testTxHistoryByHTTP(t *testing.T, client *client.HTTPClient) {
	// then get the history of the client
	address, _ := client.GetPrivKey().PubKey().Address()
	history, err := client.GetHistoryByHTTP(address.Bytes())
	require.NoError(t, err)
	// check channel public
	require.Contains(t, history.Txs, "public")
	require.Len(t, history.Txs["public"].Value, 9)
	// check channel private
	require.Contains(t, history.Txs, "private")
	require.Len(t, history.Txs["private"].Value, 9)
	// check cahnnel config
	require.Contains(t, history.Txs, core.CONFIGCHANNELID)
	require.Len(t, history.Txs[core.CONFIGCHANNELID].Value, 2)
}
func testAsset(t *testing.T, client *client.Client, peers []string) {
	algo := crypto.KeyAlgoSecp256k1

	issuerKey, err := crypto.GeneratePrivateKey(algo)
	require.NoError(t, err)
	falseIssuerKey, err := crypto.GeneratePrivateKey(algo)
	require.NoError(t, err)
	require.NotEqual(t, issuerKey, falseIssuerKey)
	receiverKey, err := crypto.GeneratePrivateKey(algo)
	require.NoError(t, err)

	issuer, err := issuerKey.PubKey().Address()
	require.NoError(t, err)
	falseIssuer, err := falseIssuerKey.PubKey().Address()
	require.NoError(t, err)
	receiver, err := receiverKey.PubKey().Address()
	require.NoError(t, err)

	err = client.CreateChannel("test", true, nil, nil, 0, 1, 10000000, peers)
	require.NoError(t, err)

	//issue to issuer itself
	coreTx := getAssetChannelTx(core.IssueContractAddress, issuer, "", uint64(10), issuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)

	balance, err := client.GetAccountBalance(issuer)
	require.NoError(t, err)
	require.Equal(t, uint64(10), balance)

	//falseissuer issue fail
	coreTx = getAssetChannelTx(core.IssueContractAddress, falseIssuer, "", uint64(10), falseIssuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)

	balance, err = client.GetAccountBalance(falseIssuer)
	require.NoError(t, err)
	require.Equal(t, uint64(0), balance)

	//test issue to channel
	coreTx = getAssetChannelTx(core.IssueContractAddress, common.ZeroAddress, "test", uint64(10), issuerKey)
	// question, what if test is not created?
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)
	balance, err = client.GetAccountBalance(common.AddressFromChannelID("test"))
	require.NoError(t, err)
	require.Equal(t, uint64(10), balance)

	//test transfer
	coreTx = getAssetChannelTx(core.TransferContractrAddress, receiver, "", uint64(5), issuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)
	balance, err = client.GetAccountBalance(receiver)
	require.NoError(t, err)
	require.Equal(t, uint64(5), balance)

	//test transfer fail
	coreTx = getAssetChannelTx(core.TransferContractrAddress, receiver, "", uint64(5), falseIssuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)
	balance, err = client.GetAccountBalance(receiver)
	require.NoError(t, err)
	require.Equal(t, uint64(5), balance)

	//test transfer to oneself
	coreTx = getAssetChannelTx(core.TransferContractrAddress, receiver, "", uint64(5), receiverKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)
	balance, err = client.GetAccountBalance(receiver)
	require.NoError(t, err)
	require.Equal(t, uint64(5), balance)

	//4.test exchangeToken a.k.a transfer to channel in orderer execution
	coreTx = getAssetChannelTx(core.TokenExchangeAddress, common.ZeroAddress, "test", uint64(5), receiverKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)

	balance, err = client.GetAccountBalance(common.AddressFromChannelID("test"))
	require.NoError(t, err)
	require.Equal(t, uint64(15), balance)

	token, err := client.GetTokenInfo(receiver, []byte("test"))
	require.NoError(t, err)
	require.Equal(t, uint64(5), token)

	//test Block Price
	coreTx, err = core.NewTx("test", common.ZeroAddress, []byte("success"), 0, "", issuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)

	//change BlockPrice of test channel's

	payload, err := json.Marshal(cc.Payload{
		ChannelID: "test",
		Profile: &cc.Profile{
			BlockPrice: 100,
		},
	})
	require.NoError(t, err)
	coreTx, err = core.NewTx(core.CONFIGCHANNELID, common.ZeroAddress, payload, 0, "", issuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)

	//now add tx that cause due
	coreTx, err = core.NewTx("test", common.ZeroAddress, []byte("cause due but pass"), 0, "", issuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)

	// now add multiple txs to ensure that orderers have executed prev tx and stopped receiving tx
	for i := 0; i < 10; i++ {
		coreTx, err = core.NewTx("test", common.ZeroAddress, []byte("multiple tx"), 0, fmt.Sprintln(i), issuerKey)
		_, _ = client.AddTx(coreTx)
	}

	//this one should fail
	coreTx, err = core.NewTx("test", common.ZeroAddress, []byte("fail"), 0, "", issuerKey)
	_, err = client.AddTx(coreTx)
	require.Error(t, err)

	//now issue money to channel account to wake it
	coreTx = getAssetChannelTx(core.IssueContractAddress, common.ZeroAddress, "test", uint64(1000000), issuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)

	coreTx, err = core.NewTx("test", common.ZeroAddress, []byte("success again"), 0, "", issuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)
}

func testAssetByHTTP(t *testing.T, client *client.HTTPClient) {
	algo := crypto.KeyAlgoSecp256k1

	issuerKey, err := crypto.GeneratePrivateKey(algo)
	require.NoError(t, err)
	falseIssuerKey, err := crypto.GeneratePrivateKey(algo)
	require.NoError(t, err)
	require.NotEqual(t, issuerKey, falseIssuerKey)
	receiverKey, err := crypto.GeneratePrivateKey(algo)
	require.NoError(t, err)

	issuer, err := issuerKey.PubKey().Address()
	require.NoError(t, err)
	falseIssuer, err := falseIssuerKey.PubKey().Address()
	require.NoError(t, err)
	receiver, err := receiverKey.PubKey().Address()
	require.NoError(t, err)

	err = client.CreateChannelByHTTP("test", true, nil, nil, 0, 1, 10000000)
	require.NoError(t, err)

	//issue to issuer itself
	coreTx := getAssetChannelTx(core.IssueContractAddress, issuer, "", uint64(10), issuerKey)
	_, err = client.AddTxByHTTP(coreTx)
	require.NoError(t, err)

	balance, err := client.GetAccountBalanceByHTTP(issuer)
	require.NoError(t, err)
	require.Equal(t, uint64(10), balance)

	//falseissuer issue fail
	coreTx = getAssetChannelTx(core.IssueContractAddress, falseIssuer, "", uint64(10), falseIssuerKey)
	_, err = client.AddTxByHTTP(coreTx)
	require.NoError(t, err)

	balance, err = client.GetAccountBalanceByHTTP(falseIssuer)
	require.NoError(t, err)
	require.Equal(t, uint64(0), balance)

	//test issue to channel
	coreTx = getAssetChannelTx(core.IssueContractAddress, common.ZeroAddress, "test", uint64(10), issuerKey)
	// question, what if test is not created?
	_, err = client.AddTxByHTTP(coreTx)
	require.NoError(t, err)
	balance, err = client.GetAccountBalanceByHTTP(common.AddressFromChannelID("test"))
	require.NoError(t, err)
	require.Equal(t, uint64(10), balance)

	//test transfer
	coreTx = getAssetChannelTx(core.TransferContractrAddress, receiver, "", uint64(5), issuerKey)
	_, err = client.AddTxByHTTP(coreTx)
	require.NoError(t, err)
	balance, err = client.GetAccountBalanceByHTTP(receiver)
	require.NoError(t, err)
	require.Equal(t, uint64(5), balance)

	//test transfer fail
	coreTx = getAssetChannelTx(core.TransferContractrAddress, receiver, "", uint64(5), falseIssuerKey)
	_, err = client.AddTxByHTTP(coreTx)
	require.NoError(t, err)
	balance, err = client.GetAccountBalanceByHTTP(receiver)
	require.NoError(t, err)
	require.Equal(t, uint64(5), balance)

	//4.test exchangeToken a.k.a transfer to channel in orderer execution
	coreTx = getAssetChannelTx(core.TokenExchangeAddress, common.ZeroAddress, "test", uint64(5), receiverKey)
	_, err = client.AddTxByHTTP(coreTx)
	require.NoError(t, err)

	balance, err = client.GetAccountBalanceByHTTP(common.AddressFromChannelID("test"))
	require.NoError(t, err)
	require.Equal(t, uint64(15), balance)

	token, err := client.GetTokenInfoByHTTP(receiver, []byte("test"))
	require.NoError(t, err)
	require.Equal(t, uint64(5), token)

	//test Block Price
	coreTx, err = core.NewTx("test", common.ZeroAddress, []byte("success"), 0, "", issuerKey)
	_, err = client.AddTxByHTTP(coreTx)
	require.NoError(t, err)

	//change BlockPrice of test channel's

	payload, err := json.Marshal(cc.Payload{
		ChannelID: "test",
		Profile: &cc.Profile{
			BlockPrice: 100,
		},
	})
	require.NoError(t, err)
	coreTx, err = core.NewTx(core.CONFIGCHANNELID, common.ZeroAddress, payload, 0, "", issuerKey)
	_, err = client.AddTxByHTTP(coreTx)
	require.NoError(t, err)

	//now add tx that cause due
	coreTx, err = core.NewTx("test", common.ZeroAddress, []byte("cause due but pass"), 0, "", issuerKey)
	_, err = client.AddTxByHTTP(coreTx)
	require.NoError(t, err)

	//this one should fail
	coreTx, err = core.NewTx("test", common.ZeroAddress, []byte("fail"), 0, "", issuerKey)
	_, err = client.AddTxByHTTP(coreTx)
	require.Error(t, err)

	//now issue money to channel account to wake it
	coreTx = getAssetChannelTx(core.IssueContractAddress, common.ZeroAddress, "test", uint64(1000000), issuerKey)
	_, err = client.AddTxByHTTP(coreTx)
	require.NoError(t, err)

	coreTx, err = core.NewTx("test", common.ZeroAddress, []byte("success again"), 0, "", issuerKey)
	_, err = client.AddTxByHTTP(coreTx)
	require.NoError(t, err)
}

func getAssetChannelTx(contract, addressInPayload common.Address, channelInPayload string, value uint64, privKey crypto.PrivateKey) *core.Tx {
	payload, _ := json.Marshal(asset.Payload{
		Address:   addressInPayload,
		ChannelID: channelInPayload,
	})
	coreTx, _ := core.NewTx(core.ASSETCHANNELID, contract, payload, value, "", privKey)
	return coreTx
}

func testAssetOld(t *testing.T, client *client.Client) {
	address, err := client.GetPrivKey().PubKey().Address()
	require.NoError(t, err)
	receiverPrivKey, err := crypto.NewPrivateKey([]byte("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"), crypto.KeyAlgoSecp256k1)
	require.NoError(t, err)
	receiverAddress, err := receiverPrivKey.PubKey().Address()
	require.NoError(t, err)
	require.NotEqual(t, receiverAddress, address)

	payload, err := json.Marshal(asset.Payload{
		//Action:    "person",
		ChannelID: "",
		Address:   receiverAddress,
	})
	tx, err := core.NewTx(core.ASSETCHANNELID, core.IssueContractAddress, payload, 10, "", client.GetPrivKey())
	require.NoError(t, err)
	status, err := client.AddTx(tx)
	require.NoError(t, err)
	require.Empty(t, status.Err)

	balance, err := client.GetAccountBalance(receiverAddress)
	require.NoError(t, err)
	require.Equal(t, balance, uint64(10))
	// then try to issue again, this should cause authentication error

	payload, err = json.Marshal(asset.Payload{
		//Action:    "person",
		ChannelID: "",
		Address:   address,
	})
	tx, err = core.NewTx(core.ASSETCHANNELID, core.IssueContractAddress, payload, 10, "", receiverPrivKey)
	require.NoError(t, err)
	status, err = client.AddTx(tx)
	require.NoError(t, err)
	require.NotEmpty(t, status.Err)

	payload, err = json.Marshal(asset.Payload{
		//Action:    "channel",
		ChannelID: "public",
	})
	tx, err = core.NewTx(core.ASSETCHANNELID, core.IssueContractAddress, payload, 10, "", client.GetPrivKey())
	require.NoError(t, err)
	status, err = client.AddTx(tx)
	require.NoError(t, err)
	require.Empty(t, status.Err)

	payload, err = json.Marshal(asset.Payload{
		//Action:    "channel",
		ChannelID: "public",
	})
	tx, err = core.NewTx(core.ASSETCHANNELID, core.TransferContractrAddress, payload, 10, "", receiverPrivKey)
	require.NoError(t, err)
	status, err = client.AddTx(tx)
	require.NoError(t, err)
	require.Empty(t, status.Err)

	balance, err = client.GetAccountBalance(address)
	require.NoError(t, err)
	require.Equal(t, balance, uint64(0))
	balance, err = client.GetAccountBalance(receiverAddress)
	require.NoError(t, err)
	require.Equal(t, balance, uint64(0))

}
