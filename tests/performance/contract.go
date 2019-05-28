package performance

import (
	client "madledger/client/lib"
	"madledger/common"
	"madledger/common/abi"
	"madledger/core/types"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	// BalanceBin is the path of balance contract bin
	BalanceBin = "../balance/Balance.bin"
	// BalanceAbi is the path of balance contract biabin
	BalanceAbi = "../balance/Balance.abi"
	// ContractAddress is the address of contract
	ContractAddress common.Address
)

// CreateContract will create a channel by the client
func CreateContract(t *testing.T, channelID string, client *client.Client) {
	contractCodes, err := readCodes(BalanceBin)
	require.NoError(t, err)

	tx, err := types.NewTx(channelID, common.ZeroAddress, contractCodes, client.GetPrivKey())
	require.NoError(t, err)

	status, err := client.AddTx(tx)
	require.NoError(t, err)
	require.Empty(t, status.Err)
	ContractAddress = common.HexToAddress(status.ContractAddress)
}

// CallContract will call a contract of a channel
func CallContract(t *testing.T, channelID string, client *client.Client, times int) {
	var payload []byte
	payload, _ = abi.GetPayloadBytes(BalanceAbi, "get", nil)
	var wg sync.WaitGroup
	for i := 0; i < times; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tx, _ := types.NewTx(channelID, ContractAddress, payload, client.GetPrivKey())
			_, err := client.AddTx(tx)
			require.NoError(t, err)
		}()
	}
	wg.Wait()
}
