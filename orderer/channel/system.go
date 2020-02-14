package channel

import (
	"encoding/json"
	"fmt"
	"errors"
	ac "madledger/blockchain/account"
	cc "madledger/blockchain/config"
	"madledger/common"
	"madledger/common/crypto"
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

	var err error

	for _, tx := range block.Transactions {
		var payload ac.Payload
		err = json.Unmarshal(tx.Data.Payload, &payload)
		if err != nil {
			fmt.Errorf("wrong tx format: %v", err)
		}
		sender, err := tx.GetSender()
		if err != nil {
			return fmt.Errorf("wrong sender address %v", sender);
		}
		receiver := tx.GetReceiver()
		//if receiver is not set, issue or transfer money to a channel
		if receiver == common.ZeroAddress {
			receiver = common.BytesToAddress([]byte(tx.Data.ChannelID))
		}

		switch payload.Action {
		case "issue":
			// avoid overflow
			issueValue := tx.Data.Value
			if issueValue < 0 {
				issueValue = 0
			}
			err = manager.issue(tx.Data.Sig.PK, receiver, uint64(issueValue))
		case "transfer":
			err = manager.db.UpdateAccount(nil)

		}
		if err != nil {
			return fmt.Errorf("err when execute account block tx %v : %v", tx, err)
		}
	}

	return nil
}

func (manager *Manager) issue(senderPKBytes []byte, receiver common.Address, value uint64) error {
	pk, err := crypto.NewPublicKey(senderPKBytes)
	if !manager.db.IsAccountAdmin(pk) && manager.db.SetAccountAdmin(pk) != nil {
		return fmt.Errorf("issue authentication failed: %v", err)
	}
	receiverAccount, err := manager.db.GetAccount(receiver)
	if err != nil {
		return nil
	}
	return receiverAccount.AddBalance(value)
}