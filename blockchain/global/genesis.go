package global

import (
	"encoding/json"
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
		tx := types.NewTxWithoutSig(types.GLOBALCHANNELID, payloadBytes, 0)
		txs = append(txs, tx)
	}

	return types.NewBlock(types.GLOBALCHANNELID, 0, types.GenesisBlockPrevHash, txs), nil
}
