package consensus

// Block is the block interface
type Block interface {
	GetNumber() uint64
	GetTxs() [][]byte
}
