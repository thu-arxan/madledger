package config

import (
	"encoding/json"
	"madledger/core/types"
)

// CreateGenesisBlock return the genesis block
// TODO: maybe there should includes some admins in the genesis block
func CreateGenesisBlock() (*types.Block, error) {
	var payloads = []Payload{Payload{
		ChannelID: types.CONFIGCHANNELID,
		Profile: &Profile{
			Public: true,
		},
		Version: 1,
	}, Payload{
		ChannelID: types.GLOBALCHANNELID,
		Profile: &Profile{
			Public: true,
		},
		Version: 1,
	}}
	var txs []*types.Tx
	for i, payload := range payloads {
		payloadBytes, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}

		accountNonce := uint64(i)
		tx := types.NewTxWithoutSig(types.CONFIGCHANNELID, payloadBytes, accountNonce)
		txs = append(txs, tx)
	}

	return types.NewBlock(types.CONFIGCHANNELID, 0, types.GenesisBlockPrevHash, txs), nil
}
