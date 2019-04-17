package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/core/types"
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
	Orderer    OrdererConfig    `yaml:"Orderer"`
	DB         struct {
		Type    string `yaml:"Type"`
		LevelDB struct {
			Dir string `yaml:"Dir"`
		} `yaml:"LevelDB"`
	} `yaml:"DB"`
	KeyStore struct {
		Key string `yaml:"Key"`
	} `yaml:"KeyStore"`
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

// OrdererConfig is the config of orderer
type OrdererConfig struct {
	Address []string `yaml:"Address"`
}

// GetOrdererConfig return the orderer config
func (cfg *Config) GetOrdererConfig() (*OrdererConfig, error) {
	return &OrdererConfig{
		Address: cfg.Orderer.Address,
	}, nil
}

// BlockChainConfig is the config of blockchain
type BlockChainConfig struct {
	Path string `yaml:"Path"`
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
	return &BlockChainConfig{
		Path: storePath,
	}, nil
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

// GetIdentity return the identity of peer
func (cfg *Config) GetIdentity() (*types.Member, error) {
	if cfg.KeyStore.Key == "" {
		return nil, errors.New("The key should not be nil")
	}
	privKey, err := crypto.LoadPrivateKeyFromFile(cfg.KeyStore.Key)
	if err != nil {
		return nil, err
	}

	identity, err := types.NewMember(privKey.PubKey(), "self")
	if err != nil {
		return nil, err
	}

	return identity, nil
}

func getDefaultLevelDBPath() string {
	storePath, _ := util.MakeFileAbs("src/madledger/peer/data/leveldb", gopath)
	return storePath
}

func getDefaultChainPath() string {
	storePath, _ := util.MakeFileAbs("src/madledger/peer/data/blocks", gopath)
	return storePath
}
