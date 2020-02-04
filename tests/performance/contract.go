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
	ContractAddress = make(map[string]common.Address, 0)
	mapLock         sync.Mutex
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
	setContractAddress(channelID, common.HexToAddress(status.ContractAddress))
}

// CreateCallContractTx will create tx
func CreateCallContractTx(channelID string, client *client.Client, size int) []*types.Tx {
	var payload []byte
	payload, _ = abi.GetPayloadBytes(BalanceAbi, "get", nil)
	var txs []*types.Tx
	for i := 0; i < size; i++ {
		tx, _ := types.NewTx(channelID, getContractAddress(channelID), payload, client.GetPrivKey())
		txs = append(txs, tx)
	}

	return txs
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
			tx, _ := types.NewTx(channelID, getContractAddress(channelID), payload, client.GetPrivKey())
			_, err := client.AddTx(tx)
			require.NoError(t, err)
		}()
	}
	wg.Wait()
}

// AddTxs add txs
func AddTxs(t *testing.T, client *client.Client, txs []*types.Tx) {
	var wg sync.WaitGroup
	for i := 0; i < len(txs); i++ {
		wg.Add(1)
		tx := txs[i]
		go func() {
			defer wg.Done()
			_, err := client.AddTx(tx)
			require.NoError(t, err)
		}()
	}
	wg.Wait()
}

func setContractAddress(channelID string, address common.Address) {
	mapLock.Lock()
	defer mapLock.Unlock()

	ContractAddress[channelID] = address
}

func getContractAddress(channelID string) common.Address {
	mapLock.Lock()
	defer mapLock.Unlock()

	return ContractAddress[channelID]
}
