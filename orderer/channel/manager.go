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
	"fmt"
	"madledger/blockchain"
	"madledger/common"
	"madledger/common/event"
	"madledger/common/util"
	"madledger/consensus"
	"madledger/consensus/raft"
	"madledger/core"
	"madledger/orderer/db"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "orderer", "package": "channel"})
)

// Manager is the manager of channel
type Manager struct {
	// ID is the id of channel
	ID string
	// db is the database
	db db.DB
	// chain manager
	cm *blockchain.Manager
	// consensus block chan
	cbc         chan consensus.Block
	init        bool
	stop        chan bool
	hub         *event.Hub
	coordinator *Coordinator

	lock                sync.RWMutex
	insufficientBalance bool
}

// NewManager is the constructor of Manager
func NewManager(id string, coordinator *Coordinator) (*Manager, error) {
	cm, err := blockchain.NewManager(id, fmt.Sprintf("%s/%s", coordinator.chainCfg.Path, id))
	if err != nil {
		return nil, err
	}
	return &Manager{
		ID:                  id,
		db:                  coordinator.db,
		cm:                  cm,
		cbc:                 make(chan consensus.Block, 1024),
		init:                false,
		stop:                make(chan bool),
		hub:                 event.NewHub(),
		coordinator:         coordinator,
		insufficientBalance: false,
	}, nil
}

// Start starts the channel
func (manager *Manager) Start() {
	log.Infof("Channel %s is starting", manager.ID)
	manager.init = true
	go manager.syncBlock()
	for {
		select {
		case cb := <-manager.cbc:
			// log.Infof("Receive block %s:%d from consensus", manager.ID, cb.GetNumber())
			// todo: if a tx is duplicated and it was added into consensus block succeed, then it may never receive response
			txs, _ := manager.getTxsFromConsensusBlock(cb)
			if len(txs) != 0 {
				prevBlock := manager.cm.GetPrevBlock()
				var block *core.Block
				if prevBlock == nil {
					block = core.NewBlock(manager.ID, 0, core.GenesisBlockPrevHash, txs)
					log.Debugf("Channel %s create new block %d, hash is %s", manager.ID, 0, util.Hex(block.Hash().Bytes()))
				} else {
					block = core.NewBlock(manager.ID, prevBlock.Header.Number+1, prevBlock.Hash().Bytes(), txs)
					log.Debugf("Channel %s create new block %d, hash is %s", manager.ID, prevBlock.Header.Number+1, util.Hex(block.Hash().Bytes()))
				}
				// If the channel is not the global channel, it should send a tx to the global channel
				if manager.ID != core.GLOBALCHANNELID {
					tx := core.NewGlobalTx(manager.ID, block.Header.Number, block.Hash())
					// 打印非config通道向global通道中添加的tx信息
					log.Debugf("Channel %s add tx %s to global channel, num: %d", manager.ID, tx.ID, block.Header.Number)
					if err := manager.coordinator.GM.AddTx(tx); err != nil {
						// todo: This is temporary fix
						if err.Error() != "The tx exist in the blockchain aleardy" && raft.GetError(err) != raft.TxInPool {
							log.Fatalf("Channel %s failed to add tx into global channel because %s", manager.ID, err)
							return
						}
					}
				}
				if err := manager.AddBlock(block); err != nil {
					log.Fatalf("Channel %s failed to run because of %s", manager.ID, err)
					return
				}
				log.Debugf("Channel %s has %d block now", manager.ID, block.Header.Number)
				manager.hub.Done(string(block.Header.Number), nil)
				for _, tx := range block.Transactions {
					manager.hub.Done(util.Hex(tx.Hash()), nil)
				}
			}
		case <-manager.stop:
			manager.init = false
			return
		}
	}
}

// Stop stop the manager
func (manager *Manager) Stop() {
	if manager.init {
		manager.stop <- true
		for manager.init {
			time.Sleep(1 * time.Millisecond)
		}
	}
}

// HasGenesisBlock return if the channel has a genesis block
func (manager *Manager) HasGenesisBlock() bool {
	return manager.cm.HasGenesisBlock()
}

// GetBlock return the block of num
func (manager *Manager) GetBlock(num uint64) (*core.Block, error) {
	return manager.cm.GetBlock(num)
}

// AddBlock add a block
func (manager *Manager) AddBlock(block *core.Block) error {
	log.Infof("start adding block %d in channel %v", block.GetNumber(), manager.ID)
	// first update db
	wb := manager.db.NewWriteBatch()
	if err := wb.AddBlock(block); err != nil {
		log.Infof("manager.db.AddBlock error: %s add block %d, %s",
			manager.ID, block.Header.Number, err.Error())
		return err
	}
	if err := manager.cm.AddBlock(block); err != nil {
		log.Infof("manager.cm.AddBlock error: %s add block %d, %s",
			manager.ID, block.Header.Number, err.Error())
		return err
	}

	if isUserChannel(manager.ID) && !isGenesisBlock(block) {
		profile, err := manager.db.GetChannelProfile(manager.ID)
		if err != nil {
			log.Infof("manager.db cannot get channel profile: %s add block %d, %s",
				manager.ID, block.Header.Number, err.Error())
			return err
		}
		if profile.BlockPrice != 0 {
			acc, err := manager.db.GetOrCreateAccount(common.AddressFromChannelID(manager.ID))
			if err != nil {
				return err
			}
			balance := acc.GetBalance()

			storagePrice := uint64(len(block.Bytes())) * profile.BlockPrice
			if balance < storagePrice {
				manager.lock.Lock()
				manager.insufficientBalance = true
				manager.lock.Unlock()
				acc.SubBalance(balance)
				acc.AddDue(storagePrice - balance)
				// zhq todo: am i doing it correct?
			} else {
				acc.SubBalance(storagePrice)
			}
			err = wb.SetAccount(acc)
			if err != nil {
				log.Infof("manager.db cannot set account: %s add block %d, %s",
					manager.ID, block.Header.Number, err.Error())
				return err
			}
		}
	}
	// TODO: Sync too early

	defer func() {
		log.Infof("AddBlock %d in orderer channel %v success", block.GetNumber(), manager.ID)
	}()

	// check is there is any need to update local state of orderer
	switch manager.ID {
	case core.CONFIGCHANNELID:
		if !isGenesisBlock(block) && !manager.coordinator.CanRun(block.Header.ChannelID, block.Header.Number) {
			manager.coordinator.Watch(block.Header.ChannelID, block.Header.Number)
		}
		if err := manager.AddConfigBlock(wb, block); err != nil {
			return err
		}
		return wb.Sync()
	case core.GLOBALCHANNELID:
		wb.Sync()
		return manager.AddGlobalBlock(block)
	case core.ASSETCHANNELID:
		if !isGenesisBlock(block) && !manager.coordinator.CanRun(block.Header.ChannelID, block.Header.Number) {
			manager.coordinator.Watch(block.Header.ChannelID, block.Header.Number)
		}
		if err := manager.AddAssetBlock(wb, block); err != nil {
			return err
		}
		return wb.Sync()
	default:
		wb.Sync()
		return nil
	}
}

// GetBlockSize return the size of blocks
func (manager *Manager) GetBlockSize() uint64 {
	return manager.cm.GetExpect()
}

// AddTx try to add a tx
func (manager *Manager) AddTx(tx *core.Tx) error {
	if manager.db.HasTx(tx) {
		return errors.New("The tx exist in the blockchain aleardy")
	}

	var insufficientBalance bool
	// c.lock.RLock()
	manager.lock.RLock()
	insufficientBalance = manager.insufficientBalance
	// c.lock.RUnlock()
	manager.lock.RUnlock()
	if insufficientBalance {
		return errors.New("Not Enough Balance In User Channel To Generate New Block")
	}

	hash := tx.Hash()
	err := manager.coordinator.Consensus.AddTx(tx)
	if err != nil {
		return err
	}

	// Note: The reason why we must do this is because we must make sure we return the result after we store the block
	// However, we may find a better way to do this if we allow there are more interactive between the consensus and orderer.
	result := manager.hub.Watch(util.Hex(hash), nil)
	if result == nil {
		return nil
	}
	return result.(*event.Result).Err
}

// FetchBlock return the block if exist
func (manager *Manager) FetchBlock(num uint64) (*core.Block, error) {
	return manager.cm.GetBlock(num)
}

// IsMember return if the member belongs to the channel
func (manager *Manager) IsMember(member *core.Member) bool {
	return manager.db.IsMember(manager.ID, member)
}

// IsAdmin return if the member is the admin of the channel
func (manager *Manager) IsAdmin(member *core.Member) bool {
	return manager.db.IsAdmin(manager.ID, member)
}

// IsSystemAdmin return if the member is the system admin
func (manager *Manager) IsSystemAdmin(member *core.Member) bool {
	return manager.db.IsAdmin(core.CONFIGCHANNELID, member)
}

// GetAccount return requested account
func (manager *Manager) GetAccount(address common.Address) (common.Account, error) {
	return manager.db.GetOrCreateAccount(address)
}

// WakeFromSufficientBalance wake up the manager if balance is enough
func (manager *Manager) WakeFromSufficientBalance() {
	manager.lock.Lock()
	manager.insufficientBalance = false
	manager.lock.Unlock()
}

// FetchBlockAsync will fetch book async.
// TODO: fix the thread unsafety
func (manager *Manager) FetchBlockAsync(num uint64) (*core.Block, error) {
	if manager.cm.GetExpect() <= num {
		manager.hub.Watch(string(num), nil)
	}

	block, err := manager.cm.GetBlock(num)
	if err == nil {
		return block, err
	}
	return nil, err
}

// syncBlock is not safe and not efficiency
// todo: should not using a channel to send block
func (manager *Manager) syncBlock() {
	var num = manager.db.GetConsensusBlock(manager.ID)
	for {
		log.Infof("Going to get block %d of channel %s from consensus", num, manager.ID)
		cb, err := manager.coordinator.Consensus.GetBlock(manager.ID, num, true)
		if err != nil {
			log.Infof("Get block %d of channel %s from consensus failed, because %s", num, manager.ID, err.Error())
			continue
		} else {
			log.Infof("Get block %d of channel %s from consensus", num, manager.ID)
		}
		num++
		manager.cbc <- cb
		// Note: solo consensus will create block from block 1 if restart, so we would not remember last consensus block
		if manager.coordinator.Consensus.Info() != "solo" {
			wb := manager.db.NewWriteBatch()
			wb.SetConsensusBlock(manager.ID, num)
			wb.Sync()
		}
	}
}

// getTxsFromConsensusBlock return txs which are legal and duplicate
func (manager *Manager) getTxsFromConsensusBlock(block consensus.Block) (legal, duplicate []*core.Tx) {
	txs := block.GetTxs()
	var count = make(map[string]bool)
	for _, tx := range txs {
		if !util.Contain(count, tx.ID) && !manager.db.HasTx(tx) {
			count[tx.ID] = true
			legal = append(legal, tx)
			// log.Infof("getTxsFromConsensusBlock: block %d in %s add tx %s",
			// 	block.GetNumber(), manager.ID, tx.ID)
		} else {
			duplicate = append(duplicate, tx)
		}
	}
	return
}

func isGenesisBlock(block *core.Block) bool {
	return block.GetNumber() == 0
}

func isUserChannel(id string) bool {
	return id != core.GLOBALCHANNELID && id != core.CONFIGCHANNELID && id != core.ASSETCHANNELID
}
