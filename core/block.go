package core

import (
	"bytes"
	"madledger/util"
)

// Block is the elements of BlockChain
type Block struct {
	// Header of Block
	Header BlockHeader
	// Transactions of Block
	// Transactions []*Tx
}

// BlockHeader is the header of Block
// Some details still need to be decided
type BlockHeader struct {
	Version int32
	// Bitcoin doesn't contains Number, Ethereum contains it and
	// HyperLedger Fabric calls it as Sequence
	Number uint64
	// Hash of the previous block header in the block chain.
	PrevBlock []byte
	// Merkle tree reference to hash of all transactions for the block.
	MerkleRoot []byte
	Time       int64
}

// NewBlock is the constructor of Block
// TODO: copy the transactions
func NewBlock(num uint64, prevHash []byte, txs []*Tx) *Block {
	merkleRootHash := CalcMerkleRoot(txs)
	blockHeader := NewBlockHeader(num, prevHash, merkleRootHash)

	block := &Block{
		Header: *blockHeader,
		// Transactions: transactions,
	}
	return block
}

// Hash return the hash of Block
func (b *Block) Hash() []byte {
	var buffer bytes.Buffer
	buffer.Write(util.Int32ToBytes(b.Header.Version))
	buffer.Write(util.Uint64ToBytes(b.Header.Number))
	buffer.Write(b.Header.PrevBlock)
	buffer.Write(b.Header.MerkleRoot)
	buffer.Write(util.Int64ToBytes(b.Header.Time))
	return util.Hash(buffer.Bytes())
}

// NewBlockHeader is the constructor of BlockHeader
// May support the version of others in the future
func NewBlockHeader(num uint64, prevHash, merkleRootHash []byte) *BlockHeader {
	return &BlockHeader{
		Version:    1,
		Number:     num,
		PrevBlock:  prevHash,
		MerkleRoot: merkleRootHash,
		Time:       util.Now(),
	}
}
