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
	"fmt"
	"madledger/common/util"
	pc "madledger/peer/config"
	peer "madledger/peer/server"
	"time"
)

var (
	peers = make(map[int]*peer.Server)
)

func startPeers(size int) error {
	for i := 0; i < size; i++ {
		cfg := getPeerConfig(i)
		s, err := peer.NewServer(cfg)
		if err != nil {
			return err
		}
		peers[i] = s
		go func() {
			s.Start()
		}()
	}

	time.Sleep(1000 * time.Millisecond)
	return nil
}

func stopPeers(size int) {
	for i := 0; i < size; i++ {
		if util.Contain(peers, i) {
			peers[i].Stop()
			delete(peers, i)
		}
	}
	time.Sleep(1000 * time.Millisecond)
}

func getPeerConfig(i int) *pc.Config {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/config/peer/peer.yaml", gopath)
	cfg, _ := pc.LoadConfig(cfgFilePath)
	chainPath, _ := util.MakeFileAbs(fmt.Sprintf("src/madledger/tests/.peer%d/data/blocks", i), gopath)
	dbPath, _ := util.MakeFileAbs(fmt.Sprintf("src/madledger/tests/.peer%d/data/leveldb", i), gopath)
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Dir = dbPath
	// then set port
	cfg.Port += i
	// then set key
	cfg.KeyStore.Key = fmt.Sprintf("config/peer/.solo_peer%d.pem", i)
	return cfg
}
