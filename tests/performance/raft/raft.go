package raft

import (
	"fmt"
	"madledger/common/util"
	"os"
	"time"

	oc "madledger/orderer/config"
	pc "madledger/peer/config"

	orderer "madledger/orderer/server"
	peer "madledger/peer/server"
)

var (
	ordererServers []*orderer.Server
	peerServer     *peer.Server
)

// Init will init environment
func Init() error {
	Clean()
	os.MkdirAll(getOrdererPath(), os.ModePerm)
	os.MkdirAll(getPeerPath(), os.ModePerm)
	os.MkdirAll(getClientsPath(), os.ModePerm)
	return newClients()
}

// Clean clean environment
func Clean() {
	os.RemoveAll(getOrdererPath())
	os.RemoveAll(getPeerPath())
	os.RemoveAll(getClientsPath())
}

// StartOrderers start orderers
func StartOrderers() error {
	for i := 1; i <= 3; i++ {
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

// StopOrderers stop orderers
func StopOrderers() {
	for i := range ordererServers {
		ordererServers[i].Stop()
		time.Sleep(500 * time.Millisecond)
	}
}

// StartPeers start peers
func StartPeers() error {
	var err error
	peerServer, err = peer.NewServer(getPeerConfig())
	if err != nil {
		return err
	}
	go func() {
		peerServer.Start()
	}()
	time.Sleep(300 * time.Millisecond)
	return nil
}

// StopPeers stop peers
func StopPeers() {
	peerServer.Stop()
	time.Sleep(300 * time.Millisecond)
}

func getOrdererPath() string {
	path, _ := util.MakeFileAbs("src/madledger/tests/performance/raft/.orderer", gopath)
	return path
}

func getPeerPath() string {
	path, _ := util.MakeFileAbs("src/madledger/tests/performance/raft/.peer", gopath)
	return path
}

func getOrdererConfig(id int) (*oc.Config, error) {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/config/orderer/raft_orderer.yaml", gopath)
	cfg, err := oc.LoadConfig(cfgFilePath)
	if err != nil {
		return nil, err
	}
	chainPath, _ := util.MakeFileAbs("data/blocks", getOrdererPath()+fmt.Sprintf("/%d", id))
	dbPath, _ := util.MakeFileAbs("data/leveldb", getOrdererPath()+fmt.Sprintf("/%d", id))
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Path = dbPath
	cfg.Consensus.Raft.Path = fmt.Sprintf("%s/%d/%s", getOrdererPath(), id, cfg.Consensus.Raft.Path)
	cfg.Consensus.Raft.ID = uint64(id)
	cfg.Port = 12345 + (id-1)*11111
	return cfg, nil
}

func getPeerConfig() *pc.Config {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/config/peer/raft_peer.yaml", gopath)
	cfg, _ := pc.LoadConfig(cfgFilePath)
	chainPath, _ := util.MakeFileAbs("data/blocks", getPeerPath())
	dbPath, _ := util.MakeFileAbs("data/leveldb", getPeerPath())
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Dir = dbPath
	key, _ := util.MakeFileAbs("src/madledger/tests/config/peer/.raft_peer.pem", gopath)
	cfg.KeyStore.Key = key
	return cfg
}
