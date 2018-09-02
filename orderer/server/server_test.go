package server

import (
	"context"
	"fmt"
	"madledger/common"
	"madledger/common/util"
	"madledger/core/types"
	"madledger/orderer/config"
	pb "madledger/protos"
	"os"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc"
)

var (
	server            *Server
	genesisBlocksHash = make(map[string]common.Hash)
)

func TestNewServer(t *testing.T) {
	var err error
	server, err = NewServer(getTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		server.Start()
	}()
	time.Sleep(100 * time.Millisecond)
}

func TestListChannelsAtNil(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal()
	}
	infos, err := client.ListChannels(context.Background(), &pb.ListChannelsRequest{
		System: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(infos.Channels) != 2 {
		t.Fatal()
	}
	for _, channel := range infos.Channels {
		switch channel.ChannelID {
		case types.GLOBALCHANNELID:
			if channel.BlockSize != 1 {
				t.Fatal(channel.BlockSize)
			}
		case types.CONFIGCHANNELID:
			if channel.BlockSize != 1 {
				t.Fatal(channel.BlockSize)
			}
		default:
			t.Fatal(fmt.Errorf("Unknown channel %s", channel.ChannelID))
		}
	}
	infos, err = client.ListChannels(context.Background(), &pb.ListChannelsRequest{
		System: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(infos.Channels) != 0 {
		t.Fatal()
	}
}

func TestFetchBlockAtNil(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal()
	}
	globalGenesisBlock, err := client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: types.GLOBALCHANNELID,
		Number:    0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if globalGenesisBlock.Header.Number != 0 {
		t.Fatal(fmt.Errorf("The number of block is %d", globalGenesisBlock.Header.Number))
	}
	// set global genesis block hash
	typesGlobalGenesisBlock, err := ConvertBlockFromPbToTypes(globalGenesisBlock)
	if err != nil {
		t.Fatal(err)
	}
	genesisBlocksHash[types.GLOBALCHANNELID] = typesGlobalGenesisBlock.Hash()
	// test bigger number
	_, err = client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: types.GLOBALCHANNELID,
		Number:    1,
	})
	if err == nil {
		t.Fatal()
	}
	// test empty channel id
	_, err = client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: "",
		Number:    1,
	})
	if err == nil {
		t.Fatal()
	}
	// test a channel which is not exist
	_, err = client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: "test",
		Number:    1,
	})
	// get genesis block of config
	configGenesisBlock, err := client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: types.CONFIGCHANNELID,
		Number:    0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if configGenesisBlock.Header.Number != 0 {
		t.Fatal(fmt.Errorf("The number of block is %d", configGenesisBlock.Header.Number))
	}
	// set config genesis block hash
	typesConfigGenesisBlock, err := ConvertBlockFromPbToTypes(configGenesisBlock)
	if err != nil {
		t.Fatal(err)
	}
	genesisBlocksHash[types.CONFIGCHANNELID] = typesConfigGenesisBlock.Hash()
	server.Stop()
	initTestEnvironment()
}

func TestServerStartAtAnotherPath(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	cfg, _ := config.LoadConfig(getTestConfigFilePath())
	chainPath, _ := util.MakeFileAbs("src/madledger/orderer/server/.data1/blocks", gopath)
	dbPath, _ := util.MakeFileAbs("src/madledger/orderer/server/.data1/leveldb", gopath)
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Dir = dbPath
	var err error
	server, err = NewServer(getTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		server.Start()
	}()
	time.Sleep(100 * time.Millisecond)
	client, err := getClient()
	if err != nil {
		t.Fatal()
	}
	// compare global genesis block
	globalGenesisBlock, _ := client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: types.GLOBALCHANNELID,
		Number:    0,
	})
	typesGlobalGenesisBlock, _ := ConvertBlockFromPbToTypes(globalGenesisBlock)
	if !reflect.DeepEqual(typesGlobalGenesisBlock.Hash().Bytes(), genesisBlocksHash[types.GLOBALCHANNELID].Bytes()) {
		t.Fatal()
	}
	// compare config genesis block
	configGenesisBlock, _ := client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: types.CONFIGCHANNELID,
		Number:    0,
	})
	typesConfigGenesisBlock, _ := ConvertBlockFromPbToTypes(configGenesisBlock)
	if !reflect.DeepEqual(typesConfigGenesisBlock.Hash().Bytes(), genesisBlocksHash[types.CONFIGCHANNELID].Bytes()) {
		t.Fatal()
	}
	server.Stop()

}

func getClient() (pb.OrdererClient, error) {
	var conn *grpc.ClientConn
	var err error
	conn, err = grpc.Dial("localhost:12345", grpc.WithInsecure(), grpc.WithTimeout(2000*time.Millisecond))
	if err != nil {
		return nil, err
	}
	client := pb.NewOrdererClient(conn)
	return client, nil
}

func getTestConfig() *config.Config {
	cfg, _ := config.LoadConfig(getTestConfigFilePath())
	cfg.BlockChain.Path = getTestChainPath()
	cfg.DB.LevelDB.Dir = getTestDBPath()
	return cfg
}

func getTestChainPath() string {
	gopath := os.Getenv("GOPATH")
	chainPath, _ := util.MakeFileAbs("src/madledger/orderer/server/.data/blocks", gopath)
	return chainPath
}

func getTestDBPath() string {
	gopath := os.Getenv("GOPATH")
	dbPath, _ := util.MakeFileAbs("src/madledger/orderer/server/.data/leveldb", gopath)
	return dbPath
}

func getTestConfigFilePath() string {
	gopath := os.Getenv("GOPATH")
	cfgFilePath, _ := util.MakeFileAbs("src/madledger/orderer/config/.orderer.yaml", gopath)
	return cfgFilePath
}

func initTestEnvironment() {
	gopath := os.Getenv("GOPATH")
	dataPath, _ := util.MakeFileAbs("src/madledger/orderer/server/.data", gopath)
	os.RemoveAll(dataPath)
}
