package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"madledger/util"
	"os"

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
	BlockChain BlockChainConfig `yaml:"BlockChain"`
	Consensus  struct {
		Type string `yaml:"Type"`
	} `yaml:"Consensus"`
	DB struct {
		Type    string `yaml:"Type"`
		LevelDB struct {
			Dir string `yaml:"Dir"`
		} `yaml:"LevelDB"`
	} `yaml:"DB"`
}

// ServerConfig is the config of server
type ServerConfig struct {
	// Listening port for the server
	Port int `yaml:"Port"`
	// Bind address for the server
	Address string `yaml:"Address"`
	// Debug
	Debug bool `yaml:"Debug"`
}

// LoadConfig load config from the config file
func LoadConfig(cfgFile string) (*Config, error) {
	cfgBytes, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetServerConfig return the ServerConfig
func (cfg *Config) GetServerConfig() (*ServerConfig, error) {
	if cfg.Port < 1024 {
		return nil, fmt.Errorf("The port can not be %d", cfg.Port)
	}
	if cfg.Address == "" {
		return nil, errors.New("The address can not be empty")
	}
	return &ServerConfig{
		Port:    cfg.Port,
		Address: cfg.Address,
		Debug:   cfg.Debug,
	}, nil
}

// BlockChainConfig is the config of blockchain
type BlockChainConfig struct {
	BatchTimeout int    `yaml:"BatchTimeout"`
	BatchSize    int    `yaml:"BatchSize"`
	Path         string `yaml:"Path"`
	Verify       bool   `yaml:"Verify"`
}

// ConsensusType is the type of consensus
type ConsensusType int

const (
	_ ConsensusType = iota
	// SOLO is the solo
	SOLO
	// RAFT is the raft
	RAFT
	// PBFT is the pbft
	PBFT
)

// ConsensusConfig is the config of consensus
type ConsensusConfig struct {
	Type ConsensusType
}

// GetBlockChainConfig return the BlockChainConfig
func (cfg *Config) GetBlockChainConfig() (*BlockChainConfig, error) {
	var storePath = cfg.BlockChain.Path
	if storePath == "" {
		if cfg.Debug {
			storePath = getDefaultChainPath()
		} else {
			return nil, errors.New("The path of blockchain is not provided")
		}
	}
	if cfg.BlockChain.BatchTimeout <= 0 {
		return nil, fmt.Errorf("The batch timeout can not be %d", cfg.BlockChain.BatchTimeout)
	}
	if cfg.BlockChain.BatchSize <= 0 {
		return nil, fmt.Errorf("The batch size can not be %d", cfg.BlockChain.BatchSize)
	}
	return &BlockChainConfig{
		BatchTimeout: cfg.BlockChain.BatchTimeout,
		BatchSize:    cfg.BlockChain.BatchSize,
		Path:         storePath,
		Verify:       cfg.BlockChain.Verify,
	}, nil
}

// GetConsensusConfig return the ConsensusConfig
func (cfg *Config) GetConsensusConfig() (*ConsensusConfig, error) {
	var consensus ConsensusConfig
	switch cfg.Consensus.Type {
	case "solo":
		consensus.Type = SOLO
	case "raft":
		consensus.Type = RAFT
		return nil, errors.New("Raft is not supported yet")
	case "pbft":
		consensus.Type = PBFT
		return nil, errors.New("Pbft is not supported yet")
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
	Dir string
}

// GetDBConfig return the DBConfig
func (cfg *Config) GetDBConfig() (*DBConfig, error) {
	var config DBConfig
	switch cfg.DB.Type {
	case "leveldb":
		config.Type = LEVELDB
		config.LevelDB.Dir = cfg.DB.LevelDB.Dir
		if config.LevelDB.Dir == "" {
			config.LevelDB.Dir = getDefaultLevelDBPath()
		}
	default:
		return nil, fmt.Errorf("Unsupport db type: %s", cfg.DB.Type)
	}
	return &config, nil
}

func getDefaultChainPath() string {
	storePath, _ := util.MakeFileAbs("src/madledger/orderer/data/blocks", gopath)
	return storePath
}

func getDefaultLevelDBPath() string {
	storePath, _ := util.MakeFileAbs("src/madledger/orderer/data/leveldb", gopath)
	return storePath
}
