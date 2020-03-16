// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package eraft

import (
	"encoding/json"
	"madledger/core"
)

// Block is the implementation of raft Block
type Block struct {
	ChannelID string
	Num       uint64
	Txs       [][]byte
}

// GetNumber is the implementation of block
func (block *Block) GetNumber() uint64 {
	return block.Num
}

// GetTxs is the implementation of block
func (block *Block) GetTxs() []*core.Tx {
	var txs []*core.Tx
	for _, txBytes := range block.Txs {
		tx, err := core.BytesToTx(txBytes)
		if err == nil {
			txs = append(txs, tx)
		} else {
			log.Infof("get tx from consensus block failed because %v", err)
		}
	}
	return txs
}

// Bytes return bytes of block
func (block *Block) Bytes() []byte {
	bytes, _ := json.Marshal(block)
	return bytes
}

// UnmarshalBlock convert bytes to Block
func UnmarshalBlock(bytes []byte) *Block {
	var block Block
	json.Unmarshal(bytes, &block)
	return &block
}
