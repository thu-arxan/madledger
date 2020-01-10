package protos

import (
	"madledger/common/util"
	"madledger/core"
)

// NewBlock is the constructor of Block
func NewBlock(block *core.Block) (*Block, error) {
	var txs = make([]*Tx, len(block.Transactions))
	if block.Transactions == nil {
		txs = nil
	} else {
		for i := range txs {
			tx, err := NewTx(block.Transactions[i])
			if err != nil {
				return nil, err
			}
			txs[i] = tx
		}
	}

	return &Block{
		Header: &BlockHeader{
			Version:    block.Header.Version,
			ChannelID:  block.Header.ChannelID,
			Number:     block.Header.Number,
			PrevBlock:  util.CopyBytes(block.Header.PrevBlock),
			MerkleRoot: util.CopyBytes(block.Header.MerkleRoot),
			Time:       block.Header.Time,
		},
		Transactions: txs,
	}, nil
}

// ToCore convert pb.Block to core.Block
func (block *Block) ToCore() (*core.Block, error) {
	var txs = make([]*core.Tx, len(block.Transactions))
	if len(txs) == 0 {
		txs = nil
	} else {
		for i := range txs {
			tx, err := block.Transactions[i].ToCore()
			if err != nil {
				return nil, err
			}
			txs[i] = tx
		}
	}

	return &core.Block{
		Header: &core.BlockHeader{
			Version:    block.Header.Version,
			ChannelID:  block.Header.ChannelID,
			Number:     block.Header.Number,
			PrevBlock:  util.CopyBytes(block.Header.PrevBlock),
			MerkleRoot: util.CopyBytes(block.Header.MerkleRoot),
			Time:       block.Header.Time,
		},
		Transactions: txs,
	}, nil
}
