package raft

import (
	"fmt"
	pb "madledger/consensus/raft/protos"
)

// Config is the config of eraft
// The raft will use some linear address to make sure the service can config and run simple
// If the chain port is 12345, then raft service will use 12346 and etcd raft node will use 12347
type Config struct {
	id uint64
	// The work path that raft need
	dir     string
	dbDir   string
	walDir  string
	snapDir string
	// peers are eraft address
	peers map[uint64]string
	// The url of node, maybe ip or domain
	url string
	// The port of eraft port
	eraftPort int
	// The port of raft port
	raftPort int
	// The listen address, it should be consensus with the blockchain service
	address string

	snapshotInterval uint64
}

// NewConfig is the constructor of Config
// works on dir and listen on address, id is the id of raft node, nodes is a url map of all nodes
func NewConfig(dir, address string, id uint64, nodes map[uint64]string) (*Config, error) {
	url, eraftPort, err := pb.ParseERaftAddress(nodes[id])
	if err != nil {
		return nil, err
	}

	return &Config{
		id:               id,
		dir:              dir,
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
func (c *Config) GetID() uint64 {
	return c.id
}

func (c *Config) getLocalRaftAddress() string {
	return fmt.Sprintf("%s:%d", c.address, c.raftPort)
}

func (c *Config) getRaftAddress() string {
	return fmt.Sprintf("%s:%d", c.url, c.raftPort)
}

func (c *Config) getLocalERaftAddress() string {
	return fmt.Sprintf("%s:%d", c.address, c.eraftPort)
}

func (c *Config) getERaftAddress() string {
	return fmt.Sprintf("%s:%d", c.url, c.eraftPort)
}

// getPeerAddress return the peer blockchain address
// func (c *Config) getPeerAddress(id uint64) string {
// 	if util.Contain(c.peers, id) {
// 		return pb.ERaftToChain(c.peers[id])
// 	}
// 	return ""
// }
