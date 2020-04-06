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
	ordererServer *orderer.Server
	peerServer    *peer.Server
)

// Init will init environment
func Init(clientSize int) error {
	Clean()
	os.MkdirAll(getOrdererPath(), os.ModePerm)
	os.MkdirAll(getPeerPath(), os.ModePerm)
	os.MkdirAll(getClientsPath(), os.ModePerm)
	clients = make([]*client.Client, clientSize)
	return newClients()
}

// Clean clean environment
func Clean() {
	os.RemoveAll(getOrdererPath())
	os.RemoveAll(getPeerPath())
	os.RemoveAll(getClientsPath())
}

// StartOrderers start orderers(solo will only have one)
func StartOrderers() error {
	cfg, err := getOrdererConfig()
	if err != nil {
		return err
	}
	ordererServer, err = orderer.NewServer(cfg)
	if err != nil {
		return err
	}
	go func() {
		ordererServer.Start()
	}()
	time.Sleep(300 * time.Millisecond)
	return nil
}

// StopOrderers stop orderers(solo will only have one)
func StopOrderers() {
	ordererServer.Stop()
	time.Sleep(500 * time.Millisecond)
}

// StartPeers start peers(solo will only have one)
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

// StopPeers stop peers(solo will only have one)
func StopPeers() {
	peerServer.Stop()
	time.Sleep(300 * time.Millisecond)
}

func getOrdererPath() string {
	path, _ := util.MakeFileAbs("src/madledger/tests/performance/solo/.orderer", gopath)
	return path
}

func getPeerPath() string {
	path, _ := util.MakeFileAbs("src/madledger/tests/performance/solo/.peer", gopath)
	return path
}

func getOrdererConfig() (*oc.Config, error) {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/config/orderer/solo_orderer.yaml", gopath)
	cfg, err := oc.LoadConfig(cfgFilePath)
	if err != nil {
		return nil, err
	}
	chainPath, _ := util.MakeFileAbs("data/blocks", getOrdererPath())
	dbPath, _ := util.MakeFileAbs("data/leveldb", getOrdererPath())
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Path = dbPath
	return cfg, nil
}

func getPeerConfig() *pc.Config {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/config/peer/solo_peer.yaml", gopath)
	cfg, _ := pc.LoadConfig(cfgFilePath)
	chainPath, _ := util.MakeFileAbs("data/blocks", getPeerPath())
	dbPath, _ := util.MakeFileAbs("data/leveldb", getPeerPath())
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Dir = dbPath
	key, _ := util.MakeFileAbs("src/madledger/tests/config/peer/.solo_peer.pem", gopath)
	cfg.KeyStore.Key = key
	return cfg
}
