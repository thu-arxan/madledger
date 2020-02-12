package raft

import (
	"madledger/consensus"
	"madledger/consensus/raft/eraft"
)

// Config is the config of consensus
type Config struct {
	id    uint64
	dir   string            // root dir for raft storage
	peers map[uint64]string // id => grpc addr
	// consensus config
	cc consensus.Config
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
		id:    id,
		dir:   dir,
		peers: nodes,
		cc:    cc,
		ec:    ec,
	}, nil
}
