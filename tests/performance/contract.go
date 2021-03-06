// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package performance

import (
	"madledger/common/abi"
	client "madledger/client/lib"
	"madledger/common"
	"madledger/core"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
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
	log             = logrus.WithFields(logrus.Fields{"app": "performance", "package": "tests"})
)

// CreateContract will create a channel by the client
func CreateContract(t *testing.T, channelID string, client *client.Client) {
	contractCodes, err := readCodes(BalanceBin)
	require.NoError(t, err)

	tx, err := core.NewTx(channelID, common.ZeroAddress, contractCodes, 0, "", client.GetPrivKey())
	require.NoError(t, err)

	status, err := client.AddTx(tx)
	require.NoError(t, err)
	require.Empty(t, status.Err)
	setContractAddress(channelID, common.HexToAddress(status.ContractAddress))
}

// CreateCallContractTx will create tx
func CreateCallContractTx(channelID string, client *client.Client, size int) []*core.Tx {
	var payload []byte
	payload, _ = abi.Pack(BalanceAbi, "get")
	var txs []*core.Tx
	for i := 0; i < size; i++ {
		tx, _ := core.NewTx(channelID, getContractAddress(channelID), payload, 0, "", client.GetPrivKey())
		txs = append(txs, tx)
	}

	return txs
}

// CallContract will call a contract of a channel
func CallContract(t *testing.T, channelID string, client *client.Client, times int) {
	var payload []byte
	payload, _ = abi.Pack(BalanceAbi, "get")
	var wg sync.WaitGroup
	for i := 0; i < times; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tx, _ := core.NewTx(channelID, getContractAddress(channelID), payload, 0, "", client.GetPrivKey())
			_, err := client.AddTx(tx)
			require.NoError(t, err)
		}()
	}
	wg.Wait()
}

// AddTxs add txs
func AddTxs(t *testing.T, client *client.Client, txs []*core.Tx) {
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
