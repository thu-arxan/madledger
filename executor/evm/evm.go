package evm

import (
	"evm"
	"madledger/common"
	"madledger/peer/db"
)

// EVM defines the common functions for evm
type EVM interface {
	// Call run code on a evm.
	//
	// MadEVM.Call(caller, callee Address, code []byte) ([]byte, error)
	/*
		// Address describe what functions that an Address implementation should provide
		type Address interface {
			// It would be better if length = 32
			// 1. Add zero in left if length < 32
			// 2. Remove left byte if length > 32(however, this may be harm)
			Bytes() []byte
		}
	*/
	Call(caller, callee common.Account, code []byte) ([]byte, error)
	// Create create a contract.
	// MadEVM.Create(caller Address) ([]byte, Address, error)
	Create(caller common.Account) ([]byte, common.Address, error)
}

// DefaultEVM ...
type DefaultEVM struct {
	runner *evm.EVM
	ctx    Context
}

// NewEVM ...
func NewEVM(ctx Context, caller common.Address, payload []byte, value uint64, gas uint64, engine db.DB, wb db.WriteBatch) EVM {
	// todo
	// bc := NewBlockchain(engine)
	// database := NewMemory(bc.NewAccount, engine, ctx)

	// ctx.setDB(engine)

	evmCtx := ctx.BlockContext()
	evmCtx.Input = payload
	evmCtx.Value = value
	evmCtx.Gas = &gas

	return &DefaultEVM{
		ctx:    ctx,
		runner: evm.New(ctx.NewBlockchain(), ctx.NewDatabase(), evmCtx),
	}
}

// Call ...
func (evm *DefaultEVM) Call(caller, callee common.Account, code []byte) ([]byte, error) {
	return evm.runner.Call(caller.GetAddress(), callee.GetAddress(), code)
}

// Create ...
func (evm *DefaultEVM) Create(caller common.Account) ([]byte, common.Address, error) {
	v, addr, err := evm.runner.Create(caller.GetAddress())
	if addr == nil {
		return v, common.ZeroAddress, err
	}
	return v, common.BytesToAddress(addr.Bytes()), err
}
