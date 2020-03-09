package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/core"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var (
	gopath = os.Getenv("GOPATH")
)

// Config is the combination of all config
type Config struct {
	Port    int    `yaml:"Port"`
	Address string `yaml:"Address"`
	Debug   bool   `yaml:"Debug"`
	// TLS
	TLS        TLSConfig        `yaml:"TLS"`
	BlockChain BlockChainConfig `yaml:"BlockChain"`
	Orderer    OrdererConfig    `yaml:"Orderer"`
	DB         struct {
		T       string `yaml:"Type"`
		Type    DBType `yaml:"-"`
		LevelDB struct {
			Dir string `yaml:"Dir"`
		} `yaml:"LevelDB"`
	} `yaml:"DB"`
	KeyStore struct {
		Key string `yaml:"Key"`
	} `yaml:"KeyStore"`
}

// TLSConfig ...
type TLSConfig struct {
	Enable  bool   `yaml:"Enable"`
	CA      string `yaml:"CA"`
	RawCert string `yaml:"Cert"`
	Key     string `yaml:"Key"`
	// Pool of CA
	Pool *x509.CertPool
	Cert *tls.Certificate
}

// ServerConfig is the config of server
type ServerConfig struct {
	// Listening port for the server
	Port int `yaml:"Port"`
	// Bind address for the server
	Address string `yaml:"Address"`
	// Debug
	Debug bool `yaml:"Debug"`
	// TLS
	TLS TLSConfig `yaml:"TLS"`
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
	if cfg.Port < 1024 {
		return nil, fmt.Errorf("The port can not be %d", cfg.Port)
	}
	if cfg.Address == "" {
		return nil, errors.New("The address can not be empty")
	}
	if err = cfg.loadTLSConfig(); err != nil {
		return nil, err
	}
	if err = cfg.loadOrdererConfig(); err != nil {
		return nil, err
	}
	if err = cfg.loadBlockChainConfig(); err != nil {
		return nil, err
	}
	if err = cfg.loadDBConfig(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// loadTLSConfig check the tls config and set necessary things
func (cfg *Config) loadTLSConfig() error {
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

// OrdererConfig is the config of orderer
type OrdererConfig struct {
	Address []string `yaml:"Address"`
}

// loadOrdererConfig check the orderer config and set necessary things
func (cfg *Config) loadOrdererConfig() error {
	if len(cfg.Orderer.Address) == 0 {
		return errors.New("orderer address is not setted")
	}
	return nil
}

// BlockChainConfig is the config of blockchain
type BlockChainConfig struct {
	Path string `yaml:"Path"`
}

// loadBlockChainConfig check the blockchain config and set necessary things
func (cfg *Config) loadBlockChainConfig() error {
	if cfg.BlockChain.Path == "" {
		if cfg.Debug {
			cfg.BlockChain.Path = getDefaultChainPath()
		} else {
			return errors.New("The path of blockchain is not provided")
		}
	}
	return nil
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
func (cfg *Config) loadDBConfig() error {
	switch cfg.DB.T {
	case "leveldb":
		cfg.DB.Type = LEVELDB
		if cfg.DB.LevelDB.Dir == "" {
			if !cfg.Debug {
				return errors.New("leveldb path is not setted")
			}
			cfg.DB.LevelDB.Dir = getDefaultLevelDBPath()
		}
	default:
		return fmt.Errorf("unsupport db type: %s", cfg.DB.T)
	}
	return nil
}

// GetIdentity return the identity of peer
func (cfg *Config) GetIdentity() (*core.Member, error) {
	if cfg.KeyStore.Key == "" {
		return nil, errors.New("The key should not be nil")
	}
	privKey, err := crypto.LoadPrivateKeyFromFile(cfg.KeyStore.Key)
	if err != nil {
		return nil, err
	}

	identity, err := core.NewMember(privKey.PubKey(), "self")
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
