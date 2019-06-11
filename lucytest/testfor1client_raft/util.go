package testfor1client_raft

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
	pc "madledger/peer/config"
	peer "madledger/peer/server"
	"madledger/common/util"
	"encoding/hex"
	client "madledger/client/lib"
	"os"
	oc "madledger/orderer/config"
	"os/exec"
	"strings"
	"time"

	"github.com/otiai10/copy"
)

var (
	// raft的orderer只需要3个
	bftOrderers [3]string
	// just 1 is enough, we set 2
	bftClients [2]*client.Client
	bftPeers   [4]*peer.Server
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

	for i := range bftOrderers {
		if err := absRAFTOrdererConfig(i); err != nil {
			return err
		}
	}

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

// startOrderer run orderer and return pid
func startOrderer(node int) string {
	before := getOrderersPid()
	go func() {
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("orderer start -c %s", getRAFTOrdererConfigPath(node)))
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
