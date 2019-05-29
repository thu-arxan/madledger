package raft

import (
	"errors"
	"madledger/consensus"
)

// Consensus is the implementation of interface
// Consensus keeps connections to raft services
type Consensus struct {
	cfg *Config
}

// NewConseneus is the constructor of Consensus
func NewConseneus(cfg *Config) *Consensus {
	return &Consensus{
		cfg: cfg,
	}
}

// AddTx is the implementation of interface
func (c *Consensus) AddTx(channelID string, tx []byte) error {
	return errors.New("Not implementation yet")
}

// AddChannel is the implementation of interface
// Note: we can ignore this function now
func (c *Consensus) AddChannel(channelID string, cfg Config) error {
	return nil
}

// GetBlock is the implementation of interface
func (c *Consensus) GetBlock(channelID string, num uint64, async bool) (consensus.Block, error) {
	return nil, errors.New("Not implementation yet")
}
