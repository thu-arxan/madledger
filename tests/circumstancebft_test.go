package tests

import (
	"fmt"
	cc "madledger/client/config"
	client "madledger/client/lib"
	"madledger/common/util"
	oc "madledger/orderer/config"
	orderer "madledger/orderer/server"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/require"
)

/*
* This test will start from a empty environment and start some orderers support bft consensus.
* This test will include operations below.
 */

var (
	bftOrderers [4]*orderer.Server
	bftClients  [4]*client.Client
)

func TestBFT(t *testing.T) {
	require.NoError(t, initBFTEnvironment())
}

func TestBFTRun(t *testing.T) {
	for i := range bftOrderers {
		server, err := newBFTOrderer(i)
		require.NoError(t, err)
		bftOrderers[i] = server
	}
	// then we can run orderers
	for i := range bftOrderers {
		go func(t *testing.T, i int) {
			err := bftOrderers[i].Start()
			require.NoError(t, err)
		}(t, i)
	}
	time.Sleep(2 * time.Second)
}

func TestBFTLoadClients(t *testing.T) {
	for i := range bftClients {
		clientPath := getBFTClientPath(i)
		cfgPath := getBFTClientConfigPath(i)
		cfg, err := cc.LoadConfig(cfgPath)
		require.NoError(t, err)
		re, _ := regexp.Compile("^.*[.]keystore")
		for i := range cfg.KeyStore.Keys {
			cfg.KeyStore.Keys[i] = clientPath + "/.keystore" + re.ReplaceAllString(cfg.KeyStore.Keys[i], "")
		}
		client, err := client.NewClientFromConfig(cfg)
		require.NoError(t, err)
		bftClients[i] = client
	}
}

func TestBFTCreateChannels(t *testing.T) {
	var wg sync.WaitGroup
	var channels []string
	for i := range bftClients {
		// each client will create 5 channels
		for m := 0; m < 5; m++ {
			wg.Add(1)
			go func(t *testing.T, i int) {
				defer wg.Done()
				client := bftClients[i]
				channel := strings.ToLower(util.RandomString(16))
				channels = append(channels, channel)
				err := client.CreateChannel(channel, true, nil, nil)
				require.NoError(t, err)
			}(t, i)
		}
	}
	wg.Wait()
	// then we will check if all channels are create successful
	time.Sleep(2 * time.Second)
	for i := range bftClients {
		wg.Add(1)
		go func(t *testing.T, i int) {
			defer wg.Done()
			client := bftClients[i]
			infos, err := client.ListChannel(false)
			require.NoError(t, err)
			require.Len(t, infos, len(channels))
			for i := range infos {
				require.True(t, util.Contain(channels, infos[i].Name))
			}
		}(t, i)
	}
	wg.Wait()
}

func TestBFTOrdererRestart(t *testing.T) {
	bftOrderers[1].Stop()
	os.RemoveAll(getBFTOrdererDataPath(1))
	server, err := newBFTOrderer(1)
	require.NoError(t, err)
	bftOrderers[1] = server
	go func(t *testing.T) {
		require.NoError(t, bftOrderers[1].Start())
	}(t)
	time.Sleep(2000 * time.Millisecond)
	for i := range bftOrderers {
		require.True(t, util.IsDirSame(getBFTOrdererBlockPath(0), getBFTOrdererBlockPath(i)), fmt.Sprintf("Orderer %d is not same with 0", i))
	}
}

func TestBFTReCreateChannels(t *testing.T) {
	// Here we recreate 5 channels
	for i := 0; i < 5; i++ {
		channel := strings.ToLower(util.RandomString(16))
		err := bftClients[0].CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
	}
	time.Sleep(2 * time.Second)
	// Then we will list channels
	for i := range bftClients {
		infos, err := bftClients[i].ListChannel(false)
		require.NoError(t, err)
		// todo: compare infos here
		fmt.Println(infos)
	}
}

func TestBFTEnd(t *testing.T) {
	for i := range bftOrderers {
		bftOrderers[i].Stop()
	}
	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/.bft"))
}

// initBFTEnvironment will remove old test folders and copy necessary folders
func initBFTEnvironment() error {
	gopath := os.Getenv("GOPATH")
	if err := os.RemoveAll(gopath + "/src/madledger/tests/.bft"); err != nil {
		return err
	}
	if err := copy.Copy(gopath+"/src/madledger/env/bft/.orderers", gopath+"/src/madledger/tests/.bft/orderers"); err != nil {
		return err
	}
	if err := copy.Copy(gopath+"/src/madledger/env/bft/.clients", gopath+"/src/madledger/tests/.bft/clients"); err != nil {
		return err
	}
	return nil
}

func newBFTOrderer(node int) (*orderer.Server, error) {
	cfgPath := getBFTOrdererConfigPath(node)
	cfg, err := oc.LoadConfig(cfgPath)
	if err != nil {
		return nil, err
	}
	cfg.BlockChain.Path = getBFTOrdererPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Path = getBFTOrdererPath(node) + "/" + cfg.DB.LevelDB.Path
	cfg.Consensus.Tendermint.Path = getBFTOrdererPath(node) + "/" + cfg.Consensus.Tendermint.Path
	return orderer.NewServer(cfg)
}

func getBFTOrdererPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/.bft/orderers/%d", gopath, node)
}

func getBFTOrdererConfigPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/.bft/orderers/%d/orderer.yaml", gopath, node)
}

func getBFTClientPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/.bft/clients/%d", gopath, node)
}

func getBFTClientConfigPath(node int) string {
	return getBFTClientPath(node) + "/client.yaml"
}

func getBFTOrdererDataPath(node int) string {
	return fmt.Sprintf("%s/data", getBFTOrdererPath(node))
}

func getBFTOrdererBlockPath(node int) string {
	return fmt.Sprintf("%s/data/blocks", getBFTOrdererPath(node))
}
