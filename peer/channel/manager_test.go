// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package channel

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"madledger/common"
	"madledger/common/util"
	"madledger/core"
	oc "madledger/orderer/config"
	pc "madledger/peer/config"
	"madledger/peer/db"
	"madledger/peer/orderer"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc/credentials"
	yaml "gopkg.in/yaml.v2"

	"github.com/stretchr/testify/require"

	bc "madledger/blockchain/config"
	gc "madledger/blockchain/global"
	pb "madledger/protos"

	"github.com/otiai10/copy"

	"google.golang.org/grpc"
)

// Note: Run this file and output logs to log.txt
// Then analyse all logs to make sure all orderers all right.

var (
	coordinator      = NewCoordinator()
	leveldb, _       = db.NewLevelDB(".data/leveldb")
	cfg, _           = getPeerConfig()
	client, _        = orderer.NewClient("localhost:9999", cfg)
	globalManager, _ = NewManager(core.GLOBALCHANNELID, ".data/blocks/"+core.GLOBALCHANNELID, nil, leveldb, []*orderer.Client{client}, coordinator)
	configManager, _ = NewManager(core.CONFIGCHANNELID, ".data/blocks/"+core.CONFIGCHANNELID, nil, leveldb, []*orderer.Client{client}, coordinator)
	testManager, _   = NewManager("test", ".data/blocks/test", nil, leveldb, []*orderer.Client{client}, coordinator)
	globalBlocks     = make(map[int]*core.Block)
	configBlocks     = make(map[int]*core.Block)
	testBlocks       = make(map[int]*core.Block)
	globalBlocksEnd  = make(chan bool, 1)
	configBlocksEnd  = make(chan bool, 1)
	testBlocksEnd    = make(chan bool, 1)
	blockSize        = 100
)

var (
	step int
)

func TestManageRun(t *testing.T) {
	// Note: we set step to 1 or 2 to control things we want, default we will do nothing
	flag.IntVar(&step, "s", 0, "step")
	flag.Parse()
	testRun(t)
	testAnalyse(t)
	os.RemoveAll(".data")
}

func testRun(t *testing.T) {
	if step != 1 {
		return
	}

	generateBlocks()
	go func() {
		err := startFakeOrderer()
		if err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(1000 * time.Millisecond)
	go globalManager.Start()
	go configManager.Start()
	go testManager.Start()
	<-globalBlocksEnd
	<-configBlocksEnd
	<-testBlocksEnd
}

// block is used for block analyse
type block struct {
	channel string
	number  int
}

func testAnalyse(t *testing.T) {
	if step != 2 {
		return
	}

	f, err := os.Open(getLogFile())
	require.NoError(t, err)

	var blocks []*block
	buf := bufio.NewReader(f)
	// read log file
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
		}
		line = strings.TrimSpace(line)
		block := findBlock(line)
		if block != nil {
			blocks = append(blocks, block)
		}
	}
	// analyse
	var globalNumber = -100
	var testNumber = -100 // set it a small number
	for _, block := range blocks {
		switch block.channel {
		case "global":
			globalNumber = block.number
		default:
			testNumber = block.number
		}
		if testNumber != 0 && testNumber > globalNumber-2 {
			t.Fatal(fmt.Sprintf("Run test block %d too early because global block is still %d", testNumber, globalNumber))
		}
	}

	require.Equal(t, blockSize+1, globalNumber)
	require.Equal(t, blockSize-1, testNumber)

	os.RemoveAll(".data")
	os.Remove(getLogFile())
}

type fakeOrderer struct {
	rpcServer *grpc.Server
}

func startFakeOrderer() error {
	orderer := new(fakeOrderer)
	addr := fmt.Sprintf("localhost:9999")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("Failed to start the orderer, because %s", err.Error())
	}
	fmt.Printf("Start the orderer at %s", addr)
	var opts []grpc.ServerOption
	cfg, err := getOrdererConfig()
	if cfg.TLS.Enable {
		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{*(cfg.TLS.Cert)},
			ClientCAs:    cfg.TLS.Pool,
		})
		opts = append(opts, grpc.Creds(creds))
	}
	orderer.rpcServer = grpc.NewServer(opts...)
	pb.RegisterOrdererServer(orderer.rpcServer, orderer)
	err = orderer.rpcServer.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

// FetchBlock is the implementation of protos
func (o *fakeOrderer) FetchBlock(ctx context.Context, req *pb.FetchBlockRequest) (*pb.Block, error) {
	switch req.ChannelID {
	case core.GLOBALCHANNELID:
		return getGlobalBlock(req.Number), nil
	case core.CONFIGCHANNELID:
		return getConfigBlock(req.Number), nil
	default:
		return getTestBlock(req.Number), nil
	}
}

// ListChannels is the implementation of protos
func (o *fakeOrderer) ListChannels(ctx context.Context, req *pb.ListChannelsRequest) (*pb.ChannelInfos, error) {
	return nil, nil
}

// CreateChannel is the implementation of protos
func (o *fakeOrderer) CreateChannel(ctx context.Context, req *pb.CreateChannelRequest) (*pb.ChannelInfo, error) {
	return nil, nil
}

// AddTx is the implementation of protos
func (o *fakeOrderer) AddTx(ctx context.Context, req *pb.AddTxRequest) (*pb.TxStatus, error) {
	return nil, nil
}

func (o *fakeOrderer) GetAccountInfo(ctx context.Context, req *pb.GetAccountInfoRequest) (*pb.AccountInfo, error) {
	return nil, nil
}

func (o *fakeOrderer) GetTxStatus(ctx context.Context, req *pb.GetTxStatusRequest) (*pb.TxStatus, error) {
	return nil, nil
}

func generateBlocks() {
	// first generate test blocks
	testGenesisBlock := core.NewBlock("test", 0, core.GenesisBlockPrevHash, nil)
	testBlocks[0] = testGenesisBlock
	for i := 1; i < blockSize; i++ {
		testBlock := core.NewBlock("test", uint64(i), testBlocks[i-1].Hash().Bytes(), nil)
		testBlocks[i] = testBlock
	}
	// then generate 2 config blocks
	// first is genesis config block
	admins, _ := bc.CreateAdmins()
	var payloads = []bc.Payload{bc.Payload{
		ChannelID: core.CONFIGCHANNELID,
		Profile: &bc.Profile{
			Public: true,
		},
		Version: 1,
	}, bc.Payload{
		ChannelID: core.GLOBALCHANNELID,
		Profile: &bc.Profile{
			Public: true,
		},
		Version: 1,
	}, bc.Payload{ // this payload is used to record the info of  system admin
		Profile: &bc.Profile{
			Public: true,
			Admins: admins,
		},
		Version: 1,
	}}
	var txs []*core.Tx
	for i, payload := range payloads {
		payloadBytes, _ := json.Marshal(&payload)
		accountNonce := uint64(i)
		tx := core.NewTxWithoutSig(core.CONFIGCHANNELID, payloadBytes, accountNonce)
		txs = append(txs, tx)
	}
	genesisConfigBlock := core.NewBlock(core.CONFIGCHANNELID, 0, core.GenesisBlockPrevHash, txs)
	configBlocks[0] = genesisConfigBlock
	// then second config block
	payloadBytes, _ := json.Marshal(bc.Payload{
		ChannelID: "test",
		Profile: &bc.Profile{
			Public: true,
		},
		Version: 1,
	})
	// create tx
	var tx = &core.Tx{
		Data: core.TxData{
			ChannelID: core.CONFIGCHANNELID,
			Nonce:     0,
			Recipient: common.ZeroAddress.Bytes(),
			Payload:   payloadBytes,
			Version:   1,
		},
		Time: util.Now(),
	}
	tx.ID = util.Hex(tx.Hash())
	configBlock := core.NewBlock(core.CONFIGCHANNELID, 1, genesisConfigBlock.Hash().Bytes(), []*core.Tx{tx})
	configBlocks[1] = configBlock
	// then genesis global blocks
	ggb, _ := gc.CreateGenesisBlock([]*core.GlobalTxPayload{&core.GlobalTxPayload{
		ChannelID: core.CONFIGCHANNELID,
		Number:    0,
		Hash:      genesisConfigBlock.Hash(),
	}})
	globalBlocks[0] = ggb
	payloadBytes, _ = json.Marshal(&core.GlobalTxPayload{
		ChannelID: core.CONFIGCHANNELID,
		Number:    1,
		Hash:      configBlocks[1].Hash(),
	})
	tx = core.NewTxWithoutSig(core.GLOBALCHANNELID, payloadBytes, 0)
	globalBlocks[1] = core.NewBlock(core.GLOBALCHANNELID, 1, globalBlocks[0].Hash().Bytes(), []*core.Tx{tx})

	// then many blocks, global block begin from num 2
	for i := 0; i < blockSize; i++ {
		tx := core.NewGlobalTx("test", uint64(i), testBlocks[i].Hash())
		globalBlock := core.NewBlock(core.GLOBALCHANNELID, uint64(i+2), globalBlocks[i+1].Hash().Bytes(), []*core.Tx{tx})
		globalBlocks[i+2] = globalBlock
	}
}

func getGlobalBlock(n uint64) *pb.Block {
	num := int(n)
	if util.Contain(globalBlocks, num) {
		randomSleep()
		block, _ := pb.NewBlock(globalBlocks[num])
		return block
	}
	globalBlocksEnd <- true
	return nil
}

func getConfigBlock(n uint64) *pb.Block {
	num := int(n)
	if util.Contain(configBlocks, num) {
		block, _ := pb.NewBlock(configBlocks[num])
		return block
	}
	configBlocksEnd <- true
	return nil
}

func getTestBlock(n uint64) *pb.Block {
	num := int(n)
	if util.Contain(testBlocks, num) {
		block, _ := pb.NewBlock(testBlocks[num])
		return block
	}
	testBlocksEnd <- true
	var c = make(chan bool)
	<-c
	return nil
}

func getLogFile() string {
	gopath := os.Getenv("GOPATH")
	logFile, _ := util.MakeFileAbs("src/madledger/peer/channel/log.txt", gopath)
	return logFile
}

func findBlock(line string) *block {
	if strings.Contains(line, "Add global block") {
		blockRegexp := regexp.MustCompile(`^.+?Add global block ([\d]+).+`)
		params := blockRegexp.FindStringSubmatch(line)
		if len(params) >= 1 {
			num, _ := strconv.Atoi(params[1])
			return &block{
				channel: "global",
				number:  num,
			}
		}
	} else if strings.Contains(line, "Run block test") {
		blockRegexp := regexp.MustCompile(`^.+?Run block test:.?([\d]+).+`)
		params := blockRegexp.FindStringSubmatch(line)
		if len(params) >= 1 {
			num, _ := strconv.Atoi(params[1])
			return &block{
				channel: "test",
				number:  num,
			}
		}
	}
	return nil
}

func randomSleep() {
	time.Sleep(time.Duration(util.RandNum(50)) * time.Millisecond)
}

func absPeerConfig(cfgPath string) error {
	// load config
	cfg, err := loadPeerConfig(cfgPath)
	if err != nil {
		return err
	}
	// change relative path into absolute path
	cfg.BlockChain.Path = getPeerPath() + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Dir = getPeerPath() + "/" + cfg.DB.LevelDB.Dir
	cfg.KeyStore.Key = getPeerPath() + "/" + cfg.KeyStore.Key
	cfg.TLS.CA = getPeerPath() + "/" + cfg.TLS.CA
	cfg.TLS.RawCert = getPeerPath() + "/" + cfg.TLS.RawCert
	cfg.TLS.Key = getPeerPath() + "/" + cfg.TLS.Key
	// rewrite peer config
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfgPath, data, os.ModePerm)
}

func getPeerPath() string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/env/raft/peers/0", gopath)
}

func getPeerConfigPath() string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/env/raft/peers/0/peer.yaml", gopath)
}

func getOrdererPath() string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/env/raft/orderers/0", gopath)
}

func getOrdererConfigPath() string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/env/raft/orderers/0/orderer.yaml", gopath)
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

func getPeerConfig() (*pc.Config, error) {
	// add absolute path
	gopath := os.Getenv("GOPATH")
	if err := copy.Copy(gopath+"/src/madledger/env/raft/.peers/0", gopath+"/src/madledger/env/raft/peers/0"); err != nil {
		return nil, err
	}
	err := absPeerConfig(getPeerConfigPath())
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	cfg, err := pc.LoadConfig(getPeerConfigPath())
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func getOrdererConfig() (*pc.Config, error) {
	// add absolute path
	gopath := os.Getenv("GOPATH")
	if err := copy.Copy(gopath+"/src/madledger/env/raft/.orderers/0", gopath+"/src/madledger/env/raft/orderers/0"); err != nil {
		return nil, err
	}
	err := absOrdererConfig(getOrdererConfigPath())
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	cfg, err := pc.LoadConfig(getOrdererConfigPath())
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func absOrdererConfig(cfgPath string) error {
	// load config
	cfg, err := loadOrdererConfig(cfgPath)
	if err != nil {
		return err
	}
	// change relative path into absolute path
	cfg.BlockChain.Path = getOrdererPath() + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Path = getOrdererPath() + "/" + cfg.DB.LevelDB.Path
	cfg.Consensus.Tendermint.Path = getOrdererPath() + "/" + cfg.Consensus.Tendermint.Path
	cfg.TLS.CA = getOrdererPath() + "/" + cfg.TLS.CA
	cfg.TLS.RawCert = getOrdererPath() + "/" + cfg.TLS.RawCert
	cfg.TLS.Key = getOrdererPath() + "/" + cfg.TLS.Key
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
