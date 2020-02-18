package raft

import (
	"fmt"
	client "madledger/client/lib"
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
	peerServers    []*peer.Server
)

// Init will init environment
func Init(clientSize, peerNum int) error {
	Clean()
	os.MkdirAll(getOrdererPath(), os.ModePerm)
	os.MkdirAll(getPeerPath(), os.ModePerm)
	os.MkdirAll(getClientsPath(), os.ModePerm)
	clients = make([]*client.Client, clientSize)
	return newClients(peerNum)
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
			if err := ordererServer.Start(); err != nil {
				panic("order start failed: " + err.Error())
			}
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
func StartPeers(num int) error {
	for i := 1; i <= num; i++ {
		peerServer, err := peer.NewServer(getPeerConfig(i))
		if err != nil {
			return err
		}
		peerServers = append(peerServers, peerServer)
		go func() {
			if err := peerServer.Start(); err != nil {
				panic("peer start failed: " + err.Error())
			}
		}()
		time.Sleep(300 * time.Millisecond)
	}
	return nil
}

// StopPeers stop peers
func StopPeers() {
	for i := range peerServers {
		peerServers[i].Stop()
		time.Sleep(300 * time.Millisecond)
	}
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

func getPeerConfig(id int) *pc.Config {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/config/peer/raft_peer.yaml", gopath)
	cfg, _ := pc.LoadConfig(cfgFilePath)
	chainPath, _ := util.MakeFileAbs("data/blocks", fmt.Sprintf("%s/%d", getPeerPath(), id))
	dbPath, _ := util.MakeFileAbs("data/leveldb", fmt.Sprintf("%s/%d", getPeerPath(), id))
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Dir = dbPath
	key, _ := util.MakeFileAbs(fmt.Sprintf("src/madledger/tests/config/peer/.raft_peer%d.pem", id), gopath)
	cfg.KeyStore.Key = key
	cfg.Port = 23333 + (id - 1)
	return cfg
}
