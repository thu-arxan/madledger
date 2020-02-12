package channel

import (
	"madledger/consensus"
	"madledger/core"
)

// GetTxsFromConsensusBlock return txs in the consensus block.
func GetTxsFromConsensusBlock(block consensus.Block) []*core.Tx {
	// var txs []*core.Tx
	// for _, txBytes := range block.GetTxs() {
	// 	tx, err := core.BytesToTx(txBytes)
	// 	if err == nil {
	// 		txs = append(txs, tx)
	// 	} else {
	// 		log.Infof("get tx from consensus block failed because %v", err)
	// 	}
	// }
	// return txs
	return block.GetTxs()
}
