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
	"reflect"
)

// AddAssetBlock add an asset block
func (manager *Manager) AddAssetBlock(block *core.Block) error {

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
			recipient = common.AddressFromChannelID(payload.ChannelID)
		}
		switch receiver {
		case core.IssueContractAddress:
			err = manager.issue(cache, tx.Data.Sig.PK, tx.Data.Sig.Algo, recipient, value)
		case core.TransferContractrAddress:
			err = manager.transfer(cache, sender, recipient, value)
		case core.TokenExchangeAddress:
			err = manager.exchangeToken(cache, sender, recipient, value, payload.ChannelID)
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

func (manager *Manager) issue(cache Cache, senderPKBytes []byte, pkAlgo crypto.Algorithm, receiver common.Address, value uint64) error {
	pk, err := crypto.NewPublicKey(senderPKBytes, pkAlgo)
	if !cache.IsAssetAdmin(pk, pkAlgo) && cache.SetAssetAdmin(pk, pkAlgo) != nil {
		return fmt.Errorf("issue authentication failed: %v", err)
	}
	if value == 0 {
		return nil
	}

	receiverAccount, err := cache.GetOrCreateAccount(receiver)
	if err != nil {
		return nil
	}
	valueLeft, err := manager.payDue(receiverAccount, value)
	if err != nil {
		return err
	}

	err = receiverAccount.AddBalance(valueLeft)
	if err != nil {
		return err
	}
	return cache.UpdateAccounts(receiverAccount)
}

func (manager *Manager) transfer(cache Cache, sender, receiver common.Address, value uint64) error {

	if value == 0 || reflect.DeepEqual(sender, receiver) {
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
	valueLeft, err := manager.payDue(receiverAccount, value)
	if err != nil {
		return err
	}

	if err = receiverAccount.AddBalance(valueLeft); err != nil {
		return err
	}

	return cache.UpdateAccounts(senderAccount, receiverAccount)
}

func (manager *Manager) exchangeToken(cache Cache, sender, receiver common.Address, value uint64, channelID string) error {
	if err := manager.transfer(cache, sender, receiver, value); err != nil {
		return err
	}

	profile, err := manager.db.GetChannelProfile(channelID)
	if err != nil {
		return nil
	}

	tokenKey := util.BytesCombine(receiver.Bytes(), []byte("token"), sender.Bytes())
	tokenBytes, err := cache.Get(tokenKey, true)
	var token uint64
	if tokenBytes != nil {
		token = uint64(binary.BigEndian.Uint64(tokenBytes))
	}
	token += profile.AssetTokenRatio * value
	val := make([]byte, 8)
	binary.BigEndian.PutUint64(val, token)

	cache.Put(tokenKey, val)
	log.Infof("exchange token completed. token left: %v", val)
	return nil
}

func (manager *Manager) payDue(acc common.Account, value uint64) (uint64, error) {
	due := acc.GetDue()
	if due == 0 {
		return value, nil
	}
	if value < due {
		return 0, acc.SubDue(value)
	}
	return value - due, acc.SubDue(due)
}
