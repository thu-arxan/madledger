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

import "madledger/core"

// Block is the implementaion of solo Block
type Block struct {
	channelID string
	num       uint64
	txs       []*core.Tx
}

// GetNumber is the implementation of block
func (block *Block) GetNumber() uint64 {
	return block.num
}

// GetTxs is the implementation of block
func (block *Block) GetTxs() []*core.Tx {
	return block.txs
}
