package channel

import (
	"encoding/json"
	"madledger/consensus"
	"madledger/core/types"
)

//GetTxsFromConsensusBlock return txs in the consensus block
func GetTxsFromConsensusBlock(block consensus.Block) []*types.Tx {
	var txs []*types.Tx
	for _, txBytes := range block.GetTxs() {
		var tx types.Tx
		err := json.Unmarshal(txBytes, &tx)
		if err == nil {
			// return nil, err
			txs = append(txs, &tx)
		}
	}
	return txs
}
