package consensus

// Consensus is the interface
type Consensus interface {
	// Start Consensus service
	Start() error
	// Stop stop the consensus
	Stop() error
	// Add a tx, tx is bytes so the consensus does not care what tx is
	AddTx(channelID string, tx []byte) error
	// AddChannel will add a consensus of a channel
	AddChannel(channelID string, cfg Config) error
	// GetBlock return the block or error right away if async is false, else return block until the block is created.
	GetBlock(channelID string, num uint64, async bool) (Block, error)
}
