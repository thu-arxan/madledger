// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package solo

import (
	"madledger/common/util"
	"madledger/consensus"
	"madledger/core"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	sc consensus.Consensus
)

func TestStart(t *testing.T) {
	var err error
	cfg := make(map[string]consensus.Config)
	cfg["test"] = consensus.DefaultConfig()
	sc, err = NewConsensus(cfg)
	require.NoError(t, err)
	sc.Start()
}

func TestAddTx(t *testing.T) {
	var txSize = 1024
	var txs []*core.Tx
	var success = make(map[string]int)
	var lock sync.Mutex

	for i := 0; i < txSize; i++ {
		tx := randomTx()
		success[tx.ID] = 0
		txs = append(txs, tx)
	}

	var wg sync.WaitGroup
	// each tx send 3 times
	for i := range txs {
		for m := 0; m < 3; m++ {
			wg.Add(1)
			tx := txs[i]
			go func() {
				defer wg.Done()
				if err := sc.AddTx(tx); err == nil {
					lock.Lock()
					success[tx.ID]++
					lock.Unlock()
				}
			}()
		}
	}
	wg.Wait()

	for i := range success {
		require.Equal(t, 1, success[i])
	}
}

func TestEnd(t *testing.T) {
	err := sc.Stop()
	require.NoError(t, err)
}

func randomTx() *core.Tx {
	return &core.Tx{
		ID: util.RandomString(32),
		Data: core.TxData{
			ChannelID: "test",
		},
	}
}
