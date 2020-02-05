package common

import (
	"encoding/json"
	"errors"
	"madledger/common/crypto/hash"
	"madledger/common/math"
)

// Account is the account of madledger
type Account interface {
	// These functions will be remained.
	GetAddress() Address
	GetBalance() uint64
	AddBalance(balance uint64) error
	SubBalance(balance uint64) error
	GetCode() []byte
	SetCode(code []byte)
	// GetNonce() uint64
	// SetNonce(nonce uint64)
	Bytes() ([]byte, error)
	// GetCodeHash return the hash of account code, please return [32]byte,
	// and return [32]byte{0, ..., 0} if code is empty
	GetCodeHash() []byte
	GetNonce() uint64
	SetNonce(nonce uint64)
	// Suicide will suicide an account
	Suicide()
	HasSuicide() bool
}

// DefaultAccount is the default implementation of Account
type DefaultAccount struct {
	Address     Address
	Balance     uint64
	Code        []byte
	Nonce       uint64
	SuicideMark bool
}

// NewDefaultAccount is the constructor of DefaultAccount
func NewDefaultAccount(addr Address) *DefaultAccount {
	return &DefaultAccount{
		Address: addr,
		Balance: 0,
		Code:    []byte{},
		Nonce:   0,
	}
}

// GetAddress is the implementation of Account
func (a *DefaultAccount) GetAddress() Address {
	return a.Address
}

// GetBalance is the implementation of Account
func (a *DefaultAccount) GetBalance() uint64 {
	return a.Balance
}

// AddBalance is the implementation of Account
func (a *DefaultAccount) AddBalance(balance uint64) error {
	if _, overflow := math.SafeAdd(a.Balance, balance); overflow {
		return errors.New("Overflow")
	}
	a.Balance += balance
	return nil
}

// SubBalance is the implementation of Account
func (a *DefaultAccount) SubBalance(balance uint64) error {
	if _, overflow := math.SafeSub(a.Balance, balance); !overflow {
		return errors.New("Overflow")
	}
	a.Balance -= balance
	return nil
}

// GetCode is the implementation of Account
func (a *DefaultAccount) GetCode() []byte {
	return a.Code
}

// SetCode is the implementation of Account
func (a *DefaultAccount) SetCode(code []byte) {
	a.Code = code
}

// Bytes is the implementation of Account
func (a *DefaultAccount) Bytes() ([]byte, error) {
	return json.Marshal(a)
}

// GetCodeHash return the hash of account code, please return [32]byte,
// and return [32]byte{0, ..., 0} if code is empty
func (a *DefaultAccount) GetCodeHash() []byte {
	bytes := make([]byte, 32)
	if len(a.Code) == 0 {
		return bytes
	}
	return hash.Hash(a.Code)
}

// GetNonce ...
func (a *DefaultAccount) GetNonce() uint64 {
	return a.Nonce
}

// SetNonce ...
func (a *DefaultAccount) SetNonce(nonce uint64) {
	a.Nonce = nonce
}

// Suicide will suicide an account
func (a *DefaultAccount) Suicide() {
	a.SuicideMark = true
}

// HasSuicide returns if account has suicided
func (a *DefaultAccount) HasSuicide() bool {
	return a.SuicideMark
}
