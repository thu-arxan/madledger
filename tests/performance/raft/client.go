package raft

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
    - localhost:23456
    - localhost:34567
  
# Address of peers
Peer:
  Address:
    - localhost:23333

# KeyStore manage some private keys
KeyStore:
  Keys:
    - <<<KEYFILE>>>
`
)

var (
	clients    = make([]*client.Client, 400)
	clientInit = false
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
	path, _ := util.MakeFileAbs("src/madledger/tests/performance/raft/.clients", gopath)
	return path
}
