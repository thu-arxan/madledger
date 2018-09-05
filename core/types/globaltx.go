package types

import (
	"encoding/json"
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
			ChannelID:    GLOBALCHANNELID,
			AccountNonce: 0,
			Recipient:    common.ZeroAddress.Bytes(),
			Payload:      payloadBytes,
			Version:      1,
			Sig:          nil,
		},
		Time: util.Now(),
	}
	return tx
}
