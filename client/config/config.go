package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"madledger/common/crypto"

	yaml "gopkg.in/yaml.v2"
)

// Config is the combination of all config
type Config struct {
	Debug    bool          `yaml:"Debug"`
	TLS      TLSConfig     `yaml:"TLS"`
	Orderer  OrdererConfig `yaml:"Orderer"`
	Peer     PeerConfig    `yaml:"Peer"`
	KeyStore struct {
		Keys []string `yaml:"Keys"`
	} `yaml:"KeyStore"`
}

// TLSConfig is the config of tls
type TLSConfig struct {
	Enable  bool   `yaml:"Enable"`
	CA      string `yaml:"CA"`
	RawCert string `yaml:"Cert"`
	Key     string `yaml:"Key"`
	// Pool of CA
	Pool *x509.CertPool
	Cert *tls.Certificate
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
	err = cfg.GetTLSConfig()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetTLSConfig set tls config
// TODO: maybe a better name
func (cfg *Config) GetTLSConfig() error {
	if cfg.TLS.Enable {
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

// GetOrdererConfig return the orderer config
func (cfg *Config) GetOrdererConfig() (*OrdererConfig, error) {
	if len(cfg.Orderer.Address) == 0 {
		return nil, errors.New("The address of orderer should not be nil")
	}
	return &OrdererConfig{
		Address: cfg.Orderer.Address,
	}, nil
}

// PeerConfig is the config of peer
type PeerConfig struct {
	Address []string `yaml:"Address"`
}

// GetPeerConfig return the peer config
func (cfg *Config) GetPeerConfig() (*PeerConfig, error) {
	if len(cfg.Peer.Address) == 0 {
		return nil, errors.New("The address of peer should not be nil")
	}
	return &PeerConfig{
		Address: cfg.Peer.Address,
	}, nil
}

// KeyStoreConfig is the config of KeyStore
type KeyStoreConfig struct {
	Keys []crypto.PrivateKey
}

// GetKeyStoreConfig return the keystore config
func (cfg *Config) GetKeyStoreConfig() (*KeyStoreConfig, error) {
	if len(cfg.KeyStore.Keys) == 0 {
		return nil, errors.New("The keys should not be nil")
	}
	var keys []crypto.PrivateKey
	for _, keyFile := range cfg.KeyStore.Keys {
		key, err := crypto.LoadPrivateKeyFromFile(keyFile)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return &KeyStoreConfig{
		Keys: keys,
	}, nil
}
