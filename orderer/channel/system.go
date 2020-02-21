package channel

import (
	"encoding/json"
	"fmt"
	ac "madledger/blockchain/asset"
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
func (manager *Manager) AddAssetBlock(block *core.Block) error {
	if block.Header.Number == 0 {
		return nil
	}

	var err error

	for _, tx := range block.Transactions {
		var payload ac.Payload
		err = json.Unmarshal(tx.Data.Payload, &payload)
		if err != nil {
			return fmt.Errorf("wrong tx format: %v", err)
		}
		sender, err := tx.GetSender()
		if err != nil {
			return fmt.Errorf("wrong sender address %v", sender);
		}
		receiver := tx.GetReceiver()
		//if receiver is not set, issue or transfer money to a channel
		if receiver == common.ZeroAddress {
			receiver = common.BytesToAddress([]byte(payload.ChannelID))
		}

		switch payload.Action {
		case "issue":
			// avoid overflow
			issueValue := tx.Data.Value
			err = manager.issue(tx.Data.Sig.PK, receiver, issueValue)
		case "transfer":
			// if value < 0, sender get money from receiver ??
			transferValue := tx.Data.Value
			err = manager.transfer(sender, receiver, transferValue)
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
	if value == 0 {
		return nil
	}
	receiverAccount, err := manager.db.GetOrCreateAccount(receiver)
	if err != nil {
		return nil
	}
	err = receiverAccount.AddBalance(value)
	if err != nil {
		return nil
	}
	return manager.db.UpdateAccounts(receiverAccount)
}

func (manager *Manager) transfer(sender, receiver common.Address, value uint64) error {
	if value == 0 {
		return nil
	}
	senderAccount, err := manager.db.GetOrCreateAccount(sender)
	if err != nil {
		return err
	}

	if err = senderAccount.SubBalance(value); err != nil {
		return err
	}

	receiverAccount, err := manager.db.GetOrCreateAccount(receiver)
	if err != nil {
		return err
	}

	if err = receiverAccount.AddBalance(value); err != nil {
		return err
	}

	return manager.db.UpdateAccounts(senderAccount, receiverAccount)
}