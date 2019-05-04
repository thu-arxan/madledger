package lucytest

import (
	"fmt"
	"github.com/stretchr/testify/require"
	cc "madledger/client/config"
	client "madledger/client/lib"
	oc "madledger/orderer/config"
	orderer "madledger/orderer/server"
	cliu "madledger/client/util"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"
)

var (
	bftOrderers [4]*orderer.Server
	bftClients  [4]*client.Client
)

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
	var channels []string
	// client-0 create 4 channels
	client0 := bftClients[0]
	for i := 1; i <= 4; i++ {
		channel := "test" + strconv.Itoa(i)
		channels = append(channels, channel)
		err := client0.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
	}

	time.Sleep(2 * time.Second)

	// then we will check if channels created by client-0 are create successful
	// query by client-0 and client-1
	infos, err := client0.ListChannel(false)
	require.NoError(t, err)
	table := cliu.NewTable()
	table.SetHeader("Name", "System", "BlockSize", "Identity")
	for _, info := range infos {
		table.AddRow(info.Name, info.System, info.BlockSize, info.Identity)
	}
	table.Render()

	client1 := bftClients[1]
	infos, err = client1.ListChannel(false)
	require.NoError(t, err)
	table = cliu.NewTable()
	table.SetHeader("Name", "System", "BlockSize", "Identity")
	for _, info := range infos {
		table.AddRow(info.Name, info.System, info.BlockSize, info.Identity)
	}
	table.Render()
}

func getBFTClientPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/env/bft/clients/%d", gopath, node)
}

func getBFTClientConfigPath(node int) string {
	return getBFTClientPath(node) + "/client.yaml"
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

func getBFTOrdererConfigPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/env/bft/orderers/%d/orderer.yaml", gopath, node)
}

func getBFTOrdererPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/env/bft/orderers/%d", gopath, node)
}
