package raft

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"madledger/common/util"
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
	// The port of block chain, and the port of raft service and hashicorp raft can be estimated
	chainPort int
	// The listen address, it should be consensus with the blockchain service
	address string

	snapshotInterval uint64
}

// tlsConfig is the tls config part
type tlsConfig struct {
	enable bool
	// Pool of CA
	pool     *x509.CertPool
	caFile   string
	cert     tls.Certificate
	certFile string
	keyFile  string
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
		chainPort:        eraftPort - 2,
		address:          address,
		snapshotInterval: 100,
	}, nil
}

// GetID return the id
func (c *Config) GetID() uint64 {
	return c.id
}

func (c *Config) getLocalRaftAddress() string {
	return fmt.Sprintf("%s:%d", c.address, c.chainPort+1)
}

func (c *Config) getRaftAddress() string {
	return fmt.Sprintf("%s:%d", c.url, c.chainPort+1)
}

func (c *Config) getLocalERaftAddress() string {
	return fmt.Sprintf("%s:%d", c.address, c.chainPort+2)
}

func (c *Config) getERaftAddress() string {
	return fmt.Sprintf("%s:%d", c.url, c.chainPort+2)
}

// getPeerAddress return the peer blockchain address
func (c *Config) getPeerAddress(id uint64) string {
	if util.Contain(c.peers, id) {
		return pb.ERaftToChain(c.peers[id])
	}
	return ""
}
