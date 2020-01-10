package core

import (
	"encoding/json"
	"errors"
	"madledger/common"
	"madledger/common/util"
)

// GlobalTxPayload is the payload of global tx
type GlobalTxPayload struct {
	ChannelID string
	Num       uint64
	Hash      common.Hash
}

// NewGlobalTx return a standard global tx
func NewGlobalTx(channelID string, num uint64, hash common.Hash) *Tx {
	var payload = GlobalTxPayload{
		ChannelID: channelID,
		Num:       num,
		Hash:      hash,
	}
	payloadBytes, _ := json.Marshal(payload)
	var tx = &Tx{
		Data: TxData{
			ChannelID: GLOBALCHANNELID,
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
func (tx *Tx) GetGlobalTxPayload() (*GlobalTxPayload, error) {
	if tx.Data.ChannelID != GLOBALCHANNELID {
		return nil, errors.New("The tx does not belog to global channel")
	}
	var payload GlobalTxPayload
	err := json.Unmarshal(tx.Data.Payload, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
