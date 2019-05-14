package testfor1client

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	cc "madledger/client/config"
	client "madledger/client/lib"
	cliu "madledger/client/util"
	"madledger/common"
	"madledger/common/abi"
	"madledger/common/util"
	"madledger/core/types"
	oc "madledger/orderer/config"
	orderer "madledger/orderer/server"
	pc "madledger/peer/config"
	peer "madledger/peer/server"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

var (
	bftOrderers [4]string
	bftClients  [4]*client.Client
	bftPeers    [4]*peer.Server
)

func TestInitEnv(t *testing.T) {
	require.NoError(t, initBFTEnvironment())
}

// initBFTEnvironment will remove old test folders and copy necessary folders
func initBFTEnvironment() error {
	// kill all Orderers
	pids := getOrderersPid()
	for _, pid := range pids {
		stopOrderer(pid)
	}

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
	if err := copy.Copy(gopath+"/src/madledger/env/bft/.peers", gopath+"/src/madledger/tests/bft/peers"); err != nil {
		return err
	}

	for i := range bftOrderers {
		if err := absBFTOrdererConfig(i); err != nil {
			return err
		}
	}

	return nil
}

func TestBFTOrdererStart(t *testing.T) {
	// then we can run orderers
	for i := range bftOrderers {
		pid := startOrderer(i)
		bftOrderers[i] = pid
	}
}

func TestBFTPeersStart(t *testing.T) {
	for i := 0; i < 4; i++ {
		cfg := getPeerConfig(i)
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

	time.Sleep(5 * time.Second)
}

func getPeerConfig(node int) *pc.Config {
	cfgFilePath := getBFTPeerConfigPath(node)
	cfg, _ := pc.LoadConfig(cfgFilePath)

	cfg.BlockChain.Path = getBFTPeerPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Dir = getBFTPeerPath(node) + "/" + cfg.DB.LevelDB.Dir

	// then set key
	cfg.KeyStore.Key = getBFTPeerPath(node) + "/" + cfg.KeyStore.Key
	return cfg
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
	for i := 0; i <= 2; i++ {
		channel := "test" + strconv.Itoa(i)
		err := client0.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
	}
	time.Sleep(2 * time.Second)

	// then we will check if channels created by client-0 are create successful
	// query by client-1
	require.NoError(t, listChannel(1))
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

// 关闭orderer 1，关闭期间通过client 0创建test3通道，然后重启orderer 1，查询数据
func TestBFTNodeRestart(t *testing.T) {
	stopOrderer(bftOrderers[1])
	os.RemoveAll(getBFTOrdererDataPath(1))

	//client 0创建test3通道
	client0 := bftClients[0]
	channel := "test3"
	err := client0.CreateChannel(channel, true, nil, nil)
	require.NoError(t, err)
	time.Sleep(2 * time.Second)

	fmt.Println("Restart orderer 1 ...")
	bftOrderers[1] = startOrderer(1)
	time.Sleep(2 * time.Second)

	// query by client 1
	require.NoError(t, listChannel(1))

}

func TestBFTCreateChannelAfterRestart(t *testing.T) {
	//client 1创建通道test4
	client1 := bftClients[1]
	channel := "test4"
	err := client1.CreateChannel(channel, true, nil, nil)
	require.NoError(t, err)
	time.Sleep(2 * time.Second)

	// query by client 1
	require.NoError(t, listChannel(1))
}

func readCodes(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(string(data))
}

func TestBFTCreateTxAfterRestart(t *testing.T) {
	//client 1创建智能合约
	contractCodes, err := readCodes(getBFTClientPath(1) + "/MyTest.bin")
	require.NoError(t, err)
	client := bftClients[1]
	tx, err := types.NewTx("test4", common.ZeroAddress, contractCodes, client.GetPrivKey())
	require.NoError(t, err)

	status, err := client.AddTx(tx)
	require.NoError(t, err)

	// Then print the status
	table := cliu.NewTable()
	table.SetHeader("BlockNumber", "BlockIndex", "ContractAddress")
	if status.Err != "" {
		table.AddRow(status.BlockNumber, status.BlockIndex, status.Err)
	} else {
		table.AddRow(status.BlockNumber, status.BlockIndex, status.ContractAddress)
	}
	table.Render()
}

func TestBFTCallTxAfterRestart(t *testing.T) {
	//client 1调用智能合约
	abiPath := fmt.Sprintf(getBFTClientPath(1) + "/MyTest.abi")
	funcName := "getNum"
	var inputs []string = make([]string, 0)
	payloadBytes, err := abi.GetPayloadBytes(abiPath, funcName, inputs)
	require.NoError(t, err)

	client := bftClients[1]
	tx, err := types.NewTx("test4", common.HexToAddress("0x0619e2393802cc99e90cf892b92a113f19af5887"), payloadBytes, client.GetPrivKey())
	require.NoError(t, err)

	status, err := client.AddTx(tx)
	require.NoError(t, err)

	// Then print the status
	table := cliu.NewTable()
	table.SetHeader("BlockNumber", "BlockIndex", "Output")
	if status.Err != "" {
		table.AddRow(status.BlockNumber, status.BlockIndex, status.Err)
	} else {
		values, err := abi.Unpacker(abiPath, funcName, status.Output)
		require.NoError(t, err)
		var output []string
		for _, value := range values {
			output = append(output, value.Value)
		}
		table.AddRow(status.BlockNumber, status.BlockIndex, output)
	}
	table.Render()
}

func TestBFTEnd(t *testing.T) {
	for _, pid := range bftOrderers {
		stopOrderer(pid)
	}
	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/bft"))
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

// startOrderer run orderer and return pid
func startOrderer(node int) string {
	before := getOrderersPid()
	go func() {
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("orderer start -c %s", getBFTOrdererConfigPath(node)))
		_, err := cmd.Output()
		if err != nil {
			panic("Run orderer failed")
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

func getBFTPeerPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/bft/peers/%d", gopath, node)
}
func getBFTPeerConfigPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/bft/peers/%d/peer.yaml", gopath, node)
}

func getBFTClientPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/bft/clients/%d", gopath, node)
}

func getBFTClientConfigPath(node int) string {
	return getBFTClientPath(node) + "/client.yaml"
}
