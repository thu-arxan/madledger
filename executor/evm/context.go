package evm

import (
	"madledger/common"
	"madledger/core/types"
)

// Context provide a context to run a contract on the evm
type Context struct {
	// Number is the number of the block
	Number uint64
	// BlockHash is the hash of the block
	BlockHash common.Hash
	// BlockTime is the time of the block
	BlockTime int64
	// GasLimit limit the use of gas, now is useless
	GasLimit uint64
	// CoinBase, set it to zero
	CoinBase common.Word256
	// diffculty is zero
	Diffculty uint64
}

// NewContext is the constructor of Context
func NewContext(block *types.Block) *Context {
	return &Context{
		Number:    block.Header.Number,
		BlockHash: block.Hash(),
		BlockTime: block.Header.Time,
		GasLimit:  0,
		CoinBase:  common.ZeroWord256,
		Diffculty: 0,
	}
}
