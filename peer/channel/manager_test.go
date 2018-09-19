package channel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"madledger/blockchain/config"
	"madledger/common"
	"madledger/common/util"
	"madledger/core/types"
	"madledger/peer/db"
	"madledger/peer/orderer"
	"net"
	"os"
	"testing"
	"time"

	cc "madledger/blockchain/config"
	gc "madledger/blockchain/global"
	pb "madledger/protos"

	"google.golang.org/grpc"
)

// Note: This file could not be used for automotic test

var (
	coordinator      = NewCoordinator()
	leveldb, _       = db.NewLevelDB(".data/leveldb")
	client, _        = orderer.NewClient("localhost:9999")
	globalManager, _ = NewManager(types.GLOBALCHANNELID, ".data/blocks/"+types.GLOBALCHANNELID, leveldb, client, coordinator)
	configManager, _ = NewManager(types.CONFIGCHANNELID, ".data/blocks/"+types.CONFIGCHANNELID, leveldb, client, coordinator)
	testManager, _   = NewManager("test", ".data/blocks/test", leveldb, client, coordinator)
	globalBlocks     = make(map[int]*types.Block)
	configBlocks     = make(map[int]*types.Block)
	testBlocks       = make(map[int]*types.Block)
	globalBlocksEnd  = make(chan bool, 1)
	configBlocksEnd  = make(chan bool, 1)
	testBlocksEnd    = make(chan bool, 1)
)

func TestAll(t *testing.T) {
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
	os.RemoveAll(".data")
}

type fakeOrderer struct {
	rpcServer *grpc.Server
}

func startFakeOrderer() error {
	orderer := new(fakeOrderer)
	addr := fmt.Sprintf("localhost:9999")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.New("Failed to start the orderer")
	}
	var opts []grpc.ServerOption
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
	case types.GLOBALCHANNELID:
		return getGlobalBlock(req.Number), nil
	case types.CONFIGCHANNELID:
		return getConfigBlock(req.Number), nil
	default:
		return getTestBlock(req.Number), nil
	}
}

// ListChannels is the implementation of protos
func (o *fakeOrderer) ListChannels(ctx context.Context, req *pb.ListChannelsRequest) (*pb.ChannelInfos, error) {
	return nil, nil
}

// AddChannel is the implementation of protos
func (o *fakeOrderer) AddChannel(ctx context.Context, req *pb.AddChannelRequest) (*pb.ChannelInfo, error) {
	return nil, nil
}

// AddTx is the implementation of protos
func (o *fakeOrderer) AddTx(ctx context.Context, req *pb.AddTxRequest) (*pb.TxStatus, error) {
	return nil, nil
}

func generateBlocks() {
	// first generate 100 test blocks
	testGenesisBlock := types.NewBlock("test", 0, types.GenesisBlockPrevHash, nil)
	testBlocks[0] = testGenesisBlock
	for i := 1; i < 100; i++ {
		testBlock := types.NewBlock("test", uint64(i), testBlocks[i-1].Hash().Bytes(), nil)
		testBlocks[i] = testBlock
	}
	// then generate 2 config blocks
	// first is genesis config block
	var payloads = []config.Payload{config.Payload{
		ChannelID: types.CONFIGCHANNELID,
		Profile: &config.Profile{
			Public: true,
		},
		Version: 1,
	}, config.Payload{
		ChannelID: types.GLOBALCHANNELID,
		Profile: &config.Profile{
			Public: true,
		},
		Version: 1,
	}}
	var txs []*types.Tx
	for i, payload := range payloads {
		payloadBytes, _ := json.Marshal(&payload)
		accountNonce := uint64(i)
		tx := types.NewTxWithoutSig(types.CONFIGCHANNELID, payloadBytes, accountNonce)
		txs = append(txs, tx)
	}
	genesisConfigBlock := types.NewBlock(types.CONFIGCHANNELID, 0, types.GenesisBlockPrevHash, txs)
	configBlocks[0] = genesisConfigBlock
	// then second config block
	payloadBytes, _ := json.Marshal(cc.Payload{
		ChannelID: "test",
		Profile: &cc.Profile{
			Public: true,
		},
		Version: 1,
	})
	// create tx
	var tx = &types.Tx{
		Data: types.TxData{
			ChannelID:    types.CONFIGCHANNELID,
			AccountNonce: 0,
			Recipient:    common.ZeroAddress.Bytes(),
			Payload:      payloadBytes,
			Version:      1,
			Sig:          nil,
		},
		Time: util.Now(),
	}
	tx.ID = util.Hex(tx.Hash())
	configBlock := types.NewBlock(types.CONFIGCHANNELID, 1, genesisConfigBlock.Hash().Bytes(), []*types.Tx{tx})
	configBlocks[1] = configBlock
	// then genesis global blocks
	ggb, _ := gc.CreateGenesisBlock([]*gc.Payload{&gc.Payload{
		ChannelID: types.CONFIGCHANNELID,
		Number:    0,
		Hash:      genesisConfigBlock.Hash(),
	}})
	globalBlocks[0] = ggb
	payloadBytes, _ = json.Marshal(&gc.Payload{
		ChannelID: types.CONFIGCHANNELID,
		Number:    1,
		Hash:      configBlocks[1].Hash(),
	})
	tx = types.NewTxWithoutSig(types.GLOBALCHANNELID, payloadBytes, 0)
	globalBlocks[1] = types.NewBlock(types.GLOBALCHANNELID, 1, globalBlocks[0].Hash().Bytes(), []*types.Tx{tx})

	// then many blocks, global block begin from num 2
	for i := 0; i < 100; i++ {
		tx := types.NewGlobalTx("test", uint64(i), testBlocks[i].Hash())
		globalBlock := types.NewBlock(types.GLOBALCHANNELID, uint64(i+2), globalBlocks[i+1].Hash().Bytes(), []*types.Tx{tx})
		globalBlocks[i+2] = globalBlock
	}
}

func getGlobalBlock(n uint64) *pb.Block {
	num := int(n)
	if util.Contain(globalBlocks, num) {
		time.Sleep(10 * time.Millisecond)
		block, _ := pb.NewBlock(globalBlocks[num])
		return block
	}
	globalBlocksEnd <- true
	var c = make(chan bool)
	<-c
	return nil
}

func getConfigBlock(n uint64) *pb.Block {
	num := int(n)
	if util.Contain(configBlocks, num) {
		block, _ := pb.NewBlock(configBlocks[num])
		return block
	}
	configBlocksEnd <- true
	var c = make(chan bool)
	<-c
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
