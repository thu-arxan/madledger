package bft

import (
	"fmt"
	"io/ioutil"
	client "madledger/client/lib"
	cutil "madledger/client/util"
	"madledger/common/util"
	oc "madledger/orderer/config"
	orderer "madledger/orderer/server"
	pc "madledger/peer/config"
	peer "madledger/peer/server"
	"os"
	"strings"
	"time"

	"github.com/otiai10/copy"
)

var (
	ordererServers []*orderer.Server
	peerServers    []*peer.Server
)

// Init will init environment
func Init(clientSize, peerNum int) error {
	Clean()
	// copy orderder config
	if err := copy.Copy(gopath+"/src/madledger/tests/performance/.orderer", getOrdererPath()); err != nil {
		return err
	}
	os.MkdirAll(getPeerPath(), os.ModePerm)
	os.MkdirAll(getClientsPath(), os.ModePerm)
	clients = make([]*client.Client, clientSize)
	return newClients(peerNum)
}

func getOrdererPath() string {
	path, _ := util.MakeFileAbs("src/madledger/tests/performance/bft/.orderer", gopath)
	return path
}

func getPeerPath() string {
	path, _ := util.MakeFileAbs("src/madledger/tests/performance/bft/.peer", gopath)
	return path
}

func getClientsPath() string {
	path, _ := util.MakeFileAbs("src/madledger/tests/performance/bft/.clients", gopath)
	return path
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

// StartOrderers start orderers
func StartOrderers() error {
	for i := 1; i <= 4; i++ {
		cfg, err := getOrdererConfig(i)
		if err != nil {
			return err
		}
		ordererServer, err := orderer.NewServer(cfg)
		if err != nil {
			return err
		}
		ordererServers = append(ordererServers, ordererServer)
		go func() {
			ordererServer.Start()
		}()
		time.Sleep(300 * time.Millisecond)
	}
	return nil
}

func getOrdererConfig(id int) (*oc.Config, error) {
	cfgFilePath := fmt.Sprintf("%s/%d/orderer.yaml", getOrdererPath(), id)
	cfg, err := oc.LoadConfig(cfgFilePath)
	if err != nil {
		return nil, err
	}
	chainPath, _ := util.MakeFileAbs("data/blocks", getOrdererPath()+fmt.Sprintf("/%d", id))
	dbPath, _ := util.MakeFileAbs("data/leveldb", getOrdererPath()+fmt.Sprintf("/%d", id))
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Path = dbPath
	cfg.Consensus.Tendermint.Path = fmt.Sprintf("%s/%d/%s", getOrdererPath(), id, cfg.Consensus.Tendermint.Path)
	cfg.Port = 12345 + (id-1)*11111

	return cfg, nil
}

// StartPeers start peers
func StartPeers(num int) error {
	for i := 1; i <= num; i++ {
		peerServer, err := peer.NewServer(getPeerConfig(i))
		if err != nil {
			return err
		}
		peerServers = append(peerServers, peerServer)
		go func() {
			peerServer.Start()
		}()
		time.Sleep(300 * time.Millisecond)
	}
	return nil
}

func getPeerConfig(id int) *pc.Config {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/config/peer/bft_peer.yaml", gopath)
	cfg, _ := pc.LoadConfig(cfgFilePath)
	chainPath, _ := util.MakeFileAbs("data/blocks", fmt.Sprintf("%s/%d", getPeerPath(), id))
	dbPath, _ := util.MakeFileAbs("data/leveldb", fmt.Sprintf("%s/%d", getPeerPath(), id))
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Dir = dbPath
	key, _ := util.MakeFileAbs(fmt.Sprintf("src/madledger/tests/config/peer/.bft_peer%d.pem", id), gopath)
	cfg.KeyStore.Key = key
	cfg.Port = 23333 + (id - 1)
	return cfg
}

// Clean clean environment
func Clean() {
	os.RemoveAll(getOrdererPath())
	os.RemoveAll(getPeerPath())
	os.RemoveAll(getClientsPath())
}
