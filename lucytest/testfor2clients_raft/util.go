package testfor2clients_raft

import (
	"encoding/hex"
	"fmt"
	"io"
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
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/otiai10/copy"
)

var (
	// we need 4 orderers to test 1 orderer breaking down
	raftOrderers [4]string
	// just 1 is enough, we set 2
	raftClients [2]*client.Client
	raftPeers   [4]*peer.Server
)

// initBFTEnvironment will remove old test folders and copy necessary folders
func initRAFTEnvironment() error {
	// kill all Orderers
	pids := getOrderersPid()
	for _, pid := range pids {
		stopOrderer(pid)
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

	return nil
}

func initPeer(node int) error {
	cfg := getPeerConfig(node)
	server, err := peer.NewServer(cfg)
	if err != nil {
		return err
	}
	raftPeers[node] = server
	return nil
}

func absRAFTOrdererConfig(node int) error {
	cfgPath := getRAFTOrdererConfigPath(node)
	cfg, err := oc.LoadConfig(cfgPath)
	if err != nil {
		return err
	}
	// yaml in raft is absolute path
	cfg.BlockChain.Path = getRAFTOrdererPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Path = getRAFTOrdererPath(node) + "/" + cfg.DB.LevelDB.Path
	cfg.Consensus.Raft.Path = getRAFTOrdererPath(node) + "/" + cfg.Consensus.Raft.Path

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfgPath, data, os.ModePerm)
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

func compareChannelBlocks() error {
	client0 := raftClients[0]
	client1 := raftClients[1]
	infos1, err := client0.ListChannel(true)
	if err != nil {
		return err
	}
	infos2, err := client1.ListChannel(true)
	if err != nil {
		return err
	}

	if len(infos1) != len(infos2) {
		return fmt.Errorf("the channel number is not consistent between orderer0 and orderer1")
	}
	for i := range infos1 {
		if infos1[i].BlockSize != infos2[i].BlockSize {
			return fmt.Errorf("the block size is not consistent between %s "+
				"in orderer0 and %s in orderer1", infos1[i].Name, infos2[i].Name)
		}
		fmt.Printf("%s in orderer0 has %d blocks, %s in orderer1 has %d blocks\n",
			infos1[i].Name, infos1[i].BlockSize, infos2[i].Name, infos2[i].BlockSize)
	}

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

func setNumForCallTx(num string) error {
	abiPath := fmt.Sprintf(getRAFTClientPath(0) + "/MyTest.abi")
	inputs := []string{num}
	payloadBytes, err := abi.GetPayloadBytes(abiPath, "setNum", inputs)
	if err != nil {
		return err
	}

	client := raftClients[0]
	addr := "0x8de6ce45b289502e16aef93313fd3082993acb1f"
	tx, err := types.NewTx("test0", common.HexToAddress(addr), payloadBytes, client.GetPrivKey())
	if err != nil {
		return err
	}

	_, err = client.AddTx(tx)
	if err != nil {
		return err
	}
	return nil
}

func getNumForCallTx(num string) error {
	abiPath := fmt.Sprintf(getRAFTClientPath(0) + "/MyTest.abi")
	var inputs []string = make([]string, 0)
	payloadBytes, err := abi.GetPayloadBytes(abiPath, "getNum", inputs)
	if err != nil {
		return err
	}

	client := raftClients[0]
	addr := "0x8de6ce45b289502e16aef93313fd3082993acb1f"
	tx, err := types.NewTx("test0", common.HexToAddress(addr), payloadBytes, client.GetPrivKey())
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
		return fmt.Errorf("call tx on channel test%d: getNum expect %s but receive %s", 0, num, output[0])
	}
	return nil
}

// startOrderer run orderer and return pid
func startOrderer(node int) string {
	before := getOrderersPid()
	go func() {
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("orderer start -c %s > %d.md",
			getRAFTOrdererConfigPath(node), node))
		//cmd:=exec.Command("gnome-terminal -e 'bash -c \"echo 'hello'; exec bash\"'")
		_, err := cmd.Output()
		if err != nil {
			panic(fmt.Sprintf("Run orderer failed:%s", err))
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
