package asset

import (
	"madledger/core"
	"encoding/json"
)

func CreateGenesisBlock(payloads []*Payload) (*core.Block, error) {
	var txs []*core.Tx
	for _, payload := range payloads {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		// all zero
		tx := core.NewTxWithoutSig(core.ASSETCHANNELID, payloadBytes, 0)
		txs = append(txs, tx)
	}

	return core.NewBlock(core.ASSETCHANNELID, 0, core.GenesisBlockPrevHash, txs), nil
}