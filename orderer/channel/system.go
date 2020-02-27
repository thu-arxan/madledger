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
			log.Errorf("wrong tx format: %v", err)
			continue
		}
		sender, err := tx.GetSender()
		if err != nil {
			log.Errorf("wrong sender address %v", sender)
			continue
		}
		receiver := tx.GetReceiver()
		//if receiver is not set, issue or transfer money to a channel
		if receiver == common.ZeroAddress {
			receiver = common.BytesToAddress([]byte(payload.ChannelID))
		}

		log.Infof("receiver is %v", receiver)

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
			// 如果有错误，那么应该在db里加一条key为txid的错误，如果正确，那么key为txid为ok
			log.Infof("shit happened to tx %v : %v", tx, err)
			log.Errorf("err when execute account block tx %v : %v", tx, err)
			continue
		}
		log.Infof("tx %v good", tx)
		manager.db.SetTxExecute(tx.ID)
	}

	return nil
}

func (manager *Manager) issue(senderPKBytes []byte, receiver common.Address, value uint64) error {
	pk, err := crypto.NewPublicKey(senderPKBytes)
	log.Infof("PK is %v", pk)
	if !manager.db.IsAssetAdmin(pk) && manager.db.SetAssetAdmin(pk) != nil {
		return fmt.Errorf("issue authentication failed: %v", err)
	}
	log.Infof("val is %v", value)
	if value == 0 {
		return nil
	}
	// todo:@zhq, should be log.Debugf or remove this when finished.
	log.Infof("rec is %v", receiver)
	receiverAccount, err := manager.db.GetOrCreateAccount(receiver)
	if err != nil {
		return nil
	}
	log.Infof("1rec acc is %v", receiverAccount)
	err = receiverAccount.AddBalance(value)
	if err != nil {
		return nil
	}
	log.Infof("2rec acc is %v", receiverAccount)
	return manager.db.UpdateAccounts(receiverAccount)
}

func (manager *Manager) transfer(sender, receiver common.Address, value uint64) error {
	log.Infof("val is %v", value)

	if value == 0 {
		return nil
	}
	log.Infof("sender is %v", sender)
	senderAccount, err := manager.db.GetOrCreateAccount(sender)
	if err != nil {
		return err
	}
	log.Infof("sender acc is %v", senderAccount)
	if err = senderAccount.SubBalance(value); err != nil {
		return err
	}
	log.Infof("rec is %v", receiver)
	receiverAccount, err := manager.db.GetOrCreateAccount(receiver)
	if err != nil {
		return err
	}
	log.Infof("rec acc is %v", receiverAccount)
	if err = receiverAccount.AddBalance(value); err != nil {
		return err
	}

	return manager.db.UpdateAccounts(senderAccount, receiverAccount)
}
