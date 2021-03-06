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

import "github.com/thu-arxan/evm"

// Cache is the tx cache
type Cache struct {
	ctx *DefaultContext
}

// NewCache ...
func NewCache(ctx *DefaultContext) evm.DB {
	return &Cache{
		ctx: ctx,
	}
}

// Exist return if the account exist
// Note: if account is suicided, return true
func (c *Cache) Exist(address evm.Address) bool {
	return c.ctx.exist(address.Bytes())
}

// GetAccount return a default account if unexist
func (c *Cache) GetAccount(address evm.Address) evm.Account {
	return c.ctx.getAccount(address.Bytes())
}

// GetStorage get stored value associated with addr+key
func (c *Cache) GetStorage(address evm.Address, key []byte) (value []byte) {
	return c.ctx.getStorage(address.Bytes(), key)
}

// NewWriteBatch ...
func (c *Cache) NewWriteBatch() evm.WriteBatch {
	return c
}

// functions for writebatch

// SetStorage ...
func (c *Cache) SetStorage(address evm.Address, key []byte, value []byte) {
	c.ctx.setStorage(address.Bytes(), key, value)
	return
}

// UpdateAccount ...
// Note: db should delete all storages if an account suicide
func (c *Cache) UpdateAccount(account evm.Account) error {
	return c.ctx.updateAccount(account)
}

// AddLog ...
func (c *Cache) AddLog(log *evm.Log) {
	c.ctx.addLog(log)
	return
}
