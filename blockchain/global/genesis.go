package global

import (
	"encoding/json"
	"madledger/core"
)

// CreateGenesisBlock return the genesis block
// maybe the address should be a special addr rather than all zero
// also the data is still need to be discussed
// TODO:
func CreateGenesisBlock(payloads []*Payload) (*core.Block, error) {
	var txs []*core.Tx
	for _, payload := range payloads {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		// all zero
		var addr core.Address
		tx, err := core.NewTx(addr, payloadBytes)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	return core.NewBlock(0, nil, txs), nil
}
