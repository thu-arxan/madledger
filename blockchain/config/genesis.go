package config

import (
	"encoding/json"
	"madledger/common"
	"madledger/common/util"
	"madledger/core/types"
)

// CreateGenesisBlock return the genesis block
// maybe the address should be a special addr rather than all zero
// also the data is still need to be discussed
// TODO: many things
func CreateGenesisBlock() (*types.Block, error) {
	var payloads = []Payload{Payload{
		ChannelID: types.CONFIGCHANNELID,
		Profile: Profile{
			Public: true,
		},
		Version: 1,
	}, Payload{
		ChannelID: types.GLOBALCHANNELID,
		Profile: Profile{
			Public: true,
		},
		Version: 1,
	}}
	var txs []*types.Tx
	for _, payload := range payloads {
		payloadBytes, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}
		// all zero
		var addr common.Address
		tx := createTx(addr, payloadBytes)
		txs = append(txs, tx)
	}

	return types.NewBlock(types.CONFIGCHANNELID, 0, nil, txs), nil
}

// This function is special prepared for genesis block, because there exists
// no signer and it can create a same block ever a signer exists
func createTx(recipient common.Address, payload []byte) *types.Tx {
	return &types.Tx{
		Data: &types.TxData{
			ChannelID:    types.CONFIGCHANNELID,
			AccountNonce: 0,
			Recipient:    recipient,
			Payload:      payload,
			Version:      1,
			Sig:          nil,
		},
		Time: util.Now(),
	}
}
