// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package bft

import (
	"fmt"
	client "madledger/client/lib"
	"madledger/common/util"
)

const (
	clientConfigTemplate = `#############################################################################
#   This is a configuration file for the MadLedger client.
#############################################################################

# Should be false or true (default: true)
Debug: true

# Address of orderers
Orderer:
  Address:
    - localhost:12345
    - localhost:23456
    - localhost:34567
  
# Address of peers
Peer:
  Address:
    <<<ADDRESS1>>>
    <<<ADDRESS2>>>
    <<<ADDRESS3>>>

# KeyStore manage some private keys
KeyStore:
  Keys:
    - <<<KEYFILE>>>
`
)

var (
	clients    []*client.Client
	clientInit = false
)

// GetClients will return clients
func GetClients() []*client.Client {
	if !clientInit {
		for i := range clients {
			cfgPath, _ := util.MakeFileAbs(fmt.Sprintf("%d/explorer-client.yaml", i), getClientsPath())
			client, err := client.NewClient(cfgPath)
			if err != nil {
				panic(err)
			}
			clients[i] = client
		}
		clientInit = true
	}
	return clients
}
