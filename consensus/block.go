package consensus

import "madledger/core"

// Block is the block interface
type Block interface {
	GetNumber() uint64
	GetTxs() []*core.Tx
}
