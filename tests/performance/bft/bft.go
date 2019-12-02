package bft

import (
	"fmt"
	tc "github.com/tendermint/tendermint/config"
	tlc "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
	tt "github.com/tendermint/tendermint/types/time"
	"io/ioutil"
	cutil "madledger/client/util"
	"madledger/common/util"
	oc "madledger/orderer/config"
	orderer "madledger/orderer/server"
	pc "madledger/peer/config"
	peer "madledger/peer/server"
	"os"
	"strings"
	"time"
)

var (
	ordererServers []*orderer.Server
	peerServers    []*peer.Server
)


// Init will init environment
func Init(num int) error {
	Clean()
	os.MkdirAll(getOrdererPath(), os.ModePerm)
	os.MkdirAll(getPeerPath(), os.ModePerm)
	os.MkdirAll(getClientsPath(), os.ModePerm)
	return newClients(num)
}

// Clean clean environment
func Clean() {
	os.RemoveAll(getOrdererPath())
	os.RemoveAll(getPeerPath())
	os.RemoveAll(getClientsPath())
}

func getOrdererPath() string {
	path, _ := util.MakeFileAbs("src/madledger/tests/performance/bft/.orderer", gopath)
	return path
}

func getPeerPath() string {
	path, _ := util.MakeFileAbs("src/madledger/tests/performance/bft/.peer", gopath)
	return path
}

func getClientsPath() string {
	path, _ := util.MakeFileAbs("src/madledger/tests/performance/bft/.clients", gopath)
	return path
}

func newClients(num int) error {
	for i := range clients {
		cp, _ := util.MakeFileAbs(fmt.Sprintf("%d", i), getClientsPath())
		os.MkdirAll(cp, os.ModePerm)
		if err := newClient(cp, num); err != nil {
			return err
		}
	}

	return nil
}

func newClient(path string, num int) error {
	cfgPath, _ := util.MakeFileAbs("client.yaml", path)
	keyStorePath, _ := util.MakeFileAbs(".keystore", path)
	os.MkdirAll(keyStorePath, os.ModePerm)

	keyPath, err := cutil.GeneratePrivateKey(keyStorePath)
	if err != nil {
		return err
	}

	var cfg = clientConfigTemplate
	for i := 1; i <= num; i++ {
		port := 23333 + (i - 1)
		cfg = strings.Replace(cfg, fmt.Sprintf("<<<ADDRESS%d>>>", i), fmt.Sprintf("- localhost:%d", port), 1)
	}
	for i := num + 1; i <= 3; i++ {
		cfg = strings.Replace(cfg, fmt.Sprintf("<<<ADDRESS%d>>>", i), "", 1)
	}
	cfg = strings.Replace(cfg, "<<<KEYFILE>>>", keyPath, 1)

	return ioutil.WriteFile(cfgPath, []byte(cfg), os.ModePerm)
}

// StartOrderers start orderers
func StartOrderers() error {
	for i := 1; i <= 4; i++ {
		cfg, err := getOrdererConfig(i)
		if err != nil {
			return err
		}
		ordererServer, err := orderer.NewServer(cfg)
		if err != nil {
			return err
		}
		ordererServers = append(ordererServers, ordererServer)
		go func() {
			ordererServer.Start()
		}()
		time.Sleep(300 * time.Millisecond)
	}
	return nil
}

func getOrdererConfig(id int) (*oc.Config, error) {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/config/orderer/bft_orderer.yaml", gopath)
	cfg, err := oc.LoadConfig(cfgFilePath)
	if err != nil {
		return nil, err
	}
	chainPath, _ := util.MakeFileAbs("data/blocks", getOrdererPath()+fmt.Sprintf("/%d", id))
	dbPath, _ := util.MakeFileAbs("data/leveldb", getOrdererPath()+fmt.Sprintf("/%d", id))
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Path = dbPath
	cfg.Port = 12345 + (id-1)*11111
	// get tendermintP2PID
	var tendermintP2PID string
	if tendermintP2PID, err = initTendermintEnv(getOrdererPath()+fmt.Sprintf("/%d", id)); err != nil {
		return nil,err
	}
	cfg.Consensus.Tendermint.Path = fmt.Sprintf("%s/%d/%s", getOrdererPath(), id, cfg.Consensus.Tendermint.Path)
	cfg.Consensus.Tendermint.ID = tendermintP2PID

	return cfg, nil
}

// initTendermintEnv will create all necessary things that tendermint needs
func initTendermintEnv(path string) (string, error) {
	tendermintPath, _ := util.MakeFileAbs(".tendermint", path)
	os.MkdirAll(tendermintPath+"/config", 0777)
	os.MkdirAll(tendermintPath+"/data", 0777)
	var conf = tc.DefaultConfig()
	privValKeyFile := tendermintPath + "/" + conf.PrivValidatorKeyFile()
	privValStateFile := tendermintPath + "/" + conf.PrivValidatorStateFile()
	var pv *privval.FilePV
	if tlc.FileExists(privValKeyFile) {
		pv = privval.LoadFilePV(privValKeyFile, privValStateFile)
	} else {
		pv = privval.GenFilePV(privValKeyFile, privValStateFile)
		pv.Save()
	}
	nodeKeyFile := tendermintPath + "/" + conf.NodeKeyFile()
	if !tlc.FileExists(nodeKeyFile) {
		if _, err := p2p.LoadOrGenNodeKey(nodeKeyFile); err != nil {
			return "", err
		}
	}

	// genesis file
	genFile := tendermintPath + "/" + conf.GenesisFile()
	if !tlc.FileExists(genFile) {
		genDoc := types.GenesisDoc{
			ChainID:         "madledger",
			GenesisTime:     tt.Now(),
			ConsensusParams: types.DefaultConsensusParams(),
		}
		genDoc.Validators = []types.GenesisValidator{{
			Address: pv.GetPubKey().Address(),
			PubKey:  pv.GetPubKey(),
			Power:   10,
		}}

		if err := genDoc.SaveAs(genFile); err != nil {
			return "", err
		}
	}

	// load node key
	nodeKey, err := p2p.LoadNodeKey(nodeKeyFile)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", nodeKey.ID()), nil
}



// StartPeers start peers
func StartPeers(num int) error {
	for i := 1; i <= num; i++ {
		peerServer, err := peer.NewServer(getPeerConfig(i))
		if err != nil {
			return err
		}
		peerServers = append(peerServers, peerServer)
		go func() {
			peerServer.Start()
		}()
		time.Sleep(300 * time.Millisecond)
	}
	return nil
}

func getPeerConfig(id int) *pc.Config {
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/tests/config/peer/bft_peer.yaml", gopath)
	cfg, _ := pc.LoadConfig(cfgFilePath)
	chainPath, _ := util.MakeFileAbs("data/blocks", fmt.Sprintf("%s/%d", getPeerPath(), id))
	dbPath, _ := util.MakeFileAbs("data/leveldb", fmt.Sprintf("%s/%d", getPeerPath(), id))
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Dir = dbPath
	key, _ := util.MakeFileAbs(fmt.Sprintf("src/madledger/tests/config/peer/.bft_peer%d.pem", id), gopath)
	cfg.KeyStore.Key = key
	cfg.Port = 23333 + (id - 1)
	return cfg
}