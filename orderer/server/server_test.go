package server

import (
	"context"
	"fmt"
	"madledger/common/util"
	"madledger/core/types"
	"madledger/orderer/config"
	pb "madledger/protos"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestNewServer(t *testing.T) {
	server, err := NewServer(getTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		server.Start()
	}()
	time.Sleep(100 * time.Millisecond)
	// then check blocks
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
	server.Stop()
	initTestEnvironment()
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
