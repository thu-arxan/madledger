// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package cmd

const (
	cfgTemplate = `#############################################################################
#   This is a configuration file for the MadLedger client.
#############################################################################

# Should be false or true (default: true)
Debug: true

# Configure for the TLS
TLS:
  # Should be true of false (default: false)
  Enable: false
  # The path of CA cert, it should not be empty if Enable is true
  CA: 
  # Cert of the client, it should not be empty if Enable is true
  Cert: 
  # Key of the client, it should not be empty if Enable is true
  Key: 

# Address of orderers
Orderer:
  Address:
    - localhost:12345
  
# Address of peers
Peer:
  Address:
    - localhost:23456

# KeyStore manage some private keys
KeyStore:
  Keys:
    - <<<KEYFILE>>>
`
)
