package testfor1client_bft

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	cc "madledger/client/config"
	client "madledger/client/lib"
	"madledger/common"
	"madledger/common/abi"
	"madledger/common/util"
	"madledger/core/types"
	oc "madledger/orderer/config"
	orderer "madledger/orderer/server"
	pc "madledger/peer/config"
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
	// just 1 is enough, we set 2
	bftClients [2]*client.Client
	bftPeers   [4]string
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

	for i := range bftClients {
		if err := absBFTClientConfig(i); err != nil {
			return err
		}
	}

	return nil
}

func absBFTPeerConfig(node int) error {
	cfgPath := getBFTPeerConfigPath(node)
	// load config
	cfg, err := loadPeerConfig(cfgPath)
	if err != nil {
		return err
	}
	// change relative path into absolute path
	cfg.BlockChain.Path = getBFTPeerPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Dir = getBFTPeerPath(node) + "/" + cfg.DB.LevelDB.Dir
	cfg.KeyStore.Key = getBFTPeerPath(node) + "/" + cfg.KeyStore.Key
	cfg.TLS.CA = getBFTPeerPath(node) + "/" + cfg.TLS.CA
	cfg.TLS.RawCert = getBFTPeerPath(node) + "/" + cfg.TLS.RawCert
	cfg.TLS.Key = getBFTPeerPath(node) + "/" + cfg.TLS.Key
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

func absBFTOrdererConfig(node int) error {
	cfgPath := getBFTOrdererConfigPath(node)
	// load config
	cfg, err := loadOrdererConfig(cfgPath)
	if err != nil {
		return err
	}
	// change relative path into absolute path
	cfg.BlockChain.Path = getBFTOrdererPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Path = getBFTOrdererPath(node) + "/" + cfg.DB.LevelDB.Path
	cfg.Consensus.Tendermint.Path = getBFTOrdererPath(node) + "/" + cfg.Consensus.Tendermint.Path
	cfg.TLS.CA = getBFTOrdererPath(node) + "/" + cfg.TLS.CA
	cfg.TLS.RawCert = getBFTOrdererPath(node) + "/" + cfg.TLS.RawCert
	cfg.TLS.Key = getBFTOrdererPath(node) + "/" + cfg.TLS.Key
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

func absBFTClientConfig(node int) error {
	cfgPath := getBFTClientConfigPath(node)
	// load config
	cfg, err := loadClientConfig(cfgPath)
	if err != nil {
		return err
	}
	// change relative path into absolute path
	cfg.KeyStore.Keys[0] = getBFTClientPath(node) + "/" + cfg.KeyStore.Keys[0]
	cfg.TLS.CA = getBFTClientPath(node) + "/" + cfg.TLS.CA
	cfg.TLS.RawCert = getBFTClientPath(node) + "/" + cfg.TLS.RawCert
	cfg.TLS.Key = getBFTClientPath(node) + "/" + cfg.TLS.Key
	// rewrite peer config
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfgPath, data, os.ModePerm)
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

func getOrderersPid() []string {
	cmd := exec.Command("/bin/sh", "-c", "pidof orderer")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	pids := strings.Split(string(output), " ")
	return pids
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

func compareChannels() error {
	infos1, err := bftClients[0].ListChannel(true)
	if err != nil {
		return err
	}
	infos2, err := bftClients[1].ListChannel(true)
	if err != nil {
		return err
	}
	if len(infos1) != len(infos2) {
		return fmt.Errorf("the count of channels is not consistent")
	}
	for i := range infos1 {
		if infos1[i].Name != infos2[i].Name {
			return fmt.Errorf("the name is not consistent")
		}
		if infos1[i].BlockSize != infos2[i].BlockSize {
			return fmt.Errorf("the blockSize is not consistent, %d in client, %d in admin", infos1[i].BlockSize, infos2[i].BlockSize)
		}
	}

	fmt.Println("CompareChannels: channels between two orderers are consistent.")
	return nil
}

func compareTxs() error {
	// get tx history from peer0
	address1, err := bftClients[0].GetPrivKey().PubKey().Address()
	if err != nil {
		return err
	}
	history1, err := bftClients[0].GetHistory(address1.Bytes())
	if err != nil {
		return err
	}
	// get tx history from peer1
	stopPeer(bftPeers[0])
	address2, err := bftClients[0].GetPrivKey().PubKey().Address()
	if err != nil {
		return err
	}

	history2, err := bftClients[0].GetHistory(address2.Bytes())
	if err != nil {
		return err
	}
	bftPeers[0] = startPeer(0)
	/*table1 := cliu.NewTable()
	table1.SetHeader("Channel", "TxID")
	for channel, txs := range history1.Txs {
		for _, id := range txs.Value {
			table1.AddRow(channel, id)
		}
	}
	table1.Render()*/

	if len(history1.Txs) != len(history2.Txs) {
		return fmt.Errorf("the count of txs is not consistent")
	}
	var txs1 = make(map[string]string)
	for channel, txs := range history1.Txs {
		for _, id := range txs.Value {
			txs1[id] = channel
			fmt.Printf("")
		}
	}
	var txs2 = make(map[string]string)
	for channel, txs := range history2.Txs {
		for _, id := range txs.Value {
			txs2[id] = channel
		}
	}

	for key, value := range txs1 {
		if v, ok := txs2[key]; !ok || v != value {
			return fmt.Errorf("the tx not consistent ")
		}
	}

	fmt.Println("CompareTxs: txs between two peers are consistent.")
	return nil
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

// startPeer run peer and return pid
func startPeer(node int) string {
	before := getPeerPid()
	go func() {
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("peer start -c %s", getBFTPeerConfigPath(node)))
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

func createChannelForCallTx() error {
	// client 0 create channel
	err := bftClients[0].CreateChannel("test0", true, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func createContractForCallTx() error {
	// client 0 create contract
	contractCodes, err := readCodes(getBFTClientPath(0) + "/MyTest.bin")
	if err != nil {
		return err
	}
	tx, err := types.NewTx("test0", common.ZeroAddress, contractCodes, bftClients[0].GetPrivKey())
	if err != nil {
		return err
	}
	_, err = bftClients[0].AddTx(tx)
	if err != nil {
		return err
	}

	return nil
}

func getNumForCallTx(node int, num string) error {
	abiPath := fmt.Sprintf(getBFTClientPath(node) + "/MyTest.abi")
	var inputs []string = make([]string, 0)
	payloadBytes, err := abi.GetPayloadBytes(abiPath, "getNum", inputs)
	if err != nil {
		return err
	}

	client := bftClients[node]
	channel := "test" + strconv.Itoa(node)
	tx, err := types.NewTx(channel, common.HexToAddress("0x8de6ce45b289502e16aef93313fd3082993acb1f"), payloadBytes,
		client.GetPrivKey())
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
		return fmt.Errorf("call tx on channel test%d: setNum expect %s but receive %s", node, num, output[0])
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
	tx, err := types.NewTx(channel, common.HexToAddress("0x8de6ce45b289502e16aef93313fd3082993acb1f"), payloadBytes,
		client.GetPrivKey())
	if err != nil {
		return err
	}

	_, err = client.AddTx(tx)
	if err != nil {
		return err
	}
	return nil
}
