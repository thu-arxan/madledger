package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var (
	gopath = os.Getenv("GOPATH")
)

// Config is the combination of all config
type Config struct {
	Port    int           `yaml:"Port"`
	Address string        `yaml:"Address"`
	Debug   bool          `yaml:"Debug"`
	Orderer OrdererConfig `yaml:"Orderer"`
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
