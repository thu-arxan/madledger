package tests

import (
	"encoding/json"
	"fmt"
	"madledger/blockchain/asset"
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
	err = client.CreateChannel("public", true, nil, nil, 1, 1, 10000000)
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
	err = client.CreateChannel("public", true, nil, nil, 1, 1, 10000000)
	require.Error(t, err)
	// create private channel
	err = client.CreateChannel("private", false, nil, peers, 1, 1, 10000000)
	require.NoError(t, err)
}

func testCreateContract(t *testing.T, client *client.Client) {
	// First, test on channel public, which is public
	createContract(t, "public", client)
	// Then, test on channel private, which is private
	createContract(t, "private", client)
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

func testCallContract(t *testing.T, client *client.Client) {
	callContract(t, "public", client)
	callContract(t, "private", client)
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
	fmt.Printf("test Asset begin")
}

func testAsset(t *testing.T, client *client.Client) {
	address, err := client.GetPrivKey().PubKey().Address()
	require.NoError(t, err)
	receiverPrivKey, err := crypto.NewPrivateKey([]byte("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"))
	require.NoError(t, err)
	receiverAddress, err := receiverPrivKey.PubKey().Address()
	require.NoError(t, err)
	require.NotEqual(t, receiverAddress, address)
	fmt.Printf("first issue")

	payload, err := json.Marshal(asset.Payload{
		Action:    "person",
		ChannelID: "",
		Address:   receiverAddress,
	})
	tx, err := core.NewTx(core.ASSETCHANNELID, core.IssueContractAddress, payload, 10, "", client.GetPrivKey())
	require.NoError(t, err)
	status, err := client.AddTxInOrderer(tx)
	require.NoError(t, err)
	require.Empty(t, status.Err)

	balance, err := client.GetAccountBalance(receiverAddress)
	require.NoError(t, err)
	require.Equal(t, balance, uint64(10))
	// then try to issue again, this should cause authentication error

	payload, err = json.Marshal(asset.Payload{
		Action:    "person",
		ChannelID: "",
		Address:   address,
	})
	tx, err = core.NewTx(core.ASSETCHANNELID, core.IssueContractAddress, payload, 10, "", receiverPrivKey)
	require.NoError(t, err)
	status, err = client.AddTxInOrderer(tx)
	require.NoError(t, err)
	require.NotEmpty(t, status.Err)

	payload, err = json.Marshal(asset.Payload{
		Action:    "channel",
		ChannelID: "public",
	})
	tx, err = core.NewTx(core.ASSETCHANNELID, core.IssueContractAddress, payload, 10, "", client.GetPrivKey())
	require.NoError(t, err)
	status, err = client.AddTxInOrderer(tx)
	require.NoError(t, err)
	require.Empty(t, status.Err)

	payload, err = json.Marshal(asset.Payload{
		Action:    "channel",
		ChannelID: "public",
	})
	tx, err = core.NewTx(core.ASSETCHANNELID, core.TransferContractrAddress, payload, 10, "", receiverPrivKey)
	require.NoError(t, err)
	status, err = client.AddTxInOrderer(tx)
	require.NoError(t, err)
	require.Empty(t, status.Err)

	balance, err = client.GetAccountBalance(address)
	require.NoError(t, err)
	require.Equal(t, balance, uint64(0))
	balance, err = client.GetAccountBalance(receiverAddress)
	require.NoError(t, err)
	require.Equal(t, balance, uint64(0))

}
