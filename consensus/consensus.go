package consensus

import "madledger/core"

// Consensus is the interface
type Consensus interface {
	// Start Consensus service
	Start() error
	// Stop stop the consensus
	Stop() error
	// Add a tx
	AddTx(tx *core.Tx) error
	// AddChannel will add a consensus of a channel
	AddChannel(channelID string, cfg Config) error
	// GetBlock return the block or error right away if async is false, else return block until the block is created.
	GetBlock(channelID string, num uint64, async bool) (Block, error)
}
