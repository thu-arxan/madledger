package channel

import (
	"madledger/consensus"
	"madledger/core"
)

// GetTxsFromConsensusBlock return txs in the consensus block.
func GetTxsFromConsensusBlock(block consensus.Block) []*core.Tx {
	var txs []*core.Tx
	for _, txBytes := range block.GetTxs() {
		tx, err := core.BytesToTx(txBytes)
		if err == nil {
			txs = append(txs, tx)
		} else {
			log.Info(err)
		}
	}
	return txs
}
