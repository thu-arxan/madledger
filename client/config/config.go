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
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"madledger/common/crypto"
	"madledger/core"

	yaml "gopkg.in/yaml.v2"
)

// Config is the combination of all config
type Config struct {
	Debug    bool          `yaml:"Debug"`
	TLS      TLSConfig     `yaml:"TLS"`
	Orderer  OrdererConfig `yaml:"Orderer"`
	Peer     PeerConfig    `yaml:"Peer"`
	KeyStore struct {
		Keys  []string            `yaml:"Keys"`
		Privs []crypto.PrivateKey `yaml:"-"`
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
	err = cfg.loadTLSConfig()
	if err != nil {
		return nil, err
	}
	if err := cfg.loadOrdererConfig(); err != nil {
		return nil, err
	}
	if err := cfg.loadPeerConfig(); err != nil {
		return nil, err
	}
	if err := cfg.loadKeys(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// loadTLSConfig load tls config
func (cfg *Config) loadTLSConfig() error {
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
	Address     []string `yaml:"Address"`
	HTTPAddress []string `yaml:"HTTPAddress"`
}

// loadOrdererConfig load the orderer config
func (cfg *Config) loadOrdererConfig() error {
	if len(cfg.Orderer.Address) == 0 {
		return errors.New("The address of orderer should not be nil")
	}
	return nil
}

// PeerConfig is the config of peer
type PeerConfig struct {
	Address     []string `yaml:"Address"`
	HTTPAddress []string `yaml:"HTTPAddress"`
}

// loadPeerConfig load the peer config
func (cfg *Config) loadPeerConfig() error {
	// if len(cfg.Peer.Address) == 0 {
	// 	return errors.New("The address of peer should not be nil")
	// }
	return nil
}

// KeyStoreConfig is the config of KeyStore
type KeyStoreConfig struct {
	Keys []crypto.PrivateKey
}

// loadKeys return the key store config
func (cfg *Config) loadKeys() error {
	if len(cfg.KeyStore.Keys) == 0 {
		return errors.New("The keys should not be nil")
	}
	var keys []crypto.PrivateKey
	for _, keyFile := range cfg.KeyStore.Keys {
		key, err := crypto.LoadPrivateKeyFromFile(keyFile)
		if err != nil {
			return err
		}
		keys = append(keys, key)
	}
	cfg.KeyStore.Privs = keys
	return nil
}

// LoadPeerAddress load config from the config file
func LoadPeerAddress(cfgFile string) (*PeerConfig, error) {
	cfgBytes, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	var cfg PeerConfig
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return nil, err
	}
	if len(cfg.Address) == 0 {
		return nil, errors.New("The address of peer is nil")
	}
	return &cfg, nil
}

// SavePeerCache .
func SavePeerCache(name string, peers []string) error {
	b, err := yaml.Marshal(&PeerConfig{
		Address:     peers,
		HTTPAddress: nil,
	})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(name+".yaml", b, 0777)
}

// MemberConfig .
type MemberConfig struct {
	Admins  []string `yaml:"Admins"`
	Members []string `yaml:"Members"`
}

// GetMembers .
func GetMembers(file string) ([]*core.Member, []*core.Member, error) {
	Bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}
	var cfg MemberConfig
	err = yaml.Unmarshal(Bytes, &cfg)
	if err != nil {
		return nil, nil, err
	}
	var admins, members []*core.Member
	for _, adminstr := range cfg.Admins {
		b, err := hex.DecodeString(adminstr)
		if err != nil {
			return nil, nil, err
		}
		pk, err := crypto.NewPublicKey(b, crypto.KeyAlgoSM2)
		if err != nil {
			return nil, nil, err
		}
		admin, err := core.NewMember(pk, "admin")
		if err != nil {
			return nil, nil, err
		}
		admins = append(admins, admin)
	}

	for _, memberstr := range cfg.Members {
		b, err := hex.DecodeString(memberstr)
		if err != nil {
			return nil, nil, err
		}
		pk, err := crypto.NewPublicKey(b, crypto.KeyAlgoSM2)
		if err != nil {
			return nil, nil, err
		}
		member, err := core.NewMember(pk, "member")
		if err != nil {
			return nil, nil, err
		}
		members = append(members, member)
	}

	return admins, members, nil
}
