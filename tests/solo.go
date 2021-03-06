// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package tests

import (
	client "madledger/client/lib"
	"madledger/common/util"
	"time"

	oc "madledger/orderer/config"
	pc "madledger/peer/config"

	orderer "madledger/orderer/server"
	peer "madledger/peer/server"
)

var (
	soloOrdererServer *orderer.Server
	soloPeerServer    *peer.Server
)

func startSoloOrderer() error {
	cfg, err := getSoloOrdererConfig()
	if err != nil {
		return err
	}
	soloOrdererServer, err = orderer.NewServer(cfg)
	if err != nil {
		return err
	}
	go func() {
		soloOrdererServer.Start()
	}()
	time.Sleep(300 * time.Millisecond)
	return nil
}

func stopSoloOrderer() {
	soloOrdererServer.Stop()
	time.Sleep(500 * time.Millisecond)
}

func startSoloPeer() error {
	var err error
	soloPeerServer, err = peer.NewServer(getSoloPeerConfig())
	if err != nil {
		return err
	}
	go func() {
		soloPeerServer.Start()
	}()
	time.Sleep(300 * time.Millisecond)
	return nil
}

func stopSoloPeer() {
	soloPeerServer.Stop()
	time.Sleep(300 * time.Millisecond)
}

func getSoloOrdererConfig() (*oc.Config, error) {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/config/orderer/solo_orderer.yaml", gopath)
	cfg, err := oc.LoadConfig(cfgFilePath)
	if err != nil {
		return nil, err
	}
	chainPath, _ := util.MakeFileAbs("src/madledger/tests/.orderer/data/blocks", gopath)
	dbPath, _ := util.MakeFileAbs("src/madledger/tests/.orderer/data/leveldb", gopath)
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Path = dbPath
	return cfg, nil
}

func getSoloPeerConfig() *pc.Config {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/config/peer/solo_peer.yaml", gopath)
	cfg, _ := pc.LoadConfig(cfgFilePath)
	chainPath, _ := util.MakeFileAbs("src/madledger/tests/.peer/data/blocks", gopath)
	dbPath, _ := util.MakeFileAbs("src/madledger/tests/.peer/data/leveldb", gopath)
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Dir = dbPath
	return cfg
}

func getSoloClient() (*client.Client, error) {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/config/client/solo_client.yaml", gopath)
	c, err := client.NewClient(cfgFilePath)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func getSoloHTTPClient() (*client.HTTPClient, error) {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/config/client/solo_client.yaml", gopath)
	c, err := client.NewHTTPClient(cfgFilePath)
	if err != nil {
		return nil, err
	}
	return c, nil
}
