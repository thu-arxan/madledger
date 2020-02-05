package evm

import (
	"evm"
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
