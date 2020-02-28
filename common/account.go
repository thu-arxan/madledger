package common

import (
	"encoding/json"
	"errors"
	"madledger/common/crypto/hash"
	"madledger/common/math"
)

// Account is the account of madledger
// type Account interface {
// 	// These functions will be remained.
// 	GetAddress() Address
// 	GetBalance() uint64
// 	AddBalance(balance uint64) error
// 	SubBalance(balance uint64) error
// 	GetCode() []byte
// 	SetCode(code []byte)
// 	// GetNonce() uint64
// 	// SetNonce(nonce uint64)
// 	Bytes() ([]byte, error)
// 	// GetCodeHash return the hash of account code, please return [32]byte,
// 	// and return [32]byte{0, ..., 0} if code is empty
// 	GetCodeHash() []byte
// 	GetNonce() uint64
// 	SetNonce(nonce uint64)
// 	// Suicide will suicide an account
// 	Suicide()
// 	HasSuicide() bool
// }

// Account is the default implementation of Account
type Account struct {
	Address     Address
	Balance     uint64
	Code        []byte
	Nonce       uint64
	SuicideMark bool
}

// NewAccount is the constructor of Account
func NewAccount(addr Address) *Account {
	return &Account{
		Address: addr,
		Balance: 0,
		Code:    []byte{},
		Nonce:   0,
	}
}

// GetAddress is the implementation of Account
func (a *Account) GetAddress() Address {
	return a.Address
}

// GetBalance is the implementation of Account
func (a *Account) GetBalance() uint64 {
	return a.Balance
}

// AddBalance is the implementation of Account
func (a *Account) AddBalance(balance uint64) error {
	if _, overflow := math.SafeAdd(a.Balance, balance); overflow {
		return errors.New("Overflow")
	}
	a.Balance += balance
	return nil
}

// SubBalance is the implementation of Account
func (a *Account) SubBalance(balance uint64) error {
	if _, overflow := math.SafeSub(a.Balance, balance); overflow {
		return errors.New("Overflow")
	}
	a.Balance -= balance
	return nil
}

// GetCode is the implementation of Account
func (a *Account) GetCode() []byte {
	return a.Code
}

// SetCode is the implementation of Account
func (a *Account) SetCode(code []byte) {
	a.Code = code
}

// Bytes is the implementation of Account
func (a *Account) Bytes() ([]byte, error) {
	return json.Marshal(a)
}

// GetCodeHash return the hash of account code, please return [32]byte,
// and return [32]byte{0, ..., 0} if code is empty
func (a *Account) GetCodeHash() []byte {
	bytes := make([]byte, 32)
	if len(a.Code) == 0 {
		return bytes
	}
	return hash.Hash(a.Code)
}

// GetNonce ...
func (a *Account) GetNonce() uint64 {
	return a.Nonce
}

// SetNonce ...
func (a *Account) SetNonce(nonce uint64) {
	a.Nonce = nonce
}

// Suicide will suicide an account
func (a *Account) Suicide() {
	a.SuicideMark = true
}

// HasSuicide returns if account has suicided
func (a *Account) HasSuicide() bool {
	return a.SuicideMark
}
