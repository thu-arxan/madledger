// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package evm

import (
	"github.com/thu-arxan/evm"
)

// Context caches data changed in block, passed to each evm for tx.
//
// Context syncs cached data into disk after all txs of block have been handled.
type Context interface {
	// BlockFinalize should be called after RunBlock.
	// In madevm, BlockFinalize will store logs for block generated during run txs of block into writebatch
	BlockFinalize() error
	// TxFinalize should be called after RunBlock.
	// In madevm, TxFinalize will remove info associated with suicide account
	// TxFinalize(db db.DB)

	// BlockContext returns evm ctx, fill in block info
	BlockContext() *evm.Context

	// NewBlockchain creates blockchain for evm.EVM
	NewBlockchain() evm.Blockchain
	// NewDatabase creates db for evm.EVM, caches data between txs in block
	NewDatabase() evm.DB
}
