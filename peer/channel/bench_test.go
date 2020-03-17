// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package channel

import (
	"encoding/hex"
	"fmt"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/core"
	"madledger/peer/db"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	benckmark          = true
	benchbenchBlockNum = 100
	benchBlockNum      = 100
	payload, _         = hex.DecodeString("6080604052600a600060005090905534801561001b5760006000fd5b50610021565b61026d806100306000396000f3fe60806040523480156100115760006000fd5b506004361061005c5760003560e01c80631003e2d21461006257806327ee58a6146100a5578063370158ea146100e857806360fe47b1146101395780636d4ce63c146101805761005c565b60006000fd5b61008f600480360360208110156100795760006000fd5b810190808035906020019092919050505061019e565b6040518082815260200191505060405180910390f35b6100d2600480360360208110156100bc5760006000fd5b81019080803590602001909291905050506101c6565b6040518082815260200191505060405180910390f35b6100f06101ee565b604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b610166600480360360208110156101505760006000fd5b8101908080359060200190929190505050610209565b604051808215151515815260200191505060405180910390f35b610188610225565b6040518082815260200191505060405180910390f35b6000816000600082828250540192505081909090555060006000505490506101c1565b919050565b6000816000600082828250540392505081909090555060006000505490506101e9565b919050565b600060003360006000505481915091509150610205565b9091565b600081600060005081909090555060019050610220565b919050565b60006000600050549050610234565b9056fea26469706673582212206d1f7e72f2d26816fe48ff60de6fa42d7b6fb40fc1603494b8c63221cd3c2c2364736f6c63430006000033")
)

func TestRunBlock(t *testing.T) {
	if !benckmark {
		return
	}
	os.RemoveAll(".benchmark")
	var blocks = make([]*core.Block, benchBlockNum)
	for i := 0; i < benchBlockNum; i++ {
		var txs = make([]*core.Tx, benchBlockNum)
		for j := range txs {
			privKey, _ := crypto.GeneratePrivateKey()
			tx, err := core.NewTx("benchmark", common.ZeroAddress, payload, 0, "", privKey)
			require.NoError(t, err)
			txs[j] = tx
		}
		blocks[i] = core.NewBlock("benchmark", uint64(i), nil, txs)
	}
	db, err := db.NewLevelDB(".benchmark")
	require.NoError(t, err)
	manager, err := NewManager("benchmark", ".benchmark", nil, db, nil, nil)
	require.NoError(t, err)
	var begin = time.Now()
	for i := range blocks {
		_, err := manager.RunBlock(blocks[i])
		require.NoError(t, err)
	}
	duration := time.Since(begin)
	fmt.Printf("run %d tx cost %v\n", benchBlockNum*benchBlockNum, duration)
	os.RemoveAll(".benchmark")
}
