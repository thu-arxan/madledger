package types

import (
	"bytes"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/util"
)

// Block is the elements of BlockChain
type Block struct {
	// Header of Block
	Header *BlockHeader `json:"Header,omitempty"`
	// Transactions of Block
	Transactions []*Tx `json:"Transactions,omitempty"`
}

// BlockHeader is the header of Block
// Some details still need to be decided
type BlockHeader struct {
	Version   int32  `json:"Version,omitempty"`
	ChannelID string `json:"ChannelID,omitempty"`
	// Bitcoin doesn't contains Number, Ethereum contains it and
	// HyperLedger Fabric calls it as Sequence
	Number uint64 `json:"Number,omitempty"`
	// Hash of the previous block header in the block chain.
	PrevBlock []byte `json:"PrevBlock,omitempty"`
	// Merkle tree reference to hash of all transactions for the block.
	MerkleRoot []byte `json:"MerkleRoot,omitempty"`
	Time       int64  `json:"Time,omitempty"`
}

// NewBlock is the constructor of Block
// TODO: copy the transactions
func NewBlock(channelID string, num uint64, prevHash []byte, txs []*Tx) *Block {
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
	// buffer.Write(util.Int64ToBytes(b.Header.Time))
	return common.BytesToHash(crypto.Hash(buffer.Bytes()))
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
