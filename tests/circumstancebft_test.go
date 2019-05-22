package tests

import (
	"fmt"
	"io/ioutil"
	cc "madledger/client/config"
	client "madledger/client/lib"
	"madledger/common"
	"madledger/common/util"
	"madledger/core/types"
	oc "madledger/orderer/config"
	pc "madledger/peer/config"
	peer "madledger/peer/server"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

/*
* This test will start from a empty environment and start some orderers support bft consensus.
* This test will include operations below.
 */

var (
	// bftOrderers store the pid of orderers
	bftOrderers [4]string
	bftClients  [4]*client.Client
	bftPeers    [4]*peer.Server
	bftChannels []string
)

func TestBFT(t *testing.T) {
	require.NoError(t, initBFTEnvironment())
}

func TestBFTRun(t *testing.T) {
	for i := range bftOrderers {
		pid := startOrderer(i)
		bftOrderers[i] = pid
	}
}

func TestBFTPeersStart(t *testing.T) {
	for i := 0; i < 4; i++ {
		cfg := getBFTPeerConfig(i)
		server, err := peer.NewServer(cfg)
		require.NoError(t, err)
		bftPeers[i] = server
	}

	for i := range bftPeers {
		go func(t *testing.T, i int) {
			err := bftPeers[i].Start()
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
	stopOrderer(bftOrderers[1])
	os.RemoveAll(getBFTOrdererDataPath(1))
	bftOrderers[1] = startOrderer(1)
	time.Sleep(2000 * time.Millisecond)
	for i := range bftOrderers {
		require.True(t, util.IsDirSame(getBFTOrdererBlockPath(0), getBFTOrdererBlockPath(i)), fmt.Sprintf("Orderer %d is not same with 0", i))
	}
}

func TestBFTReCreateChannels(t *testing.T) {
	// Here we recreate 2 channels
	for i := 0; i < 2; i++ {
		channel := strings.ToLower(util.RandomString(16))
		err := bftClients[0].CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
	}
	time.Sleep(2 * time.Second)
	var prevChannels []string
	// Then we will list channels
	for i := range bftClients {
		infos, err := bftClients[i].ListChannel(false)
		require.NoError(t, err)
		var channels []string
		for i := range infos {
			channels = append(channels, infos[i].Name)
		}
		if len(prevChannels) != 0 {
			require.Equal(t, prevChannels, channels)
		}
		prevChannels = channels
		bftChannels = channels
	}
}
func TestBFTCreateTx(t *testing.T) {
	client0 := bftClients[0]
	client1 := bftClients[1]
	for m := 1; m <= 6; m++ {
		if m == 3 { // stop orderer0
			stopOrderer(bftOrderers[0])
			require.NoError(t, os.RemoveAll(getBFTOrdererDataPath(0)))
		}
		if m == 4 { // restart orderer0
			bftOrderers[0] = startOrderer(0)
		}
		// client 0 create contract
		contractCodes, err := readCodes(getBFTClientPath(1) + "/MyTest.bin")
		require.NoError(t, err)
		channel := bftChannels[util.RandNum(len(bftChannels))]
		fmt.Printf("Create contract on channel %s(m=%d) ...\n", channel, m)
		tx, err := types.NewTx(channel, common.ZeroAddress, contractCodes, client0.GetPrivKey())
		require.NoError(t, err)
		_, err = client0.AddTx(tx)
		require.NoError(t, err)

		// client 1 create channel
		contractCodes, err = readCodes(getBFTClientPath(1) + "/MyTest.bin")
		require.NoError(t, err)
		channel = bftChannels[util.RandNum(len(bftChannels))]
		fmt.Printf("Create contract on channel %s(m=%d) ...\n", channel, m)
		tx, err = types.NewTx(channel, common.ZeroAddress, contractCodes, client1.GetPrivKey())
		require.NoError(t, err)
		_, err = client1.AddTx(tx)
		require.NoError(t, err)
	}
	time.Sleep(1000 * time.Millisecond)
	for i := range bftOrderers {
		require.True(t, util.IsDirSame(getBFTOrdererBlockPath(0), getBFTOrdererBlockPath(i)), fmt.Sprintf("Orderer %d is not same with 0", i))
	}
}
func TestBFTEnd(t *testing.T) {
	for i := range bftOrderers {
		stopOrderer(bftOrderers[i])
	}
	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/.bft"))
}

// initBFTEnvironment will remove old test folders and copy necessary folders, also it will stop orderers on running
func initBFTEnvironment() error {
	// kill all Orderers
	cmd := exec.Command("/bin/sh", "-c", "pidof orderer")
	output, err := cmd.Output()
	if err == nil {
		pids := strings.Split(string(output), " ")
		for _, pid := range pids {
			stopOrderer(pid)
		}
	}

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
	if err := copy.Copy(gopath+"/src/madledger/env/bft/.peers", gopath+"/src/madledger/tests/.bft/peers"); err != nil {
		return err
	}

	for i := range bftOrderers {
		if err := absBFTOrdererConfig(i); err != nil {
			return err
		}
	}
	return nil
}

func absBFTOrdererConfig(node int) error {
	cfgPath := getBFTOrdererConfigPath(node)
	cfg, err := oc.LoadConfig(cfgPath)
	if err != nil {
		return err
	}
	cfg.BlockChain.Path = getBFTOrdererPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Path = getBFTOrdererPath(node) + "/" + cfg.DB.LevelDB.Path
	cfg.Consensus.Tendermint.Path = getBFTOrdererPath(node) + "/" + cfg.Consensus.Tendermint.Path

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfgPath, data, os.ModePerm)
}

// startOrderer run orderer and return pid
func startOrderer(node int) string {
	before := getOrderersPid()
	go func() {
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("orderer start -c %s", getBFTOrdererConfigPath(node)))
		_, err := cmd.Output()
		if err != nil {
			panic(fmt.Sprintf("Run orderer %d failed", node))
		}
	}()

	for {
		after := getOrderersPid()
		if len(after) != len(before) {
			for _, pid := range after {
				if !util.Contain(before, pid) {
					return pid
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func getOrderersPid() []string {
	cmd := exec.Command("/bin/sh", "-c", "pidof orderer")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	pids := strings.Split(string(output), " ")
	return pids
}

// stopOrderer stop an orderer
func stopOrderer(pid string) {
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("kill -TERM %s", pid))
	cmd.Output()
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

func getBFTPeerPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/.bft/peers/%d", gopath, node)
}
func getBFTPeerConfigPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/.bft/peers/%d/peer.yaml", gopath, node)
}

func getBFTPeerConfig(node int) *pc.Config {
	cfgFilePath := getBFTPeerConfigPath(node)
	cfg, _ := pc.LoadConfig(cfgFilePath)

	cfg.BlockChain.Path = getBFTPeerPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Dir = getBFTPeerPath(node) + "/" + cfg.DB.LevelDB.Dir

	// then set key
	cfg.KeyStore.Key = getBFTPeerPath(node) + "/" + cfg.KeyStore.Key
	return cfg
}