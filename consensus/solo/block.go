package solo

// Block is the implementaion of solo Block
type Block struct {
	channelID string
	num       uint64
	txs       [][]byte
}
