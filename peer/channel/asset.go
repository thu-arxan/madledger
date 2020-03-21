package channel

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	ac "madledger/blockchain/asset"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/core"
	"madledger/peer/db"
)

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
			err = manager.issue(cache, tx.Data.Sig.PK, recipient, value)
		case core.TransferContractrAddress:
			err = manager.transfer(cache, sender, recipient, value)
		case core.TokenExchangeAddress:
			err = manager.exchangeToken(cache, sender, recipient, value)
		default:
			err = errors.New("Contract not support in _asset")
		}

		if err != nil {
			// 如果有错误，那么应该在db里加一条key为txid的错误，如果正确，那么key为txid为ok
			status.Err = err.Error()
			cache.SetTxStatus(tx, status)
			continue
		}
		cache.SetTxStatus(tx, status)
	}
	cache.PutBlock(block)
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

func (manager *Manager) exchangeToken(cache Cache, sender, receiver common.Address, value uint64) error {
	if err := manager.transfer(cache, sender, receiver, value); err != nil {
		return err
	}

	ratioKey := util.BytesCombine(receiver.Bytes(), []byte("ratio"))
	// if ratio not set, default to 1
	ratioBytes, err := cache.Get(ratioKey, true)
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
