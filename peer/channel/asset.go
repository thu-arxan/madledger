package channel

import (
	"encoding/json"
	"fmt"
	ac "madledger/blockchain/asset"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/core"
	"madledger/peer/db"
)

//todo: ab
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
		if receiver == core.IssueContractAddress { // issue to channel or person
			if payload.Action == "channel" { // issue to channel
				rec := common.BytesToAddress([]byte(payload.ChannelID))
				err = manager.issue(cache, tx.Data.Sig.PK, rec, value)
			} else if payload.Action == "person" { // issue to person
				err = manager.issue(cache, tx.Data.Sig.PK, payload.Address, value)
			} else { // wrong payload
				status.Err = fmt.Errorf("wrong payload").Error()
				cache.SetTxStatus(tx, status)
				continue
			}
		} else if receiver == core.TransferContractrAddress { // transfer to channel
			err = manager.transfer(cache, sender, payload.Address, value)
		} else { // transfer to person
			err = manager.transfer(cache, sender, receiver, value)
		}

		/*if receiver == common.ZeroAddress {
			receiver = common.BytesToAddress([]byte(payload.ChannelID))
			if receiver == common.ZeroAddress {
				log.Errorf("Not specified receiver")
				cache.SetTxStatus(tx, status)
				continue
			}
		}

		switch payload.Action {
		case "issue":
			// avoid overflow
			issueValue := tx.Data.Value
			err = manager.issue(cache, tx.Data.Sig.PK, receiver, issueValue)
		case "transfer":
			// if value < 0, sender get money from receiver ??
			transferValue := tx.Data.Value
			err = manager.transfer(cache, sender, receiver, transferValue)
		}*/

		if err != nil {
			// 如果有错误，那么应该在db里加一条key为txid的错误，如果正确，那么key为txid为ok
			log.Infof("err when execute account block tx %v : %v", tx, err)
			status.Err = err.Error()
			cache.SetTxStatus(tx, status)
			continue
		}
		cache.SetTxStatus(tx, status)
	}
	cache.Sync()
	return nil
}

func (manager *Manager) issue(cache Cache, senderPKBytes []byte, receiver common.Address, value uint64) error {
	pk, err := crypto.NewPublicKey(senderPKBytes)
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
