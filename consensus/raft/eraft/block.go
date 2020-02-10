package eraft

import (
	"encoding/json"
	"madledger/core"
)

// Block is the implementaion of raft Block
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

// Bytes return bytes of block
func (block *Block) Bytes() []byte {
	bytes, _ := json.Marshal(block)
	return bytes
}

// UnmarshalBlock convert bytes to Block
func UnmarshalBlock(bytes []byte) *Block {
	var block Block
	json.Unmarshal(bytes, &block)
	return &block
}

// // HybridBlock will include all txs of different channels
// type HybridBlock struct {
// 	Num uint64
// 	Txs [][]byte
// }

// // GetNumber return num of block
// func (b *HybridBlock) GetNumber() uint64 {
// 	return b.Num
// }

// // GetTxs return txs of block
// func (b *HybridBlock) GetTxs() [][]byte {
// 	return b.Txs
// }

// // Bytes will return bytes of hybrid block
// func (b *HybridBlock) Bytes() []byte {
// 	bytes, _ := json.Marshal(b)
// 	return bytes
// }

// // UnmarshalHybridBlock convert bytes to HybridBlock
// func UnmarshalHybridBlock(bytes []byte) *HybridBlock {
// 	var block HybridBlock
// 	json.Unmarshal(bytes, &block)
// 	return &block
// }
