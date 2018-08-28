package global

import (
	"encoding/json"
	"madledger/common"
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
		tx, err := types.NewTx(addr, payloadBytes)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	return types.NewBlock(0, nil, txs), nil
}
