package tests

import (
	"madledger/common/util"
	oc "madledger/orderer/config"
	orderer "madledger/orderer/server"
	pc "madledger/peer/config"
	peer "madledger/peer/server"
	"os"
	"testing"
	"time"

	client "madledger/client/lib"
)

/*
* CircumstanceSolo begins from a empty environment and defines some operations as below.
* 1. Create channel test.
* 2. Create a contract.
* 3. Call the contract in different ways.
* 4. During main operates, there are some necessary query to make sure everything is ok.
 */

func TestInitCircumstanceSolo(t *testing.T) {
	err := initDir(".orderer")
	if err != nil {
		t.Fatal(err)
	}
	err = initDir(".peer")
	if err != nil {
		t.Fatal(err)
	}
	err = initDir(".client")
	if err != nil {
		t.Fatal(err)
	}

}

func TestCreateChannel(t *testing.T) {
	startSoloOrderer()
	startSoloPeer()
	client, err := getSoloClient()
	if err != nil {
		t.Fatal(err)
	}
	// then add a channel
	err = client.CreateChannel("test")
	if err != nil {
		t.Fatal(err)
	}
	// then query channels, however this should modify the function of client
	// todo

	// create channel test again
	err = client.CreateChannel("test")
	if err == nil {
		t.Fatal(err)
	}
}

func TestEnd(t *testing.T) {
	os.RemoveAll(".orderer")
	os.RemoveAll(".peer")
	os.RemoveAll(".client")
}

func startSoloOrderer() error {
	cfg, err := getSoloOrdererConfig()
	if err != nil {
		return err
	}
	server, err := orderer.NewServer(cfg)
	if err != nil {
		return err
	}
	go func() {
		server.Start()
	}()
	time.Sleep(300 * time.Millisecond)
	return nil
}

func startSoloPeer() error {
	server, err := peer.NewServer(getSoloPeerConfig())
	if err != nil {
		return err
	}
	go func() {
		server.Start()
	}()
	time.Sleep(300 * time.Millisecond)
	return nil
}

func getSoloOrdererConfig() (*oc.Config, error) {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/solo_orderer.yaml", gopath)
	cfg, err := oc.LoadConfig(cfgFilePath)
	if err != nil {
		return nil, err
	}
	chainPath, _ := util.MakeFileAbs("src/madledger/tests/.orderer/data/blocks", gopath)
	dbPath, _ := util.MakeFileAbs("src/madledger/tests/.orderer/data/leveldb", gopath)
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Dir = dbPath
	return cfg, nil
}

func getSoloPeerConfig() *pc.Config {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/solo_peer.yaml", gopath)
	cfg, _ := pc.LoadConfig(cfgFilePath)
	chainPath, _ := util.MakeFileAbs("src/madledger/tests/.peer/data/blocks", gopath)
	dbPath, _ := util.MakeFileAbs("src/madledger/tests/.peer/data/leveldb", gopath)
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Dir = dbPath
	return cfg
}

func getSoloClient() (*client.Client, error) {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/solo_client.yaml", gopath)
	c, err := client.NewClient(cfgFilePath)
	if err != nil {
		return nil, err
	}
	return c, nil
}
