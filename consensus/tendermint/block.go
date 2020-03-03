package tendermint

import "madledger/core"

// Block is the implementaion of tendermint Block
type Block struct {
	ChannelID string
	Num       uint64
	Txs       [][]byte
}

// GetNumber is the implementation of block
func (block *Block) GetNumber() uint64 {
	return block.Num
}

// GetTxs is the implementation of block
func (block *Block) GetTxs() []*core.Tx {
	// return block.
	var txs []*core.Tx
	for _, txBytes := range block.Txs {
		tx, err := core.BytesToTx(txBytes)
		if err == nil {
			txs = append(txs, tx)
		} else {
			log.Infof("get tx from consensus block failed because %v", err)
		}
	}
	return txs
}
