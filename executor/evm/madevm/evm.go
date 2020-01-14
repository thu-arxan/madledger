package madevm

import (
	"evm"
	"madledger/core"
)

// MadContext is context for mad evm
type MadContext struct {
	block *core.Block
}

// NewMadContext ...
func NewMadContext(block *core.Block) *MadContext {
	return &MadContext{
		block: block,
	}
}

// MadEVM ...
type MadEVM struct {
	madEVM *evm.EVM
	ctx    *MadContext
}

// NewMadEVM ...
func NewMadEVM(ctx *MadContext) *MadEVM {
	return nil
}
