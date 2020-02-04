package consensus

import (
	"crypto/tls"
	"crypto/x509"
)

// Config is the config of consensus
type Config struct {
	Timeout int
	MaxSize int
	Resume  bool
	Number  uint64
	TLS TLSConfig
}

type TLSConfig struct {
	Enable  bool
	CA      string
	RawCert string
	Key     string
	// Pool of CA
	Pool *x509.CertPool
	Cert *tls.Certificate
}

// DefaultConfig is the DefaultConfig
func DefaultConfig() Config {
	return Config{
		Timeout: 1000,
		MaxSize: 10,
		Number:  1,
		Resume:  false,
	}
}
