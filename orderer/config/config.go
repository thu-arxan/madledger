// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"madledger/common/util"
	"os"
	"regexp"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

var (
	gopath = os.Getenv("GOPATH")
)

// Config is the combination of all config
type Config struct {
	Port       int              `yaml:"Port"`
	Address    string           `yaml:"Address"`
	Debug      bool             `yaml:"Debug"`
	TLS        TLSConfig        `yaml:"TLS"`
	BlockChain BlockChainConfig `yaml:"BlockChain"`
	Consensus  struct {
		Type       string           `yaml:"Type"`
		Tendermint TendermintConfig `yaml:"Tendermint"`
		Raft       RaftConfig       `yaml:"Raft"`
	} `yaml:"Consensus"`
	DB struct {
		Type    string `yaml:"Type"`
		LevelDB struct {
			Path string `yaml:"Path"`
		} `yaml:"LevelDB"`
	} `yaml:"DB"`
}

// BlockChainConfig is the config of blockchain
type BlockChainConfig struct {
	BatchTimeout int    `yaml:"BatchTimeout"`
	BatchSize    int    `yaml:"BatchSize"`
	Path         string `yaml:"Path"`
	Verify       bool   `yaml:"Verify"`
}

// LoadConfig load config from the config file
func LoadConfig(cfgFile string) (*Config, error) {
	cfgBytes, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	var cfg = new(Config)
	err = yaml.Unmarshal(cfgBytes, cfg)
	if err != nil {
		return nil, err
	}
	if err := cfg.checkAndSetServerConfig(); err != nil {
		return nil, err
	}
	if err := cfg.checkBlockChainConfig(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// checkAndSetServerConfig check config of server and set config of tls
func (cfg *Config) checkAndSetServerConfig() error {
	if cfg.Port < 1024 {
		return fmt.Errorf("The port can not be %d", cfg.Port)
	}
	if cfg.Address == "" {
		return errors.New("The address can not be empty")
	}
	if err := cfg.checkAndSetTLSConfig(); err != nil {
		return err
	}
	return nil
}

// checkAndSetTLSConfig check the tls config and set necessary things
func (cfg *Config) checkAndSetTLSConfig() error {
	if cfg.TLS.Enable {
		if cfg.TLS.CA == "" {
			return errors.New("The CA can not be empty")
		}
		if cfg.TLS.RawCert == "" {
			return errors.New("The cert can not be empty")
		}
		if cfg.TLS.Key == "" {
			return errors.New("The key can not be empty")
		}
		// load pool
		pool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(cfg.TLS.CA)
		if err != nil {
			return err
		}
		ok := pool.AppendCertsFromPEM(ca)
		if !ok {
			return fmt.Errorf("Failed to load ca file: %s", cfg.TLS.CA)
		}
		// load cert
		cert, err := tls.LoadX509KeyPair(cfg.TLS.RawCert, cfg.TLS.Key)
		if err != nil {
			return err
		}
		cfg.TLS.Pool = pool
		cfg.TLS.Cert = &cert
	}
	return nil
}

func (cfg *Config) checkBlockChainConfig() error {
	if cfg.BlockChain.Path == "" {
		return errors.New("The path of blockchain is not provided")
	}
	if cfg.BlockChain.BatchTimeout <= 0 {
		return fmt.Errorf("The batch timeout can not be %d", cfg.BlockChain.BatchTimeout)
	}
	if cfg.BlockChain.BatchSize <= 0 {
		return fmt.Errorf("The batch size can not be %d", cfg.BlockChain.BatchSize)
	}
	return nil
}

// TLSConfig is tls config
type TLSConfig struct {
	Enable  bool   `yaml:"Enable"`
	CA      string `yaml:"CA"`
	RawCert string `yaml:"Cert"`
	Key     string `yaml:"Key"`
	// Pool of CA
	Pool *x509.CertPool
	Cert *tls.Certificate
}

// ConsensusType is the type of consensus
type ConsensusType int

const (
	_ ConsensusType = iota
	// SOLO is the solo
	SOLO
	// RAFT is the raft
	RAFT
	// BFT is the tendermint
	BFT
)

// ConsensusConfig is the config of consensus
type ConsensusConfig struct {
	Type ConsensusType
	BFT  TendermintConfig
	Raft RaftConfig
}

// TendermintConfig is the config of tendermint
type TendermintConfig struct {
	Path string `yaml:"Path"`
	Port struct {
		P2P int `yaml:"P2P"`
		RPC int `yaml:"RPC"`
		APP int `yaml:"APP"`
	} `yaml:"Port"`
	ID         string   `yaml:"ID"`
	P2PAddress []string `yaml:"P2PAddress"`
}

// RaftConfig is the config of raft
type RaftConfig struct {
	Path string `yaml:"Path"`
	ID   uint64 `yaml:"ID"`
	// RawNodes should be an array like [1@localhost:12346]
	RawNodes []string `yaml:"Nodes"`
	Join     bool     `yaml:"Join"`
	Nodes    map[uint64]string
	// TLS
	TLS TLSConfig `yaml:"TLS"`
}

// GetBlockChainConfig return the BlockChainConfig

// GetConsensusConfig return the ConsensusConfig
func (cfg *Config) GetConsensusConfig() (*ConsensusConfig, error) {
	var consensus ConsensusConfig
	switch cfg.Consensus.Type {
	case "solo":
		consensus.Type = SOLO
	case "raft":
		consensus.Type = RAFT
		consensus.Raft = cfg.Consensus.Raft
		if consensus.Raft.ID <= 0 {
			return nil, errors.New("Raft id should not be zero")
		}
		// then we should parse RawNodes to Nodes
		consensus.Raft.Nodes = make(map[uint64]string)
		for i := range consensus.Raft.RawNodes {
			id, url, err := parseRaftNode(consensus.Raft.RawNodes[i])
			if err != nil {
				return nil, err
			}
			consensus.Raft.Nodes[id] = url
		}
		if !util.Contain(consensus.Raft.Nodes, consensus.Raft.ID) {
			return nil, errors.New("Nodes must contain itself")
		}
		consensus.Raft.TLS = cfg.TLS
		return &consensus, nil
	case "bft":
		consensus.Type = BFT
		consensus.BFT = cfg.Consensus.Tendermint
		// check some necessary things
		if len(consensus.BFT.ID) != 40 {
			return nil, fmt.Errorf("The ID(%s) of tendermint is not legal", consensus.BFT.ID)
		}
		return &consensus, nil
	default:
		return nil, fmt.Errorf("Unsupport consensus type: %s", cfg.Consensus.Type)
	}
	return &consensus, nil
}

// DBType is the type of DB
type DBType int

const (
	_ DBType = iota
	// LEVELDB is the leveldb
	LEVELDB
)

// DBConfig is the config of db
type DBConfig struct {
	Type    DBType
	LevelDB LevelDBConfig
}

// LevelDBConfig is the config of leveldb
type LevelDBConfig struct {
	Path string
}

// GetDBConfig return the DBConfig
func (cfg *Config) GetDBConfig() (*DBConfig, error) {
	var config DBConfig
	switch cfg.DB.Type {
	case "leveldb":
		config.Type = LEVELDB
		config.LevelDB.Path = cfg.DB.LevelDB.Path
		if config.LevelDB.Path == "" {
			return nil, errors.New("The path of leveldb is not provided")
		}
	default:
		return nil, fmt.Errorf("Unsupport db type: %s", cfg.DB.Type)
	}
	return &config, nil
}

func parseRaftNode(node string) (uint64, string, error) {
	params := regexp.MustCompile(`^([\d]+)@(.+):([0-9]+)$`).FindStringSubmatch(node)
	if len(params) != 4 {
		return 0, "", errors.New("Wrong format")
	}
	id, err := strconv.ParseUint(params[1], 10, 64)
	if err != nil || id == 0 {
		return 0, "", errors.New("Wrong format")
	}
	return id, fmt.Sprintf("%s:%s", params[2], params[3]), nil
}
