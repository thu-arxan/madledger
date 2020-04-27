// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package global

import (
	"encoding/json"
	"madledger/core"
)

// CreateGenesisBlock return the genesis block.
func CreateGenesisBlock(payloads []*core.GlobalTxPayload) (*core.Block, error) {
	var txs []*core.Tx
	for _, payload := range payloads {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		// all zero
		tx := core.NewTxWithoutSig(core.GLOBALCHANNELID, payloadBytes, 0)
		txs = append(txs, tx)
	}

	return core.NewBlock(core.GLOBALCHANNELID, 0, core.GenesisBlockPrevHash, txs), nil
}
