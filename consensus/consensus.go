package consensus

// Consensus is the interface
type Consensus interface {
	// Start Consensus service
	Start() error
	// Add a tx, tx is bytes so the consensus does not care what tx is
	AddTx(channelID string, tx []byte) error
	// SyncBlocks provide a way for channel manager sync blocks
	SyncBlocks(channelID string, ch *chan Block) error
	// GetNumber helps shorten the gap between orderer and consensus, if not support return 0 is recommend
	GetNumber(channelID string) (uint64, error)
	// Stop stop the consensus
	Stop() error
	// AddChannel will add a consensus of a channel
	AddChannel(channelID string, cfg Config) error
}
