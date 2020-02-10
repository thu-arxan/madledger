package solo

import "madledger/core"

// Block is the implementaion of solo Block
type Block struct {
	channelID string
	num       uint64
	txs       []*core.Tx
}

// GetNumber is the implementation of block
func (block *Block) GetNumber() uint64 {
	return block.num
}

// GetTxs is the implementation of block
func (block *Block) GetTxs() []*core.Tx {
	return block.txs
}
