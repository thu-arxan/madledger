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
	"encoding/json"
	"errors"
	"fmt"
	ac "madledger/blockchain/asset"
	cc "madledger/blockchain/config"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/consensus"
	"madledger/core"
	"reflect"
)

// AddConfigBlock add a config block
// The block is formated, so there is no need to verify
func (manager *Manager) AddConfigBlock(block *core.Block) error {
	nums := make(map[string][]uint64)
	if block.Header.Number == 0 {
		return nil
	}
	for _, tx := range block.Transactions {
		// this kind of tx is about consensus configuration change
		// will have different kind of payload
		if txType, err := core.GetTxType(common.BytesToAddress(tx.Data.Recipient).String()); err == nil && txType == core.CONSENSUS {
			continue
		}
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
		//todo: ab update channel may modify blockPrice of user channel
		//         may need authentication check
		// also should this use write batch?
		err := manager.db.UpdateChannel(channelID, payload.Profile)
		if err != nil {
			return err
		}
	}
	manager.coordinator.Unlocks(nums)

	return nil
}

func (manager *Manager) AddGlobalBlock(block *core.Block) error {
	nums := make(map[string][]uint64)
	for _, tx := range block.Transactions {
		payload, err := tx.GetGlobalTxPayload()
		if err != nil {
			return err
		}
		switch payload.ChannelID {
		case core.CONFIGCHANNELID, core.ASSETCHANNELID:
			manager.coordinator.Unlocks(map[string][]uint64{payload.ChannelID: []uint64{payload.Num}})
		default:
			nums[payload.ChannelID] = append(nums[payload.ChannelID], payload.Num)
		}
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

	for _, tx := range block.Transactions {
		receiver := tx.GetReceiver()
		var payload ac.Payload
		err = json.Unmarshal(tx.Data.Payload, &payload)
		if err != nil {
			log.Infof("wrong tx format: %v", err)
			continue
		}
		sender, err := tx.GetSender()
		if err != nil {
			log.Infof("wrong sender address %v", sender)
			continue
		}
		//if receiver is not set, issue or transfer money to a channel
		value := tx.Data.Value
		recipient := payload.Address
		if recipient == common.ZeroAddress {
			recipient = common.AddressFromChannelID(payload.ChannelID)
		}
		switch receiver {
		case core.IssueContractAddress:
			err = manager.issue(cache, tx.Data.Sig.PK, tx.Data.Sig.Algo, recipient, value, payload.ChannelID)
		case core.TransferContractrAddress:
			err = manager.transfer(cache, sender, recipient, value, payload.ChannelID)
		case core.TokenExchangeAddress:
			err = manager.exchangeToken(cache, sender, recipient, value, payload.ChannelID)
		default:
			err = errors.New("Contract not support in _asset")
		}
	}
	return cache.Sync()
}

func (manager *Manager) issue(cache Cache, senderPKBytes []byte, pkAlgo crypto.Algorithm, receiver common.Address, value uint64, channelID string) error {
	pk, err := crypto.NewPublicKey(senderPKBytes, pkAlgo)
	if !cache.IsAssetAdmin(pk, pkAlgo) && cache.SetAssetAdmin(pk, pkAlgo) != nil {
		return fmt.Errorf("issue authentication failed: %v", err)
	}
	if value == 0 {
		return nil
	}

	receiverAccount, err := cache.GetOrCreateAccount(receiver)
	if err != nil {
		return err
	}

	valueLeft, err := manager.payDueAndTryWakeChannel(receiverAccount, value, channelID)
	if err != nil {
		return err
	}

	err = receiverAccount.AddBalance(valueLeft)
	if err != nil {
		return err
	}
	return cache.UpdateAccounts(receiverAccount)
}

func (manager *Manager) transfer(cache Cache, sender, receiver common.Address, value uint64, channelID string) error {

	if value == 0 || reflect.DeepEqual(sender, receiver){
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

	valueLeft, err := manager.payDueAndTryWakeChannel(receiverAccount, value, channelID)
	if err != nil {
		return err
	}

	if err = receiverAccount.AddBalance(valueLeft); err != nil {
		return err
	}

	return cache.UpdateAccounts(senderAccount, receiverAccount)
}

func (manager *Manager) exchangeToken(cache Cache, sender, receiver common.Address, value uint64, channelID string) error {
	if err := manager.transfer(cache, sender, receiver, value, channelID); err != nil {
		return err
	}

	// orderer can't modify token because they don't know the exact value of token in every channel!!
	return nil
}

func (manager *Manager) payDueAndTryWakeChannel(acc common.Account, value uint64, channelID string) (uint64, error) {
	due := acc.GetDue()
	if due == 0 {
		return value, nil
	}
	if value < due {
		return 0, acc.SubDue(value)
	}
	if err := manager.coordinator.WakeDueChannel(channelID); err != nil {
		log.Warnf("channel awake error: %v", err)
	}
	log.Infof("wake channel %v", channelID)
	return value - due, acc.SubDue(due)
}
