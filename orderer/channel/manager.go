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
	"encoding/binary"
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
}

// NewManager is the constructor of Manager
func NewManager(id string, coordinator *Coordinator) (*Manager, error) {
	cm, err := blockchain.NewManager(id, fmt.Sprintf("%s/%s", coordinator.chainCfg.Path, id))
	if err != nil {
		return nil, err
	}
	return &Manager{
		ID:          id,
		db:          coordinator.db,
		cm:          cm,
		cbc:         make(chan consensus.Block, 1024),
		init:        false,
		stop:        make(chan bool),
		hub:         event.NewHub(),
		coordinator: coordinator,
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
	var price uint64
	BlockPriceBytes, err := manager.db.Get([]byte("system@blockprice"), true)
	if err != nil {
		log.Infof("manager.db not available: %s add block %d, %s",
			manager.ID, block.Header.Number, err.Error())
		return err
	}
	if BlockPriceBytes != nil && isUserChannel(manager.ID) && !isGenesisBlock(block) {
		acc, err := manager.db.GetOrCreateAccount(common.BytesToAddress([]byte(manager.ID)))
		if err != nil {
			return err
		}
		left := acc.GetBalance()
		BlockPrice := uint64(binary.BigEndian.Uint64(BlockPriceBytes))

		price = uint64(len(block.Bytes())) * BlockPrice
		if left < price {
			errMsg := fmt.Sprintf("insuffuicient balance in channel %s", manager.ID)
			return errors.New(errMsg)
		}
	}
	// first update db
	if err := manager.db.AddBlock(block); err != nil {
		log.Infof("manager.db.AddBlock error: %s add block %d, %s",
			manager.ID, block.Header.Number, err.Error())
		return err
	}
	if err := manager.cm.AddBlock(block); err != nil {
		log.Infof("manager.cm.AddBlock error: %s add block %d, %s",
			manager.ID, block.Header.Number, err.Error())
		return err
	}

	//  after adding block, sub the channel balance
	if price != 0 && isUserChannel(manager.ID) && !isGenesisBlock(block) {
		if err := manager.subChannelAsset(manager.ID, price); err != nil {
			return err
		}
	}
	defer func() {
		log.Infof("AddBlock %d in orderer channel %v success", block.GetNumber(), manager.ID)
	}()

	// check is there is any need to update local state of orderer
	switch manager.ID {
	case core.CONFIGCHANNELID:
		if !isGenesisBlock(block) && !manager.coordinator.CanRun(block.Header.ChannelID, block.Header.Number) {
			manager.coordinator.Watch(block.Header.ChannelID, block.Header.Number)
		}
		return manager.AddConfigBlock(block)
	case core.GLOBALCHANNELID:
		return manager.AddGlobalBlock(block)
	case core.ASSETCHANNELID:
		if !isGenesisBlock(block) && !manager.coordinator.CanRun(block.Header.ChannelID, block.Header.Number) {
			manager.coordinator.Watch(block.Header.ChannelID, block.Header.Number)
		}
		return manager.AddAssetBlock(block)
	default:
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
	return manager.db.IsSystemAdmin(member)
}

// GetAccount return requested account
func (manager *Manager) GetAccount(address common.Address) (common.Account, error) {
	return manager.db.GetOrCreateAccount(address)
}

// GetTxStatus return request txStatus
func (manager *Manager) GetTxStatus(channelID, txID string) (*db.TxStatus, error) {
	return manager.db.GetTxStatus(channelID, txID)
}

// GetMaxGas used by userchannel
func (manager *Manager) GetMaxGas() uint64 {
	// this param should be set when every user channel is created
	maxGasBytes, err := manager.db.Get(util.BytesCombine([]byte(manager.ID), []byte("maxGas")), false)
	if err != nil {
		return 0
	}
	return binary.BigEndian.Uint64(maxGasBytes)
}

// GetGasPrice used by userchannel
func (manager *Manager) GetGasPrice() uint64 {
	gasPriceBytes, err := manager.db.Get(util.BytesCombine([]byte(manager.ID), []byte("gasPrice")), false)
	if err != nil {
		return 0
	}
	return binary.BigEndian.Uint64(gasPriceBytes)
}

// GetAssetTokenRatio used by userchannel
func (manager *Manager) GetAssetTokenRatio() uint64 {
	ratioBytes, err := manager.db.Get(util.BytesCombine([]byte(manager.ID), []byte("ratio")), false)
	if err != nil {
		return 0
	}
	return binary.BigEndian.Uint64(ratioBytes)
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
// todo: the manager should not begin from 1 and should not using a channel to send block
func (manager *Manager) syncBlock() {
	var num uint64 = 1
	for {
		log.Infof("Going to get block %d of channel %s from consensus", num, manager.ID)
		cb, err := manager.coordinator.Consensus.GetBlock(manager.ID, num, true)
		if err != nil {
			log.Infof("Get block %d of channel %s from consensus failed, because %s", num, manager.ID, err.Error())
			//fmt.Println(err)
			continue
		} else {
			log.Infof("Get block %d of channel %s from consensus", num, manager.ID)
		}
		num++
		manager.cbc <- cb
	}
}

// getTxsFromConsensusBlock return txs which are legal and duplicate
func (manager *Manager) getTxsFromConsensusBlock(block consensus.Block) (legal, duplicate []*core.Tx) {
	txs := GetTxsFromConsensusBlock(block)
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

func (manager *Manager) subChannelAsset(id string, price uint64) error {
	wb := manager.db.NewWriteBatch()
	if !isUserChannel(manager.ID) {
		return nil
	}
	acc, err := manager.db.GetOrCreateAccount(common.BytesToAddress([]byte(id)))
	if err != nil {
		return err
	}
	if err := acc.SubBalance(price); err != nil {
		return err
	}
	if err := wb.UpdateAccounts(acc); err != nil {
		return err
	}
	return wb.Sync()
}

func isUserChannel(id string) bool {
	return id != core.GLOBALCHANNELID && id != core.CONFIGCHANNELID && id != core.ASSETCHANNELID
}
