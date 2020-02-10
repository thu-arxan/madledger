package core

//todo: ab need to add account tx?

import (
	"encoding/json"
	"errors"
	"madledger/common"
	"madledger/common/util"
)

// AccountTxPayload is the payload of account tx
type AccountTxPayload struct {
	ChannelID string
}

// NewAccountTx return a standard account tx
func NewAccountTx(channelID string, num uint64, hash common.Hash) *Tx {
	var payload = AccountTxPayload{
		ChannelID: channelID,
	}
	payloadBytes, _ := json.Marshal(payload)
	var tx = &Tx{
		Data: TxData{
			ChannelID: ACCOUNTCHANNELID,
			Nonce:     0,
			Recipient: common.ZeroAddress.Bytes(),
			Payload:   payloadBytes,
			Version:   1,
		},
		Time: util.Now(),
	}
	tx.ID = util.Hex(tx.Hash())
	return tx
}

// GetGlobalTxPayload return the payload of tx
func (tx *Tx) GetAccountTxPayload() (*AccountTxPayload, error) {
	if tx.Data.ChannelID != ACCOUNTCHANNELID {
		return nil, errors.New("The tx does not belog to global channel")
	}
	var payload AccountTxPayload
	err := json.Unmarshal(tx.Data.Payload, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
