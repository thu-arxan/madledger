package global

import (
	"encoding/json"
	"madledger/common"
	"madledger/common/util"
	"madledger/core/types"
)

// CreateGenesisBlock return the genesis block
// maybe the address should be a special addr rather than all zero
// also the data is still need to be discussed
// TODO:
func CreateGenesisBlock(payloads []*Payload) (*types.Block, error) {
	var txs []*types.Tx
	for _, payload := range payloads {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		// all zero
		var addr common.Address
		tx := createTx(addr, payloadBytes)
		txs = append(txs, tx)
	}

	return types.NewBlock(types.GLOBALCHANNELID, 0, types.GenesisBlockPrevHash, txs), nil
}

// This function is special prepared for genesis block, because there exists
// no signer and it can create a same block ever a signer exists
func createTx(recipient common.Address, payload []byte) *types.Tx {
	return &types.Tx{
		Data: types.TxData{
			ChannelID:    types.GLOBALCHANNELID,
			AccountNonce: 0,
			Recipient:    recipient.Bytes(),
			Payload:      payload,
			Version:      1,
			Sig:          nil,
		},
		Time: util.Now(),
	}
}
