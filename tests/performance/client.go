package performance

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
	soloClientConfigTemplate = `#############################################################################
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
	clients [200]*client.Client
)

func init() {
	gopath := os.Getenv("GOPATH")
	path, _ := util.MakeFileAbs("src/madledger/tests/performance/.clients", gopath)
	os.RemoveAll(path)
	os.MkdirAll(path, os.ModePerm)
	for i := range clients {
		// client path
		cp, _ := util.MakeFileAbs(fmt.Sprintf("%d", i), path)
		os.MkdirAll(cp, os.ModePerm)
		if err := newSoloClient(cp); err != nil {
			panic(err)
		}
	}
}

func newSoloClient(path string) error {
	cfgPath, _ := util.MakeFileAbs("client.yaml", path)
	keyStorePath, _ := util.MakeFileAbs(".keystore", path)
	os.MkdirAll(keyStorePath, os.ModePerm)
	keyPath, err := cutil.GeneratePrivateKey(keyStorePath)
	if err != nil {
		return err
	}
	var cfg = soloClientConfigTemplate
	cfg = strings.Replace(cfg, "<<<KEYFILE>>>", keyPath, 1)

	return ioutil.WriteFile(cfgPath, []byte(cfg), os.ModePerm)
}
