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
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"madledger/common"
	"madledger/core"
	"madledger/peer/db"

	"github.com/thu-arxan/evm"
	"github.com/thu-arxan/evm/util"

	"github.com/syndtr/goleveldb/leveldb"
)

// DefaultContext is the default implementation for Context
type DefaultContext struct {
	queryEngine db.DB         // queryEngine used to query data from db
	wb          db.WriteBatch // wb used to put data into db after finishing block

	channelID string

	// evm.Context
	block  *core.Block
	evmCtx *evm.Context
	logs   []*evm.Log

	accounts map[string]*accountInfo
}

type storageData struct {
	value   []byte
	updated bool
}

type accountInfo struct {
	account *Account
	storage map[string]*storageData // key => value
	updated bool
}

// NewContext is the constructor of context
func NewContext(block *core.Block, engine db.DB, wb db.WriteBatch) Context {
	return &DefaultContext{
		queryEngine: engine,
		wb:          wb,
		block:       block,
		channelID:   block.Header.ChannelID,
		evmCtx: &evm.Context{
			BlockHeight: block.GetNumber(),
			BlockTime:   block.Header.Time,
			Difficulty:  0,
			GasLimit:    10000000000,
			GasPrice:    0,
			CoinBase:    nil,
		},
		accounts: make(map[string]*accountInfo),
	}
}

// BlockFinalize should be called after RunBlock.
// In madevm, BlockFinalize will store logs for block generated during run txs of block into writebatch
func (ctx *DefaultContext) BlockFinalize() error {
	if len(ctx.logs) != 0 {
		logs, err := json.Marshal(ctx.logs)
		if err != nil {
			return err
		}
		ctx.wb.Put([]byte(fmt.Sprintf("block_log_%s_%d", ctx.channelID, ctx.block.GetNumber())), logs)
	}

	for addr, acc := range ctx.accounts {
		// sync account
		if acc.updated {
			if err := ctx.wb.SetAccount(acc.account.CommonAccount()); err != nil {
				return err
			}
		}
		commonAddr := bytesToCommonAddress([]byte(addr))
		if acc.account.HasSuicide() {
			// remove suicide account's storage
			ctx.wb.RemoveAccountStorage(commonAddr)
			continue
		}
		// sync storage
		for key, v := range acc.storage {
			if v.updated {
				if err := ctx.wb.SetStorage(commonAddr, bytesToCommomWord256([]byte(key)), bytesToCommomWord256(v.value)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// BlockContext returns evm ctx, fill in block info
func (ctx *DefaultContext) BlockContext() *evm.Context {
	return ctx.evmCtx
}

// NewBlockchain creates blockchain for evm.EVM
func (ctx *DefaultContext) NewBlockchain() evm.Blockchain {
	return NewBlockchain(ctx.queryEngine, ctx.channelID)
}

// NewDatabase creates db for evm.EVM, caches data between txs in block
func (ctx *DefaultContext) NewDatabase() evm.DB {
	return NewCache(ctx)
}

// for evm.DB, just query and cache

func bytesToCommonAddress(addr []byte) common.Address {
	return common.BytesToAddress(addr)
}

func bytesToCommomWord256(data []byte) common.Word256 {
	word256, err := common.BytesToWord256(data)
	if err != nil {
		log.Errorf("Fatal error! Failed to convert bytes(%d) to word256, %s, err: %v", len(data), hex.EncodeToString(data), err)
	}
	return word256
}

// exist returns if address exist
func (ctx *DefaultContext) exist(address []byte) bool {
	acc := ctx.accounts[string(address)]
	if acc != nil {
		return true
	}
	return ctx.queryEngine.AccountExist(bytesToCommonAddress(address))
}

func (ctx *DefaultContext) getOrSetAccountInfo(addr []byte) *accountInfo {
	accInfo := ctx.accounts[string(addr)]
	if accInfo != nil {
		return accInfo
	}
	ctx.getAccount(addr)
	return ctx.accounts[string(addr)]
}

func (ctx *DefaultContext) getAccount(address []byte) evm.Account {
	addrStr := string(address)

	if acc := ctx.accounts[addrStr]; acc != nil {
		return acc.account
	}

	// query from db
	account, err := ctx.queryEngine.GetAccount(bytesToCommonAddress(address))
	if err != nil {
		// err returns, default one
		log.Errorf("Fatal! failed to query account for %s, err: %v", string(address), err)
		defaultAcc := NewAccount(BytesToAddress(address))
		ctx.accounts[addrStr] = &accountInfo{
			account: defaultAcc,
			updated: true,
			storage: make(map[string]*storageData),
		}
	}

	ctx.accounts[addrStr] = &accountInfo{
		account: NewAccountFromCommon(account),
		storage: make(map[string]*storageData),
	}

	return ctx.accounts[addrStr].account
}

func (ctx *DefaultContext) getStorage(addr, key []byte) []byte {
	accInfo := ctx.getOrSetAccountInfo(addr)
	if util.Contain(accInfo.storage, string(key)) {
		return accInfo.storage[string(key)].value
	}

	// query from db
	value, err := ctx.queryEngine.GetStorage(bytesToCommonAddress(addr), bytesToCommomWord256(key))
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("Fatal error! Failed to query value to %s for addr(%s), err: %v", hex.EncodeToString(key), hex.EncodeToString(addr), err)
	}

	accInfo.storage[string(key)] = &storageData{
		value: value.Bytes(),
	}
	return value.Bytes()
}

// for evm.WriteBatch, stored into cache, sync when BlockFinalize
func (ctx *DefaultContext) setStorage(addr, key, value []byte) {
	// todo: removed account?
	accInfo := ctx.getOrSetAccountInfo(addr)
	if accInfo.account.HasSuicide() {
		log.Errorf("Fatal error, set storage on a suicide account(%s), key: %s, value: %s", string(addr), string(key), string(value))
	}
	accInfo.storage[string(key)] = &storageData{
		value:   value,
		updated: true,
	}
}

func (ctx *DefaultContext) updateAccount(account evm.Account) error {
	addr := account.GetAddress().Bytes()
	accInfo := ctx.getOrSetAccountInfo(addr)

	if accInfo.account.HasSuicide() {
		return fmt.Errorf("Fatal error, UpdateAccount on a suicide account: %s", string(addr))
	}

	acc, ok := account.(*Account)
	if !ok {
		return errors.New("invalid account type, executor/evm.Account expected")
	}

	accInfo.account = acc
	accInfo.updated = true
	return nil
}

func (ctx *DefaultContext) addLog(log *evm.Log) {
	ctx.logs = append(ctx.logs, log)
}
