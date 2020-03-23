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
	"encoding/json"
	"errors"
	"fmt"
	ac "madledger/blockchain/asset"
	cc "madledger/blockchain/config"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/consensus"
	"madledger/core"
	"madledger/orderer/db"
)

// AddConfigBlock add a config block
// The block is formated, so there is no need to verify
func (manager *Manager) AddConfigBlock(block *core.Block) error {
	nums := make(map[string][]uint64)
	if block.Header.Number == 0 {
		return nil
	}
	for _, tx := range block.Transactions {
		var payload cc.Payload
		json.Unmarshal(tx.Data.Payload, &payload)
		var channelID = payload.ChannelID
		// This is a create channel tx,从leveldb中查询是否已经存在channelID
		// 这里并没有对channelID已经存在做出响应,而是在coordinator的createChannel做出响应
		if !manager.db.HasChannel(channelID) {
			// then start the consensus
			err := manager.coordinator.Consensus.AddChannel(channelID, consensus.Config{
				Timeout: manager.coordinator.chainCfg.BatchTimeout,
				MaxSize: manager.coordinator.chainCfg.BatchSize,
				Number:  1,
				Resume:  false,
			})
			channel, err := NewManager(channelID, manager.coordinator)
			if err != nil {
				return err
			}
			// create genesis block here
			// Note: the genesis block will contain no tx
			genesisBlock := core.NewBlock(channelID, 0, core.GenesisBlockPrevHash, []*core.Tx{})

			err = channel.AddBlock(genesisBlock)
			if err != nil {
				return err
			}
			// then start the channel
			go func() {
				log.Infof("system/AddConfigBlock: start channel %s", channelID)
				channel.Start()
			}()
			// 更新coordinator.Managers(map类型)
			manager.coordinator.setChannel(channelID, channel)
			nums[payload.ChannelID] = []uint64{0}
		}
		//todo: ab where to add this kind of tx
		//         may need authentication check
		// this variable only used by orderer to charge block storage
		if payload.BlockPrice != 0 {
			blockPriceBytes := make([]byte, 8)
			binary.BigEndian.PutUint64(blockPriceBytes, payload.BlockPrice)
			err := manager.db.Put([]byte("system@blockprice"), blockPriceBytes)
			if err != nil {
				return nil
			}
		}
		// 更新leveldb
		err := manager.db.UpdateChannel(channelID, payload.Profile)
		if err != nil {
			return err
		}
	}
	manager.coordinator.Unlocks(nums)

	return nil
}

// AddGlobalBlock add a global block
// Note: It should not add block file again.
// TODO: update something in the db
func (manager *Manager) AddGlobalBlock(block *core.Block) error {
	nums := make(map[string][]uint64)
	for _, tx := range block.Transactions {
		payload, err := tx.GetGlobalTxPayload()
		if err != nil {
			return err
		}
		nums[payload.ChannelID] = append(nums[payload.ChannelID], payload.Num)
	}
	manager.coordinator.Unlocks(nums)

	return nil
}

// AddAssetBlock add an asset block
func (manager *Manager) AddAssetBlock(block *core.Block) error {
	if block.Header.Number == 0 {
		return nil
	}
	cache := NewCache(manager.db)
	var err error

	for i, tx := range block.Transactions {
		receiver := tx.GetReceiver()
		status := &db.TxStatus{
			Err:             "",
			BlockNumber:     block.Header.Number,
			BlockIndex:      i,
			Output:          nil,
			ContractAddress: receiver.String(),
		}

		var payload ac.Payload
		err = json.Unmarshal(tx.Data.Payload, &payload)
		if err != nil {
			log.Infof("wrong tx format: %v", err)
			status.Err = err.Error()
			cache.SetTxStatus(tx, status)
			continue
		}
		sender, err := tx.GetSender()
		if err != nil {
			log.Infof("wrong sender address %v", sender)
			status.Err = err.Error()
			cache.SetTxStatus(tx, status)
			continue
		}
		//if receiver is not set, issue or transfer money to a channel
		value := tx.Data.Value
		recipient := payload.Address
		if recipient == common.ZeroAddress {
			recipient = common.BytesToAddress([]byte(payload.ChannelID))
		}
		switch receiver {
		case core.IssueContractAddress:
			err = manager.issue(cache, tx.Data.Sig.PK, tx.Data.Sig.Algo, recipient, value)
		case core.TransferContractrAddress:
			err = manager.transfer(cache, sender, recipient, value)
		case core.TokenExchangeAddress:
			err = manager.exchangeToken(cache, sender, recipient, value)
		default:
			err = errors.New("Contract not support in _asset")
		}

		if err != nil {
			status.Err = err.Error()
			cache.SetTxStatus(tx, status)
			continue
		}
		cache.SetTxStatus(tx, status)
	}
	return cache.Sync()
}


func (manager *Manager) issue(cache Cache, senderPKBytes []byte, pkAlgo crypto.Algorithm, receiver common.Address, value uint64) error {
	pk, err := crypto.NewPublicKey(senderPKBytes, pkAlgo)
	if !cache.IsAssetAdmin(pk) && cache.SetAssetAdmin(pk) != nil {
		return fmt.Errorf("issue authentication failed: %v", err)
	}
	if value == 0 {
		return nil
	}

	receiverAccount, err := cache.GetOrCreateAccount(receiver)
	if err != nil {
		return nil
	}
	err = receiverAccount.AddBalance(value)
	if err != nil {
		return nil
	}
	return cache.UpdateAccounts(receiverAccount)
}

func (manager *Manager) transfer(cache Cache, sender, receiver common.Address, value uint64) error {

	if value == 0 {
		return nil
	}
	senderAccount, err := cache.GetOrCreateAccount(sender)
	if err != nil {
		return err
	}
	if err = senderAccount.SubBalance(value); err != nil {
		return err
	}
	receiverAccount, err := cache.GetOrCreateAccount(receiver)
	if err != nil {
		return err
	}
	if err = receiverAccount.AddBalance(value); err != nil {
		return err
	}

	return cache.UpdateAccounts(senderAccount, receiverAccount)
}

func (manager *Manager) exchangeToken(cache Cache, sender, receiver common.Address, value uint64) error {
	if err := manager.transfer(cache, sender, receiver, value); err != nil {
		return err
	}

	ratioKey := util.BytesCombine(receiver.Bytes(), []byte("ratio"))
	// ratio is set to a value (maybe default by 1) when channel init
	ratioBytes, err := cache.Get(ratioKey, false)
	if err != nil {
		return err
	}

	var ratio uint64
	if ratioBytes != nil {
		ratio = uint64(binary.BigEndian.Uint64(ratioBytes))
	} else {
		ratio = 1
		var ratioVal = make([]byte, 8)
		binary.BigEndian.PutUint64(ratioVal, ratio)
		cache.Put(ratioKey, ratioVal)
	}

	log.Infof("exchangeToken get token / asset ratio %d", ratio)
	var val = make([]byte, 8)
	binary.BigEndian.PutUint64(val, ratio*value)
	cache.Put(util.BytesCombine(receiver.Bytes(), []byte("token"), sender.Bytes()), val)
	return nil
}
