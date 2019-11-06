package testfor1client_bft

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint/abci/types"
	"io/ioutil"
	client "madledger/client/lib"
	cliu "madledger/client/util"
	"madledger/common"
	"madledger/common/abi"
	"madledger/common/util"
	coreTypes "madledger/core/types"
	oc "madledger/orderer/config"
	orderer "madledger/orderer/server"
	pc "madledger/peer/config"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/otiai10/copy"
	"gopkg.in/yaml.v2"
)

var (
	bftOrderers [4]string
	// just 1 is enough, we set 2
	bftClient *client.Client
	bftAdmin  *client.Client
	//bftPeers  [4]*peer.Server
	bftPeers  [4]string
)

// initBFTEnvironment will remove old test folders and copy necessary folders
func initBFTEnvironment() error {
	// kill all Orderers
	pids := getOrderersPid()
	for _, pid := range pids {
		stopOrderer(pid)
	}

	// kill all Peers
	pids = getPeerPid()
	for _, pid := range pids {
		stopPeer(pid)
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

	for i := range bftPeers {
		if err := absBFTPeerConfig(i); err != nil {
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

func absBFTPeerConfig(node int) error {
	cfgPath := getBFTPeerConfigPath(node)
	cfg, err := pc.LoadConfig(cfgPath)
	if err != nil {
		return err
	}
	cfg.BlockChain.Path = getBFTPeerPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Dir = getBFTPeerPath(node) + "/" + cfg.DB.LevelDB.Dir
	cfg.KeyStore.Key = getBFTPeerPath(node) + "/" + cfg.KeyStore.Key

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

func readCodes(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(string(data))
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

func getPeerPid() []string {
	cmd := exec.Command("/bin/sh", "-c", "pidof peer")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	pids := strings.Split(string(output), " ")
	return pids
}

// stopPeer stop a peer
func stopPeer(pid string) {
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("kill -TERM %s", pid))
	cmd.Output()
}

// startPeer run peer and return pid
func startPeer(node int) string {
	before := getPeerPid()
	go func() {
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("peer start -c %s", getBFTPeerConfigPath(node)))
		_, err := cmd.Output()
		if err != nil {
			panic("Run peer failed")
		}
	}()

	for {
		after := getPeerPid()
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

func getPeerConfig(node int) *pc.Config {
	cfgFilePath := getBFTPeerConfigPath(node)
	cfg, _ := pc.LoadConfig(cfgFilePath)

	cfg.BlockChain.Path = getBFTPeerPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Dir = getBFTPeerPath(node) + "/" + cfg.DB.LevelDB.Dir

	// then set key
	cfg.KeyStore.Key = getBFTPeerPath(node) + "/" + cfg.KeyStore.Key
	return cfg
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

func getBFTPeerDataPath(node int) string {
	return fmt.Sprintf("%s/data", getBFTPeerPath(node))
}

func getBFTPeerConfigPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/bft/peers/%d/peer.yaml", gopath, node)
}

func getBFTClientPath(node string) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/bft/clients/%s", gopath, node)
}

func getBFTClientConfigPath(node string) string {
	return getBFTClientPath(node) + "/client.yaml"
}

func addOrRemoveNode(pubKey string, power int64,channel string) error {
	// construct PubKey
	data, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return err
	}
	pubkey := types.PubKey{
		Type: "ed25519",
		Data: data,
	}

	// construct ValidatorUpdate
	validatorUpdate, err := json.Marshal(types.ValidatorUpdate{
		PubKey: pubkey,
		Power:  power,
	})
	tx, err := coreTypes.NewTx(channel, coreTypes.CfgTendermintAddress, validatorUpdate, bftAdmin.GetPrivKey(), coreTypes.VALIDATOR)
	if err != nil {
		return err
	}

	status, err := bftAdmin.AddTx(tx)
	if err != nil {
		return err
	}
	// Then print the status
	table := cliu.NewTable()
	table.SetHeader("BlockNumber", "BlockIndex", "ValidatorAddOk")
	if status.Err != "" {
		table.AddRow(status.BlockNumber, status.BlockIndex, status.Err)
	} else {
		table.AddRow(status.BlockNumber, status.BlockIndex, "ok")
	}
	table.Render()
	return nil
}

func listChannel(client *client.Client) error {
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


func createChannelForCallTx() error {
	// client 0 create channel
	err := bftClient.CreateChannel("test0", true, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func createContractForCallTx(node string) error {
	// client 0 create contract
	contractCodes, err := readCodes(getBFTClientPath(node) + "/MyTest.bin")
	if err != nil {
		return err
	}
	tx, err := coreTypes.NewTx("test"+node, common.ZeroAddress, contractCodes, bftClient.GetPrivKey(), coreTypes.NORMAL)
	if err != nil {
		return err
	}
	_, err = bftClient.AddTx(tx)
	if err != nil {
		return err
	}

	return nil
}

func compareChannels() error {
	infos1, err := bftClient.ListChannel(true)
	if err != nil {
		return err
	}
	infos2, err := bftAdmin.ListChannel(true)
	if err != nil {
		return err
	}
	if len(infos1) != len(infos2) {
		return fmt.Errorf("the number is not consistent")
	}
	for i := range infos1 {
		if infos1[i].Name != infos2[i].Name {
			return fmt.Errorf("the name is not consistent")
		}
		if infos1[i].BlockSize != infos2[i].BlockSize {
			return fmt.Errorf("the blockSize is not consistent, %d in client, %d in admin", infos1[i].BlockSize, infos2[i].BlockSize)
		}
	}

	fmt.Println("CompareChannels: channels between two orderers are the same.")
	return nil
}

func getNumForCallTx(node string, num string) error {
	abiPath := fmt.Sprintf(getBFTClientPath("0") + "/MyTest.abi")
	var inputs = make([]string, 0)
	payloadBytes, err := abi.GetPayloadBytes(abiPath, "getNum", inputs)
	if err != nil {
		return err
	}

	channel := "test" + node
	tx, err := coreTypes.NewTx(channel, common.HexToAddress("0x8de6ce45b289502e16aef93313fd3082993acb1f"), payloadBytes,
		bftClient.GetPrivKey(), coreTypes.NORMAL)
	if err != nil {
		return err
	}

	status, err := bftClient.AddTx(tx)
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
		return fmt.Errorf("call tx on channel %s: setNum expect %s but receive %s", channel, num, output[0])
	}
	return nil
}

func setNumForCallTx(node string, num string) error {
	abiPath := fmt.Sprintf(getBFTClientPath("0") + "/MyTest.abi")
	inputs := []string{num}
	payloadBytes, err := abi.GetPayloadBytes(abiPath, "setNum", inputs)
	if err != nil {
		return err
	}

	channel := "test" + node
	tx, err := coreTypes.NewTx(channel, common.HexToAddress("0x8de6ce45b289502e16aef93313fd3082993acb1f"), payloadBytes,
		bftClient.GetPrivKey(), coreTypes.NORMAL)
	if err != nil {
		return err
	}

	_, err = bftClient.AddTx(tx)
	if err != nil {
		return err
	}
	return nil
}
