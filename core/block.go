// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package core

import (
	"bytes"
	"encoding/json"
	"madledger/common"
	"madledger/common/crypto/hash"
	"madledger/common/util"
)

// Block is the elements of BlockChain
type Block struct {
	// Header of Block
	Header *BlockHeader `json:"header,omitempty"`
	// Transactions of Block
	Transactions []*Tx `json:"transactions,omitempty"`
}

// BlockHeader is the header of Block
// Some details still need to be decided
type BlockHeader struct {
	Version   int32  `json:"version,omitempty"`
	ChannelID string `json:"channelID,omitempty"`
	// Bitcoin doesn't contains Number, Ethereum contains it and
	// HyperLedger Fabric calls it as Sequence
	Number uint64 `json:"number,omitempty"`
	// Hash of the previous block header in the block chain.
	PrevBlock []byte `json:"prevBlock,omitempty"`
	// Merkle tree reference to hash of all transactions for the block.
	MerkleRoot []byte `json:"merkleRoot,omitempty"`
	Time       int64  `json:"time,omitempty"`
}

// NewBlock is the constructor of Block
func NewBlock(channelID string, num uint64, prevHash []byte, txs []*Tx) *Block {
	if len(prevHash) == 0 {
		prevHash = GenesisBlockPrevHash
	}
	merkleRootHash := CalcMerkleRoot(txs)
	blockHeader := NewBlockHeader(channelID, num, prevHash, merkleRootHash)

	block := &Block{
		Header:       blockHeader,
		Transactions: txs,
	}
	return block
}

// Hash return the hash of Block
func (b *Block) Hash() common.Hash {
	var buffer bytes.Buffer
	buffer.Write(util.Int32ToBytes(b.Header.Version))
	buffer.Write(util.Uint64ToBytes(b.Header.Number))
	buffer.Write(b.Header.PrevBlock)
	buffer.Write(b.Header.MerkleRoot)
	// Time should not be included
	// Note: Some consensus has same block time while block has same number,
	// while some consensus may have different block time even block has same number.
	// So we should set block time to same thing if we want support evm timestamp instruction in consensus which
	// block time is not consensused.
	// buffer.Write(util.Int64ToBytes(b.Header.Time))
	return common.BytesToHash(hash.SM3(buffer.Bytes()))
}

// NewBlockHeader is the constructor of BlockHeader
// May support the version of others in the future
func NewBlockHeader(channelID string, num uint64, prevHash, merkleRootHash []byte) *BlockHeader {
	return &BlockHeader{
		Version:    1,
		ChannelID:  channelID,
		Number:     num,
		PrevBlock:  prevHash,
		MerkleRoot: merkleRootHash,
		Time:       util.Now(),
	}
}

// Bytes return the bytes of Block
func (b *Block) Bytes() []byte {
	bytes, _ := json.Marshal(b)
	return bytes
}

// GetNumber return the number of Block
func (b *Block) GetNumber() uint64 {
	return b.Header.Number
}

// GetMerkleRoot return merkle root
func (b *Block) GetMerkleRoot() []byte {
	return b.Header.MerkleRoot
}

// UnmarshalBlock unmarshal json-encoded block
func UnmarshalBlock(data []byte) (*Block, error) {
	var block Block

	if err := json.Unmarshal(data, &block); err != nil {
		return nil, err
	}
	return &block, nil
}
