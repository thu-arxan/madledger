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

func newClients(peerNum int) error {
	for i := range clients {
		cp, _ := util.MakeFileAbs(fmt.Sprintf("%d", i), getClientsPath())
		os.MkdirAll(cp, os.ModePerm)
		if err := newClient(cp, peerNum); err != nil {
			return err
		}
	}

	return nil
}

func newClient(path string, peerNum int) error {
	cfgPath, _ := util.MakeFileAbs("client.yaml", path)
	keyStorePath, _ := util.MakeFileAbs(".keystore", path)
	os.MkdirAll(keyStorePath, os.ModePerm)

	keyPath, err := cutil.GeneratePrivateKey(keyStorePath)
	if err != nil {
		return err
	}

	var cfg = clientConfigTemplate
	for i := 1; i <= peerNum; i++ {
		port := 23333 + (i - 1)
		cfg = strings.Replace(cfg, fmt.Sprintf("<<<ADDRESS%d>>>", i), fmt.Sprintf("- localhost:%d", port), 1)
	}
	for i := peerNum + 1; i <= 3; i++ {
		cfg = strings.Replace(cfg, fmt.Sprintf("<<<ADDRESS%d>>>", i), "", 1)
	}
	cfg = strings.Replace(cfg, "<<<KEYFILE>>>", keyPath, 1)

	return ioutil.WriteFile(cfgPath, []byte(cfg), os.ModePerm)
}

func getClientsPath() string {
	path, _ := util.MakeFileAbs("src/madledger/tests/performance/raft/.clients", gopath)
	return path
}
