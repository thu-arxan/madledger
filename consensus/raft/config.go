package raft

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"madledger/common/util"
	"madledger/consensus"
	pb "madledger/consensus/raft/protos"
)

// Config is the config of eraft
// The raft will use some linear address to make sure the service can config and run simple
// If the chain port is 12345, then raft service will use 12346 and etcd raft node will use 12347
type Config struct {
	id   uint64
	join bool
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
	tls     tlsConfig

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
func NewConfig(cfg *consensus.Config) (*Config, error) {
	// chainCfg := cfg.BlockChain

	// if !chainCfg.Raft.Enable {
	// 	return nil, errors.New("The raft is not enable")
	// }

	// dir := fmt.Sprintf("%s/raft", chainCfg.Path)

	// url, eraftPort, err := pb.ParseERaftAddress(chainCfg.Raft.Nodes[chainCfg.Raft.ID])
	// if err != nil {
	// 	return nil, err
	// }

	// var tlsCfg tlsConfig
	// if cfg.TLS.Enable {
	// 	tlsCfg = tlsConfig{
	// 		enable:   cfg.TLS.Enable,
	// 		pool:     cfg.TLS.Pool,
	// 		caFile:   cfg.TLS.CA,
	// 		cert:     *cfg.TLS.Cert,
	// 		certFile: cfg.TLS.RawCert,
	// 		keyFile:  cfg.TLS.Key,
	// 	}
	// } else {
	// 	tlsCfg = tlsConfig{
	// 		enable: false,
	// 	}
	// }

	// return &Config{
	// 	id:               chainCfg.Raft.ID,
	// 	join:             chainCfg.Raft.Join,
	// 	dir:              dir,
	// 	dbDir:            fmt.Sprintf("%s/db", dir),
	// 	walDir:           fmt.Sprintf("%s/wal", dir),
	// 	snapDir:          fmt.Sprintf("%s/snap", dir),
	// 	peers:            cfg.BlockChain.Raft.Nodes,
	// 	url:              url,
	// 	chainPort:        eraftPort - 2,
	// 	address:          cfg.Address,
	// 	tls:              tlsCfg,
	// 	snapshotInterval: 100,
	// }, nil
	return nil, errors.New("Not implementation yet")
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
