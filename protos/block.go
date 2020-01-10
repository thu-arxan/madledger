package protos

import (
	"madledger/common/util"
	"madledger/core/types"
)

// NewBlock is the constructor of Block
// todo: fix sig nil and data nil
func NewBlock(block *types.Block) (*Block, error) {
	var txs = make([]*Tx, len(block.Transactions))
	if block.Transactions == nil {
		txs = nil
	} else {
		for i := range txs {
			tx := block.Transactions[i]
			txs[i] = &Tx{
				ID:   tx.ID,
				Data: NewTxData(&(tx.Data)),
				Time: tx.Time,
			}
		}
	}

	return &Block{
		Header: &BlockHeader{
			Version:    block.Header.Version,
			ChannelID:  block.Header.ChannelID,
			Number:     block.Header.Number,
			PrevBlock:  util.CopyBytes(block.Header.PrevBlock),
			MerkleRoot: util.CopyBytes(block.Header.MerkleRoot),
			Time:       uint64(block.Header.Time),
		},
		Transactions: txs,
	}, nil
}

// ConvertToTypes convert pb.Block to types.Block
// todo: fix nil
func (block *Block) ConvertToTypes() (*types.Block, error) {
	var txs = make([]*types.Tx, len(block.Transactions))
	if len(txs) == 0 {
		txs = nil
	} else {
		for i := range txs {
			tx := block.Transactions[i]
			txs[i] = &types.Tx{
				ID:   tx.ID,
				Time: tx.Time,
			}
			if tx.Data != nil {
				txs[i].Data = *(tx.Data.ToTypes())
			}
		}
	}

	return &types.Block{
		Header: &types.BlockHeader{
			Version:    block.Header.Version,
			ChannelID:  block.Header.ChannelID,
			Number:     block.Header.Number,
			PrevBlock:  block.Header.PrevBlock,
			MerkleRoot: block.Header.MerkleRoot,
			Time:       int64(block.Header.Time),
		},
		Transactions: txs,
	}, nil
}
