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
	"madledger/common"

	"github.com/thu-arxan/evm"
)

// Account is the implemantation of evm.Account, wraps common.Account.
type Account struct {
	account *common.Account
}

// NewAccount is the constructor of Account
// create a default account
func NewAccount(addr *Address) *Account {
	return &Account{
		account: common.NewAccount(common.BytesToAddress(addr.Bytes())),
	}
}

// NewAccountFromCommon ...
func NewAccountFromCommon(acc *common.Account) *Account {
	return &Account{
		account: acc,
	}
}

// CommonAccount ...
func (a *Account) CommonAccount() *common.Account {
	return a.account
}

// SetCode is the implementation of interface
func (a *Account) SetCode(code []byte) {
	a.account.SetCode(code)
}

// GetAddress is the implementation of interface
func (a *Account) GetAddress() evm.Address {
	return a.account.GetAddress()
}

// GetBalance is the implementation of interface
func (a *Account) GetBalance() uint64 {
	return a.account.GetBalance()
}

// GetCode is the implementation of interface
func (a *Account) GetCode() []byte {
	return a.account.GetCode()
}

// GetCodeHash return the hash of account code, please return [32]byte, and return [32]byte{0, ..., 0} if code is empty
func (a *Account) GetCodeHash() []byte {
	return a.account.GetCodeHash()
}

// AddBalance is the implementation of interface
// Note: In fact, we should avoid overflow
func (a *Account) AddBalance(balance uint64) error {
	return a.account.AddBalance(balance)
}

// SubBalance is the implementation of interface
func (a *Account) SubBalance(balance uint64) error {
	return a.account.SubBalance(balance)
}

// GetNonce is the implementation of interface
func (a *Account) GetNonce() uint64 {
	return a.account.GetNonce()
}

// SetNonce is the implementation of interface
func (a *Account) SetNonce(nonce uint64) {
	a.account.SetNonce(nonce)
}

// Suicide is the implementation of interface
func (a *Account) Suicide() {
	a.account.Suicide()
}

// HasSuicide is the implementation of interface
func (a *Account) HasSuicide() bool {
	return a.account.HasSuicide()
}

// Marshal marshal account into bytes
func (a *Account) Marshal() ([]byte, error) {
	return a.account.Bytes()
}
