package raft

import (
	"fmt"
	"madledger/consensus"
	pb "madledger/consensus/raft/protos"
)

// Config is the config of consensus
type Config struct {
	dir string
	// consensus config
	cc      consensus.Config
	address string
	// eraft config
	ec *EraftConfig
}

// NewConfig is the constructor of Config
func NewConfig(dir, address string, id uint64, nodes map[uint64]string, cc consensus.Config) (*Config, error) {
	ec, err := NewEraftConfig(dir, address, id, nodes)
	if err != nil {
		return nil, err
	}

	return &Config{
		address: address,
		dir:     dir,
		cc:      cc,
		ec:      ec,
	}, nil
}

// EraftConfig is the config of eraft
// The raft will use some linear address to make sure the service can config and run simple
// If the chain port is 12345, then raft service will use 12346 and etcd raft node will use 12347
type EraftConfig struct {
	id      uint64
	dbDir   string
	walDir  string
	snapDir string
	// peers are eraft address
	peers map[uint64]string
	// The url of node, maybe ip or domain
	url string
	// The port of eraft port
	eraftPort int
	// The port of raft port, raftPort = eraftPort + 1
	raftPort int
	// The listen address, it should be consensus with the blockchain service
	address string

	snapshotInterval uint64
}

// NewEraftConfig is the constructor of EraftConfig
// works on dir and listen on address, id is the id of raft node, nodes is a url map of all nodes
func NewEraftConfig(dir, address string, id uint64, nodes map[uint64]string) (*EraftConfig, error) {
	url, eraftPort, err := pb.ParseERaftAddress(nodes[id])
	if err != nil {
		return nil, err
	}

	return &EraftConfig{
		id:               id,
		dbDir:            fmt.Sprintf("%s/db", dir),
		walDir:           fmt.Sprintf("%s/wal", dir),
		snapDir:          fmt.Sprintf("%s/snap", dir),
		peers:            nodes,
		url:              url,
		eraftPort:        eraftPort,
		raftPort:         eraftPort + 1,
		address:          address,
		snapshotInterval: 100,
	}, nil
}

// GetID return the id
func (c *EraftConfig) GetID() uint64 {
	return c.id
}

func (c *EraftConfig) getLocalRaftAddress() string {
	return fmt.Sprintf("%s:%d", c.address, c.raftPort)
}

func (c *EraftConfig) getRaftAddress() string {
	return fmt.Sprintf("%s:%d", c.url, c.raftPort)
}

func (c *EraftConfig) getLocalERaftAddress() string {
	return fmt.Sprintf("%s:%d", c.address, c.eraftPort)
}

func (c *EraftConfig) getERaftAddress() string {
	return fmt.Sprintf("%s:%d", c.url, c.eraftPort)
}

// getPeerAddress return the peer blockchain address
// func (c *Config) getPeerAddress(id uint64) string {
// 	if util.Contain(c.peers, id) {
// 		return pb.ERaftToChain(c.peers[id])
// 	}
// 	return ""
// }
