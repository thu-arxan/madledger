// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package common

import (
	"encoding/json"
	"errors"
	"madledger/common/crypto/hash"
	"madledger/common/math"
)

// Account is the Account of Madledger
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
	// TODO: Should we replace it with sm3?
	return hash.SHA256(a.Code)
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
