package testfor2clients_bft

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	client "madledger/client/lib"
	"madledger/common"
	"madledger/common/abi"
	"madledger/common/util"
	"madledger/core/types"
	oc "madledger/orderer/config"
	pc "madledger/peer/config"
	peer "madledger/peer/server"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/otiai10/copy"
	"gopkg.in/yaml.v2"
)

var (
	bftOrderers [4]string
	bftClients  [4]*client.Client
	bftPeers    [4]*peer.Server
)

// initBFTEnvironment will remove old test folders and copy necessary folders
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

func createChannelForCallTx() error {
	// client 0 create channel
	client := bftClients[0]
	err := client.CreateChannel("test0", true, nil, nil)
	if err != nil {
		return err
	}

	// client 1 create channel
	client = bftClients[1]
	err = client.CreateChannel("test1", true, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func initPeer(node int) error {
	cfg := getPeerConfig(node)
	server, err := peer.NewServer(cfg)
	if err != nil {
		return err
	}
	bftPeers[node] = server
	return nil
}

func createContractForCallTx() error {
	// client 0 create contract
	client := bftClients[0]
	contractCodes, err := readCodes(getBFTClientPath(0) + "/MyTest.bin")
	if err != nil {
		return err
	}
	tx, err := types.NewTx("test0", common.ZeroAddress, contractCodes, client.GetPrivKey())
	if err != nil {
		return err
	}
	_, err = client.AddTx(tx)
	if err != nil {
		return err
	}

	// client 1 create channel
	client = bftClients[1]
	contractCodes, err = readCodes(getBFTClientPath(1) + "/MyTest.bin")
	if err != nil {
		return err
	}
	tx, err = types.NewTx("test1", common.ZeroAddress, contractCodes, client.GetPrivKey())
	if err != nil {
		return err
	}

	_, err = client.AddTx(tx)
	if err != nil {
		return err
	}
	return nil
}

func setNumForCallTx(node int, num string) error {
	abiPath := fmt.Sprintf(getBFTClientPath(node) + "/MyTest.abi")
	inputs := []string{num}
	payloadBytes, err := abi.GetPayloadBytes(abiPath, "setNum", inputs)
	if err != nil {
		return err
	}

	client := bftClients[node]
	channel := "test" + strconv.Itoa(node)
	addr := ""
	if node == 0 {
		addr = "0x8de6ce45b289502e16aef93313fd3082993acb1f"
	} else {
		addr = "0x1b66001e01d3c8d3893187fee59e3bea1d9bdd9b"
	}
	tx, err := types.NewTx(channel, common.HexToAddress(addr), payloadBytes, client.GetPrivKey())
	if err != nil {
		return err
	}

	_, err = client.AddTx(tx)
	if err != nil {
		return err
	}
	return nil
}

func getNumForCallTx(node int, num string) error {
	abiPath := fmt.Sprintf(getBFTClientPath(node) + "/MyTest.abi")
	var inputs []string = make([]string,0)
	payloadBytes, err := abi.GetPayloadBytes(abiPath, "getNum", inputs)
	if err != nil {
		return err
	}

	client := bftClients[node]
	channel := "test" + strconv.Itoa(node)
	addr := ""
	if node == 0 {
		addr = "0x8de6ce45b289502e16aef93313fd3082993acb1f"
	} else {
		addr = "0x1b66001e01d3c8d3893187fee59e3bea1d9bdd9b"
	}
	tx, err := types.NewTx(channel, common.HexToAddress(addr), payloadBytes, client.GetPrivKey())
	if err != nil {
		return err
	}

	status, err := client.AddTx(tx)
	if err != nil {
		return err
	}

	values, err := abi.Unpacker(abiPath, "getNum", status.Output)
	if err != nil {
		return err
	}

	var output []string
	for _, value := range values {
		output = append(output, value.Value)
	}
	if output[0] != num {
		return fmt.Errorf("call tx on channel test%d: getNum expect %s but receive %s", node, num, output[0])
	}
	return nil
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
			if infos[i].Name != "_config" && infos[i].Name != "_global" {
				if !util.Contain(channels, infos[i].Name) {
					return fmt.Errorf("channel name doesn't exit in channels")
				}
			}
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

func getOrderersPid() []string {
	cmd := exec.Command("/bin/sh", "-c", "pidof orderer")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	pids := strings.Split(string(output), " ")
	return pids
}

func getBFTPeerDataPath(node int) string {
	return fmt.Sprintf("%s/data", getBFTPeerPath(node))
}


// stopOrderer stop an orderer
func stopOrderer(pid string) {
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("kill -TERM %s", pid))
	cmd.Output()
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

func getBFTPeerPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/bft/peers/%d", gopath, node)
}
func getBFTPeerConfigPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/bft/peers/%d/peer.yaml", gopath, node)
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

func readCodes(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(string(data))
}
