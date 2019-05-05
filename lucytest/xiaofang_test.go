package lucytest

import (
	"fmt"
	cc "madledger/client/config"
	client "madledger/client/lib"
	cliu "madledger/client/util"
	oc "madledger/orderer/config"
	orderer "madledger/orderer/server"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/require"
)

var (
	bftOrderers [4]*orderer.Server
	bftClients  [4]*client.Client
)

func TestInitEnv(t *testing.T) {
	require.NoError(t, initBFTEnvironment())
}

// initBFTEnvironment will remove old test folders and copy necessary folders
func initBFTEnvironment() error {
	gopath := os.Getenv("GOPATH")
	if err := os.RemoveAll(gopath + "/src/madledger/tests/bft"); err != nil {
		return err
	}
	if err := copy.Copy(gopath+"/src/madledger/env/bft/.orderers", gopath+"/src/madledger/tests/bft/orderers"); err != nil {
		return err
	}
	if err := copy.Copy(gopath+"/src/madledger/env/bft/.clients", gopath+"/src/madledger/tests/bft/clients"); err != nil {
		return err
	}
	return nil
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
	time.Sleep(5 * time.Second)
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
	// client-0 create 4 channels
	client0 := bftClients[0]
	for i := 1; i <= 4; i++ {
		channel := "test" + strconv.Itoa(i)
		err := client0.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
	}

	time.Sleep(2 * time.Second)

	// then we will check if channels created by client-0 are create successful
	// query by client-0 and client-1
	require.NoError(t, listChannel())
}

func listChannel() error {
	client0 := bftClients[0]
	infos, err := client0.ListChannel(true)
	if err != nil {
		return err
	}
	table := cliu.NewTable()
	table.SetHeader("Name", "System", "BlockSize", "Identity")
	for _, info := range infos {
		table.AddRow(info.Name, info.System, info.BlockSize, info.Identity)
	}
	table.Render()

	client1 := bftClients[1]
	infos, err = client1.ListChannel(true)
	if err != nil {
		return err
	}
	table = cliu.NewTable()
	table.SetHeader("Name", "System", "BlockSize", "Identity")
	for _, info := range infos {
		table.AddRow(info.Name, info.System, info.BlockSize, info.Identity)
	}
	table.Render()

	return nil
}

// 关闭orderer 1，关闭期间通过client 0创建test5通道，然后重启orderer 1，查询数据
func TestNodeRestart(t *testing.T) {
	/*bftOrderers[1].Stop()

	//client 0创建test5通道
	client0 := bftClients[0]
	channel := "test5"
	err := client0.CreateChannel(channel, true, nil, nil)
	require.NoError(t, err)
	time.Sleep(2 * time.Second)

	fmt.Println("Restart orderer 1 ...")
	require.NoError(t, bftOrderers[1].Start())
	time.Sleep(5 * time.Second)

	require.NoError(t, listChannel())*/

	/*for i := range bftOrderers {
		require.True(t, util.IsDirSame(getBFTOrdererBlockPath(0), getBFTOrdererBlockPath(i)), fmt.Sprintf("Orderer %d is not same with 0", i))
	}*/

	bftOrderers[1].Stop()
	os.RemoveAll(getBFTOrdererDataPath(1))

	//client 0创建test5通道
	client0 := bftClients[0]
	channel := "test5"
	err := client0.CreateChannel(channel, true, nil, nil)
	require.NoError(t, err)
	time.Sleep(2 * time.Second)

	fmt.Println("Restart orderer 1 ...")
	server, err := newBFTOrderer(1)
	require.NoError(t, err)
	bftOrderers[1] = server
	go func(t *testing.T) {
		require.NoError(t, bftOrderers[1].Start())
	}(t)
	time.Sleep(5 * time.Second)

	require.NoError(t, listChannel())

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

func getBFTOrdererDataPath(node int) string {
	return fmt.Sprintf("%s/data", getBFTOrdererPath(node))
}

func getBFTOrdererPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/bft/orderers/%d", gopath, node)
}

func getBFTOrdererBlockPath(node int) string {
	return fmt.Sprintf("%s/data/blocks", getBFTOrdererPath(node))
}

func getBFTOrdererConfigPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/bft/orderers/%d/orderer.yaml", gopath, node)
}

func getBFTClientPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/bft/clients/%d", gopath, node)
}

func getBFTClientConfigPath(node int) string {
	return getBFTClientPath(node) + "/client.yaml"
}
