package testfor2clients_bft

import (
	"fmt"
	cc "madledger/client/config"
	client "madledger/client/lib"
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
)

func TestInitEnv(t *testing.T) {
	require.NoError(t, initBFTEnvironment())
}

// initBFTEnvironment will remove old test folders and copy necessary folders
func initBFTEnvironment() error {
	gopath := os.Getenv("GOPATH")
	if err := os.RemoveAll(gopath + "/src/madledger/tests/raft"); err != nil {
		return err
	}
	if err := copy.Copy(gopath+"/src/madledger/env/raft/.orderers", gopath+"/src/madledger/tests/raft/orderers"); err != nil {
		return err
	}
	if err := copy.Copy(gopath+"/src/madledger/env/raft/.clients", gopath+"/src/madledger/tests/raft/clients"); err != nil {
		return err
	}
	if err := copy.Copy(gopath+"/src/madledger/env/raft/.peers", gopath+"/src/madledger/tests/peers/clients"); err != nil {
		return err
	}
	for i := range raftOrderers {
		if err := absRAFTOrdererConfig(i); err != nil {
			return err
		}
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
	for i := range raftClients {
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
		raftClients[i] = client
	}
}

func TestBFTCreateChannels(t *testing.T) {
	// client 0 and client 1 create channels concurrently
	client0 := raftClients[0]
	var channels []string
	for m := 0; m < 8; m++ {
		if m == 3 { // stop orderer0
			go func(t *testing.T) {
				bftOrderers[0].Stop()
				require.NoError(t, os.RemoveAll(getBFTOrdererDataPath(0)))
			}(t)
		}
		if m == 5 { // restart orderer0

				fmt.Println("Restart orderer 0 ...")
				server, err := newBFTOrderer(0)
				require.NoError(t, err)
				bftOrderers[0] = server
				err = bftOrderers[0].Start()
				require.NoError(t, err)
		}
		// client 0 create channel
		channel := "test" + strconv.Itoa(m)
		channels = append(channels, channel)
		err := client0.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)

	}
	time.Sleep(2 * time.Second)

	// then we will check if channels are create successful
	require.NoError(t, compareChannels(channels))
}

/*func TestTendermintDB(t *testing.T) {
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
}*/

func newBFTOrderer(node int) (*orderer.Server, error) {
	cfgPath := getBFTOrdererConfigPath(node)
	cfg, err := oc.LoadConfig(cfgPath)
	if err != nil {
		return nil, err
	}
	return orderer.NewServer(cfg)
}

func getBFTOrdererDataPath(node int) string {
	return fmt.Sprintf("%s/data", getBFTOrdererPath(node))
}

func getBFTOrdererPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/raft/orderers/%d", gopath, node)
}

func getBFTOrdererBlockPath(node int) string {
	return fmt.Sprintf("%s/data/blocks", getBFTOrdererPath(node))
}

func getBFTOrdererConfigPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/raft/orderers/%d/orderer.yaml", gopath, node)
}

func getBFTClientPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/raft/clients/%d", gopath, node)
}

func getBFTClientConfigPath(node int) string {
	return getBFTClientPath(node) + "/client.yaml"
}
