package channel

import (
	"madledger/consensus"
	"madledger/core/types"
)

// GetTxsFromConsensusBlock return txs in the consensus block.
func GetTxsFromConsensusBlock(block consensus.Block) []*types.Tx {
	var txs []*types.Tx
	for _, txBytes := range block.GetTxs() {
		tx, err := types.BytesToTx(txBytes)
		if err == nil {
			txs = append(txs, tx)
		} else {
			log.Info(err)
		}
	}
	return txs
}
