package tendermint

import (
	"fmt"
	"madledger/consensus"
)

// Note: We just reach consensus in one chain of tendermint now.
// So we will not distinguish channel in the consensus layer and we consensus on
// one block then we split it into different blocks of different channels.

// Consensus is the implementation of tendermint
type Consensus struct {
	app  *Glue
	node *Node
}

// NewConsensus is the constructor of tendermint.Consensus
// Note: We are not going to support different configs of different channels.
// TODO: Not finished yet
func NewConsensus(channels map[string]consensus.Config, cfg *Config) (consensus.Consensus, error) {
	app, err := NewGlue(fmt.Sprintf("%s/.glue", cfg.Dir), cfg.Port.App)
	if err != nil {
		return nil, err
	}
	node, err := NewNode(cfg, app)
	if err != nil {
		return nil, err
	}
	return &Consensus{
		app:  app,
		node: node,
	}, nil
}

// Start is the implementation of interface
func (c *Consensus) Start() error {
	return nil
}

// AddChannel add a channel
// Because we do not care the channel in the consensus layer, so the function will do nothing now.
func (c *Consensus) AddChannel(channelID string, cfg consensus.Config) error {
	return nil
}

// AddTx is the implementation of interface
func (c *Consensus) AddTx(channelID string, tx []byte) error {
	return c.app.AddTx(channelID, tx)
}

// SyncBlocks is the implementation of interface
func (c *Consensus) SyncBlocks(channelID string, ch *chan consensus.Block) error {
	return nil
}

// GetNumber is the implementation of interface
func (c *Consensus) GetNumber(channelID string) (uint64, error) {
	return 0, nil
}

// Stop is the implementation of interface
func (c *Consensus) Stop() error {
	return nil
}
