package raft

import (
	"madledger/consensus"
	"madledger/consensus/raft/eraft"
)

// Config is the config of consensus
type Config struct {
	id    uint64
	dir   string
	dbDir string
	// consensus config
	cc      consensus.Config
	address string
	// eraft config
	ec *eraft.EraftConfig
}

// NewConfig is the constructor of Config
func NewConfig(dir, address string, id uint64, nodes map[uint64]string, join bool, cc consensus.Config) (*Config, error) {
	ec, err := eraft.NewEraftConfig(dir, address, id, nodes, join)
	if err != nil {
		return nil, err
	}

	return &Config{
		address: address,
		dir:     dir,
		cc:      cc,
		ec:      ec,
		id:      id,
	}, nil
}
