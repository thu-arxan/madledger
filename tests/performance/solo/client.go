// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package solo

import (
	"fmt"
	"io/ioutil"
	client "madledger/client/lib"
	cutil "madledger/client/util"
	"madledger/common/util"
	"os"
	"strings"
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

var (
	clientInit = false
	clients    []*client.Client
)

// GetClients will return clients
func GetClients() []*client.Client {
	if !clientInit {
		for i := range clients {
			cfgPath, _ := util.MakeFileAbs(fmt.Sprintf("%d/client.yaml", i), getClientsPath())
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

func newClients() error {
	for i := range clients {
		cp, _ := util.MakeFileAbs(fmt.Sprintf("%d", i), getClientsPath())
		os.MkdirAll(cp, os.ModePerm)
		if err := newClient(cp); err != nil {
			return err
		}
	}

	return nil
}

func newClient(path string) error {
	cfgPath, _ := util.MakeFileAbs("client.yaml", path)
	keyStorePath, _ := util.MakeFileAbs(".keystore", path)
	os.MkdirAll(keyStorePath, os.ModePerm)

	keyPath, err := cutil.GeneratePrivateKey(keyStorePath)
	if err != nil {
		return err
	}

	var cfg = clientConfigTemplate
	cfg = strings.Replace(cfg, "<<<KEYFILE>>>", keyPath, 1)

	return ioutil.WriteFile(cfgPath, []byte(cfg), os.ModePerm)
}

func getClientsPath() string {
	path, _ := util.MakeFileAbs("src/madledger/tests/performance/solo/.clients", gopath)
	return path
}
