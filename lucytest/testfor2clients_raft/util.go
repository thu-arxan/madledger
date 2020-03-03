package testfor2clients_raft

import (
	"encoding/hex"
	"evm/abi"
	"fmt"
	"io"
	"io/ioutil"
	cc "madledger/client/config"
	client "madledger/client/lib"
	"madledger/common"
	"madledger/common/util"
	"madledger/core"
	oc "madledger/orderer/config"
	pc "madledger/peer/config"
	"os"
	"os/exec"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/otiai10/copy"
)

var (
	// we need 4 orderers to test 1 orderer breaking down
	raftOrderers [4]string
	// just 1 is enough, we set 2
	raftClients [2]*client.Client
	raftPeers   [4]string
)

// initBFTEnvironment will remove old test folders and copy necessary folders
func initRAFTEnvironment() error {
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
	if err := copy.Copy(gopath+"/src/madledger/env/raft/.orderers", gopath+"/src/madledger/tests/raft/orderers"); err != nil {
		return err
	}
	if err := copy.Copy(gopath+"/src/madledger/env/raft/.clients", gopath+"/src/madledger/tests/raft/clients"); err != nil {
		return err
	}
	if err := copy.Copy(gopath+"/src/madledger/env/raft/.peers", gopath+"/src/madledger/tests/raft/peers"); err != nil {
		return err
	}

	for i := range raftOrderers {
		if err := absRAFTOrdererConfig(i); err != nil {
			return err
		}
	}
	for i := range raftPeers {
		if err := absRAFTPeerConfig(i); err != nil {
			return err
		}
	}
	for i := range raftClients {
		if err := absRAFTClientConfig(i); err != nil {
			return err
		}
	}

	return nil
}

func absRAFTClientConfig(node int) error {
	cfgPath := getRAFTClientConfigPath(node)
	// load config
	cfg, err := loadClientConfig(cfgPath)
	if err != nil {
		return err
	}
	// change relative path into absolute path
	cfg.KeyStore.Keys[0] = getRAFTClientPath(node) + "/" + cfg.KeyStore.Keys[0]
	cfg.TLS.CA = getRAFTClientPath(node) + "/" + cfg.TLS.CA
	cfg.TLS.RawCert = getRAFTClientPath(node) + "/" + cfg.TLS.RawCert
	cfg.TLS.Key = getRAFTClientPath(node) + "/" + cfg.TLS.Key
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

func absRAFTPeerConfig(node int) error {
	cfgPath := getRAFTPeerConfigPath(node)
	// load config
	cfg, err := loadPeerConfig(cfgPath)
	if err != nil {
		return err
	}
	// change relative path into absolute path
	cfg.BlockChain.Path = getRAFTPeerPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Dir = getRAFTPeerPath(node) + "/" + cfg.DB.LevelDB.Dir
	cfg.KeyStore.Key = getRAFTPeerPath(node) + "/" + cfg.KeyStore.Key
	cfg.TLS.CA = getRAFTPeerPath(node) + "/" + cfg.TLS.CA
	cfg.TLS.RawCert = getRAFTPeerPath(node) + "/" + cfg.TLS.RawCert
	cfg.TLS.Key = getRAFTPeerPath(node) + "/" + cfg.TLS.Key
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
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("peer start -c %s", getRAFTPeerConfigPath(node)))
		_, err := cmd.Output()
		if err != nil && !strings.Contains(err.Error(), "exit status 143") {
			panic(fmt.Sprintf("Run peer failed: %s", err))
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

func absRAFTOrdererConfig(node int) error {
	cfgPath := getRAFTOrdererConfigPath(node)
	// load config
	cfg, err := loadOrdererConfig(cfgPath)
	if err != nil {
		return err
	}
	// change relative path into absolute path
	cfg.BlockChain.Path = getRAFTOrdererPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Path = getRAFTOrdererPath(node) + "/" + cfg.DB.LevelDB.Path
	cfg.Consensus.Raft.Path = getRAFTOrdererPath(node) + "/" + cfg.Consensus.Raft.Path
	cfg.TLS.CA = getRAFTOrdererPath(node) + "/" + cfg.TLS.CA
	cfg.TLS.RawCert = getRAFTOrdererPath(node) + "/" + cfg.TLS.RawCert
	cfg.TLS.Key = getRAFTOrdererPath(node) + "/" + cfg.TLS.Key
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

func getRAFTOrdererConfigPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/raft/orderers/%d/orderer.yaml", gopath, node)
}

func getRAFTOrdererPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/raft/orderers/%d", gopath, node)
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

func getPeerConfig(node int) *pc.Config {
	cfgFilePath := getRAFTPeerConfigPath(node)
	cfg, _ := pc.LoadConfig(cfgFilePath)

	cfg.BlockChain.Path = getRAFTPeerPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Dir = getRAFTPeerPath(node) + "/" + cfg.DB.LevelDB.Dir

	// then set key
	cfg.KeyStore.Key = getRAFTPeerPath(node) + "/" + cfg.KeyStore.Key
	return cfg
}

func getRAFTPeerConfigPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/raft/peers/%d/peer.yaml", gopath, node)
}

func copyFile(src string, dst string) (writelne int64, err error) {
	srcFile, err := os.Open(src)

	if err != nil {
		fmt.Printf("open file err = %v\n", err)
		return
	}

	defer srcFile.Close()

	//open destination file
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Printf("open file err = %v\n", err)
		return
	}

	defer dstFile.Close()
	return io.Copy(dstFile, srcFile)
}

func readCodes(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(string(data))
}

func getRAFTPeerPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/raft/peers/%d", gopath, node)
}

func getRAFTPeerDataPath(node int) string {
	return fmt.Sprintf("%s/data", getRAFTPeerPath(node))
}

func getRAFTClientPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/raft/clients/%d", gopath, node)
}

func getRAFTClientConfigPath(node int) string {
	return getRAFTClientPath(node) + "/client.yaml"
}

func getRAFTOrdererDataPath(node int) string {
	return fmt.Sprintf("%s/data", getRAFTOrdererPath(node))
}

func getRAFTOrdererWALPath(node int) string {
	return fmt.Sprintf("%s/.raft/wal", getRAFTOrdererPath(node))
}

func getRAFTOrdererSnapPath(node int) string {
	return fmt.Sprintf("%s/.raft/snap", getRAFTOrdererPath(node))
}

func getRAFTOrdererDBPath(node int) string {
	return fmt.Sprintf("%s/.raft/db", getRAFTOrdererPath(node))
}

// stopOrderer stop an orderer
func stopOrderer(pid string) {
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("kill -TERM %s", pid))
	cmd.Output()
}

func compareChannelName(channels []string) error {
	lenChannels := len(channels) + 2
	for i := range raftClients {
		client := raftClients[i]
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

func compareClientTx(len int, channelName string) error {
	client := raftClients[0]
	address, err := client.GetPrivKey().PubKey().Address()
	if err != nil {
		return err
	}

	history, err := client.GetHistory(address.Bytes())
	if err != nil {
		return err
	}

	count := 0
	for channel, txs := range history.Txs {
		if channel == channelName {
			for _, id := range txs.Value {
				fmt.Printf("(%s, %s)\n", channel, id)
				count++
			}
			if count != len {
				return fmt.Errorf("channel %s should have %d txs but there is %d", channel, len, count)
			}
		}
	}

	return nil
}

func compareChannels() error {
	infos1, err := raftClients[0].ListChannel(true)
	if err != nil {
		return err
	}
	infos2, err := raftClients[1].ListChannel(true)
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

func backupMdFile1(path string) error {
	_, err := copyFile("./0.md", path+"0.md")
	if err != nil {
		return err
	}
	err = os.Remove("./0.md")
	if err != nil {
		return err
	}

	_, err = copyFile("./00.md", path+"00.md")
	if err != nil {
		return err
	}
	err = os.Remove("./00.md")
	if err != nil {
		return err
	}

	_, err = copyFile("./1.md", path+"1.md")
	if err != nil {
		return err
	}
	err = os.Remove("./1.md")
	if err != nil {
		return err
	}

	_, err = copyFile("./000.md", path+"000.md")
	if err != nil {
		return err
	}
	err = os.Remove("./000.md")
	if err != nil {
		return err
	}

	_, err = copyFile("./2.md", path+"2.md")
	if err != nil {
		return err
	}
	err = os.Remove("./2.md")
	if err != nil {
		return err
	}

	_, err = copyFile("./0000.md", path+"0000.md")
	if err != nil {
		return err
	}
	err = os.Remove("./0000.md")
	if err != nil {
		return err
	}

	_, err = copyFile("./3.md", path+"3.md")
	if err != nil {
		return err
	}
	err = os.Remove("./3.md")
	if err != nil {
		return err
	}

	return nil
}

func backupMdFile2(path string) error {
	_, err := copyFile("./0.md", path+"0.md")
	if err != nil {
		return err
	}
	err = os.Remove("./0.md")
	if err != nil {
		return err
	}

	_, err = copyFile("./1.md", path+"1.md")
	if err != nil {
		return err
	}
	err = os.Remove("./1.md")
	if err != nil {
		return err
	}

	_, err = copyFile("./2.md", path+"2.md")
	if err != nil {
		return err
	}
	err = os.Remove("./2.md")
	if err != nil {
		return err
	}

	_, err = copyFile("./3.md", path+"3.md")
	if err != nil {
		return err
	}
	err = os.Remove("./3.md")
	if err != nil {
		return err
	}

	return nil
}

func setNumForCallTx(node int, num string) error {
	abiPath := fmt.Sprintf(getRAFTClientPath(node) + "/MyTest.abi")
	inputs := []string{num}
	payloadBytes, err := abi.Pack(abiPath, "setNum", inputs...)
	if err != nil {
		return err
	}

	client := raftClients[0]
	if node == 1 {
		client = raftClients[1]
	}
	addr := "0x8de6ce45b289502e16aef93313fd3082993acb1f"
	if node == 1 {
		addr = "0x1b66001e01d3c8d3893187fee59e3bea1d9bdd9b"
	}
	channel := "test0"
	if node == 1 {
		channel = "test1"
	}
	tx, err := core.NewTx(channel, common.HexToAddress(addr), payloadBytes, 0, "", client.GetPrivKey())
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
	abiPath := fmt.Sprintf(getRAFTClientPath(node) + "/MyTest.abi")
	var inputs = make([]string, 0)
	payloadBytes, err := abi.Pack(abiPath, "getNum", inputs...)
	if err != nil {
		return err
	}

	client := raftClients[0]
	if node == 1 {
		client = raftClients[1]
	}
	addr := "0x8de6ce45b289502e16aef93313fd3082993acb1f"
	if node == 1 {
		addr = "0x1b66001e01d3c8d3893187fee59e3bea1d9bdd9b"
	}
	channel := "test0"
	if node == 1 {
		channel = "test1"
	}
	tx, err := core.NewTx(channel, common.HexToAddress(addr), payloadBytes, 0, "", client.GetPrivKey())
	if err != nil {
		return err
	}

	status, err := client.AddTx(tx)
	if err != nil {
		return err
	}

	output, err := abi.Unpack(abiPath, "getNum", status.Output)
	if err != nil {
		return err
	}

	if output[0] != num {
		return fmt.Errorf("call tx on channel test%d: getNum expect %s but receive %s", 0, num, output[0])
	}
	return nil
}

// startOrderer run orderer and return pid
func startOrderer(node int) string {
	before := getOrderersPid()
	go func() {
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("orderer start -c %s", getRAFTOrdererConfigPath(node)))
		//cmd:=exec.Command("gnome-terminal -e 'bash -c \"echo 'hello'; exec bash\"'")
		_, err := cmd.Output()
		if err != nil && !strings.Contains(err.Error(), "exit status 2") {
			panic(fmt.Sprintf("Run orderer failed, because %s", err))
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

func compareTxs() error {
	// get tx history from peer0
	address1, err := raftClients[0].GetPrivKey().PubKey().Address()
	if err != nil {
		return err
	}
	history1, err := raftClients[0].GetHistory(address1.Bytes())
	if err != nil {
		return err
	}
	// get tx history from peer1
	stopPeer(raftPeers[0])
	address2, err := raftClients[0].GetPrivKey().PubKey().Address()
	if err != nil {
		return err
	}

	history2, err := raftClients[0].GetHistory(address2.Bytes())
	if err != nil {
		return err
	}
	raftPeers[0] = startPeer(0)

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
