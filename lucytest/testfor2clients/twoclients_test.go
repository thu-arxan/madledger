package testfor2clients

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	cc "madledger/client/config"
	client "madledger/client/lib"
	cliu "madledger/client/util"
	comu "madledger/common/util"
	oc "madledger/orderer/config"
	orderer "madledger/orderer/server"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	// client 0 and client 1 create channels concurrently
	client0 := bftClients[0]
	client1 := bftClients[1]
	var channels []string
	for m := 1; m <= 8; m++ {
		if m == 4 { // stop orderer0
			go func(t *testing.T) {
				bftOrderers[0].Stop()
				require.NoError(t, os.RemoveAll(getBFTOrdererDataPath(0)))
			}(t)
		}
		if m == 6 { // restart orderer0
			go func(t *testing.T) {
				fmt.Println("Restart orderer 0 ...")
				server, err := newBFTOrderer(0)
				require.NoError(t, err)

				bftOrderers[0] = server
				err = bftOrderers[0].Start()
				require.NoError(t, err)
			}(t)
		}
		// client 0 create channel
		channel := "test" + strconv.Itoa(m) + "0"
		channels = append(channels, channel)
		err := client0.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)

		// client 1 create channel
		channel = "test" + strconv.Itoa(m) + "1"
		channels = append(channels, channel)
		err = client1.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)

	}
	time.Sleep(20 * time.Second)

	// then we will check if channels are create successful
	require.NoError(t, compareChannels(channels))
}

func compareChannels(channels []string) error {
	lenChannels := len(channels) + 2
	for i := 0; i < 2; i++ {
		client := bftClients[i]
		infos, err := client.ListChannel(true)
		if err != nil {
			return err
		}

		if len(infos) != lenChannels {
			return fmt.Errorf("the number is not consistent")
		}

		for i := range infos {
			if (infos[i].Name != "_config" && infos[i].Name != "_global") && !comu.Contain(channels, infos[i].Name) {
				return fmt.Errorf("channel name doesn't exit in channels")
			}
		}

	}

	return nil
}

func listChannel(node int) error {
	client := bftClients[node]
	infos, err := client.ListChannel(true)
	if err != nil {
		return err
	}
	table := cliu.NewTable()
	table.SetHeader("Name", "System", "BlockSize", "Identity")
	for _, info := range infos {
		table.AddRow(info.Name, info.System, info.BlockSize, info.Identity)
	}
	table.Render()

	return nil
}

func TestTendermintDB(t *testing.T) {
	for i := 0; i < 2; i++ {
		//path := fmt.Sprint(getBFTOrdererPath(i)+"/.tendermint/.glue")
		path := fmt.Sprintf("/home/hadoop/GOPATH/src/madledger/tests/bft/orderers/%d/data/leveldb", i)
		db, err := leveldb.OpenFile(path, nil)
		require.NoError(t, err)
		iter := db.NewIterator(nil, nil)
		// 遍历key-value
		fmt.Println("Get glue db from ", path)
		for iter.Next() {
			key := string(iter.Key())
			value := iter.Value()
			if strings.HasPrefix(string(key), "number") {
				number, _ := comu.BytesToUint64(value)
				fmt.Println(string(key), ", ", number)
			} else {
				fmt.Println(string(key), ", ", string(value))
			}
		}
		err = iter.Error()
		require.NoError(t, err)
		iter.Release()
		db.Close()
	}
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
