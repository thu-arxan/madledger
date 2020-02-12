package channel

import (
	"encoding/json"
	ac "madledger/blockchain/account"
	cc "madledger/blockchain/config"
	"madledger/consensus"
	"madledger/core"
)

// AddConfigBlock add a config block
// The block is formated, so there is no need to verify
func (manager *Manager) AddConfigBlock(block *core.Block) error {
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
			manager.coordinator.Managers[channelID] = channel
		}
		// 更新leveldb
		err := manager.db.UpdateChannel(channelID, payload.Profile)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddGlobalBlock add a global block
// Note: It should not add block file again.
// TODO: update something in the db
func (manager *Manager) AddGlobalBlock(block *core.Block) error {
	return nil
}

// AddAccountBlock add an account block
// TODO: ab
func (manager *Manager) AddAccountBlock(block *core.Block) error {
	if block.Header.Number == 0 {
		return nil
	}

	for _, tx := range block.Transactions {
		var payload ac.Payload
		json.Unmarshal(tx.Data.Payload, &payload)
		var channelID= tx.Data.ChannelID
		var action= payload.Action
		var value= tx.Data.Value
		var err error
		sender, _ := tx.GetSender()
		receiver := tx.GetReceiver()
		if action == "issue" {
			err = manager.db.UpdateAccountIssue(channelID, sender, value)
		} else if action == "transfer" {
			err = manager.db.UpdateAccountTransfer(channelID, sender, receiver, value)
		}
		if err != nil {
			return err
		}

	}

	return nil
}