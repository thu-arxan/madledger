package testfor2clients_CfgRaft

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	cc "madledger/client/config"
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
	"regexp"
	"strings"
	"time"

	"go.etcd.io/etcd/raft/raftpb"

	"github.com/otiai10/copy"
	yaml "gopkg.in/yaml.v2"
)

var (
	raftOrderers [4]string
	raftClient   [3]*client.Client
	raftAdmin    *client.Client
	raftPeers    [4]string
)

// initBFTEnvironment will remove old test folders and copy necessary folders
func initRaftEnvironment() error {
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
	if err := os.RemoveAll(gopath + "/src/madledger/tests/raft"); err != nil {
		return err
	}
	if err := copy.Copy(gopath+"/src/madledger/env/raft/.orderers1", gopath+"/src/madledger/tests/raft/orderers"); err != nil {
		return err
	}
	if err := copy.Copy(gopath+"/src/madledger/env/raft/.clients", gopath+"/src/madledger/tests/raft/clients"); err != nil {
		return err
	}
	if err := copy.Copy(gopath+"/src/madledger/env/raft/.peers", gopath+"/src/madledger/tests/raft/peers"); err != nil {
		return err
	}

	for i := range raftOrderers {
		if err := absRaftOrdererConfig(i); err != nil {
			return err
		}
	}

	for i := range raftPeers {
		if err := absRaftPeerConfig(i); err != nil {
			return err
		}
	}

	if err := absRaftClientConfig("0"); err != nil {
		return err
	}
	if err := absRaftClientConfig("1"); err != nil {
		return err
	}
	if err := absRaftClientConfig("3"); err != nil {
		return err
	}
	if err := absRaftClientConfig("admin"); err != nil {
		return err
	}

	return nil
}

func loadClientConfig(cfgPath string) (*cc.Config, error) {
	cfgBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	var cfg cc.Config
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func absRaftClientConfig(node string) error {
	cfgPath := getRaftClientConfigPath(node)
	// load config
	cfg, err := loadClientConfig(cfgPath)
	if err != nil {
		return err
	}
	// change relative path into absolute path
	cfg.KeyStore.Keys[0] = getRaftClientPath(node) + "/" + cfg.KeyStore.Keys[0]
	cfg.TLS.CA = getRaftClientPath(node) + "/" + cfg.TLS.CA
	cfg.TLS.RawCert = getRaftClientPath(node) + "/" + cfg.TLS.RawCert
	cfg.TLS.Key = getRaftClientPath(node) + "/" + cfg.TLS.Key
	// rewrite peer config
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfgPath, data, os.ModePerm)
}

func absRaftOrdererConfig(node int) error {
	cfgPath := getRaftOrdererConfigPath(node)
	// load config
	cfg, err := loadOrdererConfig(cfgPath)
	if err != nil {
		return err
	}
	// change relative path into absolute path
	cfg.BlockChain.Path = getRaftOrdererPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Path = getRaftOrdererPath(node) + "/" + cfg.DB.LevelDB.Path
	cfg.Consensus.Raft.Path = getRaftOrdererPath(node) + "/" + cfg.Consensus.Raft.Path
	cfg.TLS.CA = getRaftOrdererPath(node) + "/" + cfg.TLS.CA
	cfg.TLS.RawCert = getRaftOrdererPath(node) + "/" + cfg.TLS.RawCert
	cfg.TLS.Key = getRaftOrdererPath(node) + "/" + cfg.TLS.Key
	// rewrite orderer config
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfgPath, data, os.ModePerm)
}

func loadOrdererConfig(cfgPath string) (*oc.Config, error) {
	cfgBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	var cfg oc.Config
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func absRaftPeerConfig(node int) error {
	cfgPath := getRaftPeerConfigPath(node)
	// load config
	cfg, err := loadPeerConfig(cfgPath)
	if err != nil {
		return err
	}
	// change relative path into absolute path
	cfg.BlockChain.Path = getRaftPeerPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Dir = getRaftPeerPath(node) + "/" + cfg.DB.LevelDB.Dir
	cfg.KeyStore.Key = getRaftPeerPath(node) + "/" + cfg.KeyStore.Key
	cfg.TLS.CA = getRaftPeerPath(node) + "/" + cfg.TLS.CA
	cfg.TLS.RawCert = getRaftPeerPath(node) + "/" + cfg.TLS.RawCert
	cfg.TLS.Key = getRaftPeerPath(node) + "/" + cfg.TLS.Key
	// rewrite peer config
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfgPath, data, os.ModePerm)
}

func loadPeerConfig(cfgPath string) (*pc.Config, error) {
	cfgBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	var cfg pc.Config
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
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
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("orderer start -c %s", getRaftOrdererConfigPath(node)))
		_, err := cmd.Output()
		if err != nil {
			fmt.Printf("Run orderer failed, because %s\n", err.Error())
			if !strings.Contains(err.Error(), "exit status") {
				panic(fmt.Sprintf("Run orderer failed, because %s\n", err.Error()))
			}
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
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("peer start -c %s", getRaftPeerConfigPath(node)))
		_, err := cmd.Output()
		if err != nil {
			fmt.Printf("Run peer failed, because %s\n", err.Error())
			if !strings.Contains(err.Error(), "exit status") {
				panic(fmt.Sprintf("Run peer failed, because %s\n", err.Error()))
			}
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
	cfgFilePath := getRaftPeerConfigPath(node)
	cfg, _ := pc.LoadConfig(cfgFilePath)

	cfg.BlockChain.Path = getRaftPeerPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Dir = getRaftPeerPath(node) + "/" + cfg.DB.LevelDB.Dir

	// then set key
	cfg.KeyStore.Key = getRaftPeerPath(node) + "/" + cfg.KeyStore.Key
	return cfg
}

func newBFTOrderer(node int) (*orderer.Server, error) {
	cfgPath := getRaftOrdererConfigPath(node)
	cfg, err := oc.LoadConfig(cfgPath)
	if err != nil {
		return nil, err
	}
	cfg.BlockChain.Path = getRaftOrdererPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Path = getRaftOrdererPath(node) + "/" + cfg.DB.LevelDB.Path
	cfg.Consensus.Tendermint.Path = getRaftOrdererPath(node) + "/" + cfg.Consensus.Tendermint.Path
	return orderer.NewServer(cfg)
}

func getBFTOrdererDataPath(node int) string {
	return fmt.Sprintf("%s/data", getRaftOrdererPath(node))
}

func getRaftOrdererPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/raft/orderers/%d", gopath, node)
}

func getBFTOrdererBlockPath(node int) string {
	return fmt.Sprintf("%s/data/blocks", getRaftOrdererPath(node))
}

func getRaftOrdererConfigPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/raft/orderers/%d/orderer.yaml", gopath, node)
}

func getRaftPeerPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/raft/peers/%d", gopath, node)
}

func getRaftPeerDataPath(node int) string {
	return fmt.Sprintf("%s/data", getRaftPeerPath(node))
}

func getRaftPeerConfigPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/raft/peers/%d/peer.yaml", gopath, node)
}

func getRaftClientPath(node string) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/raft/clients/%s", gopath, node)
}

func getRaftClientConfigPath(node string) string {
	return getRaftClientPath(node) + "/client.yaml"
}

func addNode(nodeID uint64, url string, channel string) error {
	// construct ConfChange
	cc, err := json.Marshal(raftpb.ConfChange{
		Type:    raftpb.ConfChangeAddNode,
		NodeID:  nodeID,
		Context: []byte(url),
	})
	if err != nil {
		return err
	}

	tx, err := coreTypes.NewTx(channel, coreTypes.CfgRaftAddress, cc, 0, "", raftAdmin.GetPrivKey())
	if err != nil {
		return err
	}

	status, err := raftAdmin.AddTx(tx)
	if err != nil {
		return err
	}
	// Then print the status
	table := cliu.NewTable()
	table.SetHeader("BlockNumber", "BlockIndex", "NodeAddOK")
	if status.Err != "" {
		table.AddRow(status.BlockNumber, status.BlockIndex, status.Err)
	} else {
		table.AddRow(status.BlockNumber, status.BlockIndex, "ok")
	}
	table.Render()
	return nil
}
func removeNode(nodeID uint64, channel string) error {
	// construct ConfChange
	cc, err := json.Marshal(raftpb.ConfChange{
		Type:   raftpb.ConfChangeRemoveNode,
		NodeID: nodeID,
	})
	if err != nil {
		return err
	}

	tx, err := coreTypes.NewTx(channel, coreTypes.CfgRaftAddress, cc, 0, "", raftAdmin.GetPrivKey())
	if err != nil {
		return err
	}

	status, err := raftAdmin.AddTx(tx)
	if err != nil {
		return err
	}
	// Then print the status
	table := cliu.NewTable()
	table.SetHeader("BlockNumber", "BlockIndex", "NodeRemoveOK")
	if status.Err != "" {
		table.AddRow(status.BlockNumber, status.BlockIndex, status.Err)
	} else {
		table.AddRow(status.BlockNumber, status.BlockIndex, "ok")
	}
	table.Render()
	return nil
}

func loadClient(node string, index int) error {
	clientPath := getRaftClientPath(node)
	cfgPath := getRaftClientConfigPath(node)
	cfg, err := cc.LoadConfig(cfgPath)
	if err != nil {
		return err
	}
	re, _ := regexp.Compile("^.*[.]keystore")
	for i := range cfg.KeyStore.Keys {
		cfg.KeyStore.Keys[i] = clientPath + "/.keystore" + re.ReplaceAllString(cfg.KeyStore.Keys[i], "")
	}
	client, err := client.NewClientFromConfig(cfg)
	if err != nil {
		return err
	}
	raftClient[index] = client
	return nil
}

func createContractForCallTx(channel string, node string, client *client.Client) error {
	// client create contract
	contractCodes, err := readCodes(getRaftClientPath(node) + "/MyTest.bin")
	if err != nil {
		return err
	}
	tx, err := coreTypes.NewTx(channel, common.ZeroAddress, contractCodes, 0, "", client.GetPrivKey())
	if err != nil {
		return err
	}
	_, err = client.AddTx(tx)
	if err != nil {
		return err
	}

	return nil
}

func compareChannels(client1 *client.Client, client2 *client.Client) error {
	infos1, err := client1.ListChannel(true)
	if err != nil {
		return err
	}
	infos2, err := client2.ListChannel(true)
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
			return fmt.Errorf("the blockSize is not consistent, %d in client1, %d in client2", infos1[i].BlockSize, infos2[i].BlockSize)
		}
	}

	fmt.Println("CompareChannels: channels between two orderers are the same.")
	return nil
}

func getNumForCallTx(node string, num string) error {
	abiPath := fmt.Sprintf(getRaftClientPath("0") + "/MyTest.abi")
	var inputs = make([]string, 0)
	payloadBytes, err := abi.GetPayloadBytes(abiPath, "getNum", inputs)
	if err != nil {
		return err
	}

	channel := "test" + node
	tx, err := coreTypes.NewTx(channel, common.HexToAddress("0x8de6ce45b289502e16aef93313fd3082993acb1f"), payloadBytes,
		0, "", raftClient[0].GetPrivKey())
	if err != nil {
		return err
	}

	status, err := raftClient[0].AddTx(tx)
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
	abiPath := fmt.Sprintf(getRaftClientPath("0") + "/MyTest.abi")
	inputs := []string{num}
	payloadBytes, err := abi.GetPayloadBytes(abiPath, "setNum", inputs)
	if err != nil {
		return err
	}

	channel := "test" + node
	tx, err := coreTypes.NewTx(channel, common.HexToAddress("0x8de6ce45b289502e16aef93313fd3082993acb1f"), payloadBytes,
		0, "", raftClient[0].GetPrivKey())
	if err != nil {
		return err
	}

	_, err = raftClient[0].AddTx(tx)
	if err != nil {
		return err
	}
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
