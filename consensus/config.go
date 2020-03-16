// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

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
	TLS     TLSConfig
}

// TLSConfig ...
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
