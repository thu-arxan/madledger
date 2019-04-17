package tendermint

// Block is the implementaion of tendermint Block
type Block struct {
	channelID string
	num       uint64
	txs       [][]byte
}

// GetNumber is the implementation of block
func (block *Block) GetNumber() uint64 {
	return block.num
}

// GetTxs is the implementation of block
func (block *Block) GetTxs() [][]byte {
	return block.txs
}
