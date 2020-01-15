package evm

import (
	"fmt"
	"madledger/common"
	"madledger/core"
	"madledger/executor/evm/madevm"
	"madledger/executor/evm/wildevm"
	"madledger/peer/db"
)

// Types defines which implemention to be used used
type Types uint32

const (
	// Wild marks using the original evm
	Wild = iota
	// Mad marks using the newest evm(in vendor)
	Mad
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
	Call(caller, callee common.Account, code, input []byte, value uint64) ([]byte, error)
	// Create create a contract.
	// MadEVM.Create(caller Address) ([]byte, Address, error)
	Create(caller common.Account, code, input []byte, value uint64) ([]byte, common.Address, error)
}

// Context defines the common functions for evm context
type Context interface {
}

// NewEVM is the constructor of evm, types marks which implemention to use
// func NewEVM(context Context, origin common.Address, db wildevm.StateDB, wb db.WriteBatch, types Types) EVM {
// 	switch types {
// 	case Wild:
// 		return wildevm.NewEVM(context, origin, db, wb)
// 	case Mad:
// 	default:
// 		panic(fmt.Errorf("NewEvm: invalid types %d", types))
// 	}
// 	return nil
// }

// NewWildEVM is the constructor of wild evm
func NewWildEVM(context Context, origin common.Address, db wildevm.StateDB, wb db.WriteBatch) EVM {
	ctx, ok := context.(*wildevm.Context)
	if !ok {
		panic("invalid context type, wildevm.Context expected")
	}
	return wildevm.NewEVM(ctx, origin, db, wb)
}

// NewContext is the constructor of context
func NewContext(block *core.Block, types Types) Context {
	switch types {
	case Wild:
		return wildevm.NewContext(block)
	case Mad:
		return madevm.NewMadContext(block)
	default:
		panic(fmt.Errorf("NewContext: invalid types %d", types))
	}
	// return nil
}
