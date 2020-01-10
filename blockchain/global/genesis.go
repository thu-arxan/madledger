package global

import (
	"encoding/json"
	"madledger/core"
)

// CreateGenesisBlock return the genesis block.
func CreateGenesisBlock(payloads []*Payload) (*core.Block, error) {
	var txs []*core.Tx
	for _, payload := range payloads {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		// all zero
		tx := core.NewTxWithoutSig(core.GLOBALCHANNELID, payloadBytes, 0)
		txs = append(txs, tx)
	}

	return core.NewBlock(core.GLOBALCHANNELID, 0, core.GenesisBlockPrevHash, txs), nil
}
