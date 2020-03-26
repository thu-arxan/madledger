// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package channel

import (
	"errors"
	"madledger/blockchain"
	"madledger/common"
	"madledger/core"

	"madledger/executor/evm"
	"madledger/peer/db"
	"madledger/peer/orderer"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "peer", "package": "channel"})
)

// Manager is the manager of channel
type Manager struct {
	signalCh chan bool
	stopCh   chan bool

	identity *core.Member
	// id is the id of channel
	id string
	// db is the database
	db db.DB
	// chain manager
	cm          *blockchain.Manager
	clients     []*orderer.Client
	coordinator *Coordinator
}

// NewManager is the constructor of Manager
func NewManager(id, dir string, identity *core.Member, db db.DB, clients []*orderer.Client, coordinator *Coordinator) (*Manager, error) {
	cm, err := blockchain.NewManager(id, dir)
	if err != nil {
		return nil, err
	}
	return &Manager{
		signalCh:    make(chan bool, 1),
		stopCh:      make(chan bool, 1),
		identity:    identity,
		id:          id,
		db:          db,
		cm:          cm,
		clients:     clients,
		coordinator: coordinator,
	}, nil
}

// Start start the manager.
func (m *Manager) Start() {
	log.Infof("channel %s is starting...", m.id)
	for {
		block, err := m.fetchBlock()
		if err == nil {
			// fmt.Println("Succeed to fetch block", m.id, ":", block.Header.Number)
			m.AddBlock(block)
		} else if err.Error() == "Stop" {
			m.stopCh <- true
			return
		} else {
			log.Warnf("failed to fetch block: %d, err: %v", m.cm.GetExpect(), err)
		}
	}
}

// Stop will stop the manager
// TODO: find a good way to stop
func (m *Manager) Stop() {
	m.signalCh <- true
	<-m.stopCh
}

// AddBlock add a block
func (m *Manager) AddBlock(block *core.Block) error {
	// add into the blockchain
	err := m.cm.AddBlock(block)
	if err != nil {
		return err
	}
	switch block.Header.ChannelID {
	case core.GLOBALCHANNELID:
		return m.AddGlobalBlock(block)
	case core.CONFIGCHANNELID:
		if !isGenesisBlock(block) && !m.coordinator.CanRun(block.Header.ChannelID, block.Header.Number) {
			m.coordinator.Watch(block.Header.ChannelID, block.Header.Number)
		}
		return m.AddConfigBlock(block)
	case core.ASSETCHANNELID:
		if !isGenesisBlock(block) && !m.coordinator.CanRun(block.Header.ChannelID, block.Header.Number) {
			m.coordinator.Watch(block.Header.ChannelID, block.Header.Number)
		}
		return m.AddAssetBlock(block)
	default:
		if !m.coordinator.CanRun(block.Header.ChannelID, block.Header.Number) {
			m.coordinator.Watch(block.Header.ChannelID, block.Header.Number)
		}
		log.Infof("Run block %s: %d", m.id, block.Header.Number)
		wb, err := m.RunBlock(block)
		if err != nil {
			return err
		}

		wb.PutBlock(block)
		return wb.Sync()
	}
	return nil
}

func isGenesisBlock(block *core.Block) bool {
	return block.GetNumber() == 0
}

// RunBlock will carry out all txs in the block.
// It will return after the block is runned.
// In the future, this will contains chains which rely on something or nothing
func (m *Manager) RunBlock(block *core.Block) (db.WriteBatch, error) {
	cache := NewCache(m.db)
	context := evm.NewContext(block, cache.db, cache.wb)
	defer context.BlockFinalize()
	// first parallel get sender to speed up
	threadSize := runtime.NumCPU()
	if threadSize < 2 {
		threadSize = 2
	}
	var ch = make(chan bool, threadSize)
	var wg sync.WaitGroup
	for i := range block.Transactions {
		wg.Add(1)
		tx := block.Transactions[i]
		ch <- true
		go func() {
			defer func() {
				<-ch
				wg.Done()
			}()
			tx.GetSender()
		}()
	}
	wg.Wait()

	profile, err := m.db.GetChannelProfile(m.id)
	if err != nil {
		return nil, err
	}

	for i, tx := range block.Transactions {
		senderAddress, err := tx.GetSender()
		status := &db.TxStatus{
			Err:         "",
			BlockNumber: block.Header.Number,
			BlockIndex:  i,
			Output:      nil,
		}
		if err != nil {
			status.Err = err.Error()
			cache.SetTxStatus(tx, status)
			continue
		}
		receiverAddress := tx.GetReceiver()

		sender, err := m.db.GetAccount(senderAddress)
		if err != nil {
			status.Err = err.Error()
			cache.SetTxStatus(tx, status)
			continue
		}

		if receiverAddress.String() == core.CfgTendermintAddress.String() {
			cache.SetTxStatus(tx, status)
			continue
		}

		if receiverAddress.String() == core.CfgRaftAddress.String() {
			cache.SetTxStatus(tx, status)
			continue
		}

		// 用户的参数：tx.Data.Gas (user gas limit)
		// 通道的参数：maxGas (channel gas limit), gasPrice
		// gas limit = min (user, channel)
		// 获取sender的token，如果比gas limit * gas price 小，那么不能执行，直接下一个tx
		// 记录进入evm前的gas limit
		// 用出来之后用前减后可得到具体消耗了多少gas
		// 然后将token -= gas * gas price，存到cache中

		gasLimit := profile.MaxGas
		if gasLimit > tx.Data.Gas {
			gasLimit = tx.Data.Gas
		}
		tokenLeft, err := cache.GetToken(m.id, senderAddress)
		if err != nil {
			status.Err = err.Error()
			cache.SetTxStatus(tx, status)
			continue
		}
		if tokenLeft < gasLimit*profile.GasPrice {
			status.Err = "Not enough token"
			cache.SetTxStatus(tx, status)
			continue
		}

		evm := evm.NewEVM(context, senderAddress, tx.Data.Payload, tx.Data.Value, gasLimit, cache.db, cache.wb)

		if receiverAddress.String() != common.ZeroAddress.String() {
			// if the length of payload is not zero, this is a contract call
			if len(tx.Data.Payload) != 0 && !m.db.AccountExist(receiverAddress) {
				status.Err = "Invalid Address"
				cache.SetTxStatus(tx, status)
				continue
			}

			receiver, err := m.db.GetAccount(receiverAddress)
			if err != nil {
				status.Err = err.Error()
				cache.SetTxStatus(tx, status)
				continue
			}
			output, err := evm.Call(sender, receiver, receiver.GetCode())
			status.Output = output
			if err != nil {
				status.Err = err.Error()
			}
			cache.SetTxStatus(tx, status)
		} else {
			output, addr, err := evm.Create(sender)
			status.Output = output
			status.ContractAddress = addr.String()
			if err != nil {
				status.Err = err.Error()
			}
			cache.SetTxStatus(tx, status)
		}
		gasUsed := gasLimit - *context.BlockContext().Gas
		tokenLeft -= gasUsed * profile.GasPrice
		cache.SetToken(m.id, senderAddress, tokenLeft)
	}
	return cache.wb, nil
}

// todo: here we should support evil orderer
func (m *Manager) fetchBlock() (*core.Block, error) {
	var lock sync.RWMutex
	var errs = make(chan error, len(m.clients))
	var blocks = make(chan *core.Block, len(m.clients))
	closed := false
	id := m.id
	expect := m.cm.GetExpect()
	for i := range m.clients {
		go func(i int) {
			block, err := m.clients[i].FetchBlock(id, expect, true)
			lock.RLock()
			defer lock.RUnlock()
			if closed {
				return
			}
			if err != nil {
				errs <- err
			} else {
				blocks <- block
			}
		}(i)
	}

	fails := 0

	for {
		defer func() {
			lock.Lock()
			if !closed {
				close(errs)
				close(blocks)
				closed = true
			}
			lock.Unlock()
		}()
		// log.Infof("get %s %d, closed: %t, fail: %d, %d", id, except, closed, fails, len(m.clients))
		select {
		case block := <-blocks:
			return block, nil
		case err := <-errs:
			fails++
			if fails == len(m.clients) {
				return nil, err
			}
		case <-m.signalCh:
			return nil, errors.New("Stop")
		}
	}
}
