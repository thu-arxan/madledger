package server

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	cc "madledger/blockchain/config"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/core/types"
	"madledger/orderer/config"
	pb "madledger/protos"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
)

const (
	secp256k1String = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
)

var (
	rawSecp256k1Bytes, _ = hex.DecodeString(secp256k1String)
	rawPrivKey           = rawSecp256k1Bytes
	privKey, _           = crypto.NewPrivateKey(rawPrivKey)
	pubKeyBytes, _       = privKey.PubKey().Bytes()
)

var (
	server            *Server
	genesisBlocksHash = make(map[string]common.Hash)
)

func TestNewServer(t *testing.T) {
	initTestEnvironment(".data")
	initTestEnvironment(".data1")
	var err error
	server, err = NewServer(getTestConfig())
	require.NoError(t, err)

	go func() {
		server.Start()
	}()
	time.Sleep(300 * time.Millisecond)
}

func TestListChannelsAtNil(t *testing.T) {
	client, err := getClient()
	require.NoError(t, err)

	infos, err := client.ListChannels(context.Background(), &pb.ListChannelsRequest{
		System: true,
		PK:     pubKeyBytes,
	})
	require.NoError(t, err)
	require.Len(t, infos.Channels, 2)

	for _, channel := range infos.Channels {
		switch channel.ChannelID {
		case types.GLOBALCHANNELID, types.CONFIGCHANNELID:
			require.Equal(t, channel.BlockSize, uint64(1))
		default:
			t.Fatal(fmt.Errorf("Unknown channel %s", channel.ChannelID))
		}
	}
	infos, err = client.ListChannels(context.Background(), &pb.ListChannelsRequest{
		System: false,
		PK:     pubKeyBytes,
	})
	require.NoError(t, err)
	require.Len(t, infos.Channels, 0)
}

func TestFetchBlockAtNil(t *testing.T) {
	client, err := getClient()
	require.NoError(t, err)

	globalGenesisBlock, err := client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: types.GLOBALCHANNELID,
		Number:    0,
	})
	require.NoError(t, err)
	require.Equal(t, globalGenesisBlock.Header.Number, uint64(0))
	// set global genesis block hash
	typesGlobalGenesisBlock, err := globalGenesisBlock.ConvertToTypes()
	require.NoError(t, err)

	genesisBlocksHash[types.GLOBALCHANNELID] = typesGlobalGenesisBlock.Hash()
	// test bigger number
	_, err = client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: types.GLOBALCHANNELID,
		Number:    1,
	})
	require.Error(t, err)
	// test empty channel id
	_, err = client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: "",
		Number:    1,
	})
	require.Error(t, err)
	// test a channel which is not exist
	_, err = client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: "test",
		Number:    1,
	})
	require.Error(t, err)
	// get genesis block of config
	configGenesisBlock, err := client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: types.CONFIGCHANNELID,
		Number:    0,
	})
	require.NoError(t, err)
	require.Equal(t, configGenesisBlock.Header.Number, uint64(0))
	// set config genesis block hash
	typesConfigGenesisBlock, err := configGenesisBlock.ConvertToTypes()
	require.NoError(t, err)

	genesisBlocksHash[types.CONFIGCHANNELID] = typesConfigGenesisBlock.Hash()
	server.Stop()
}

func TestServerRestart(t *testing.T) {
	var err error
	server, err = NewServer(getTestConfig())
	require.NoError(t, err)

	go func() {
		err := server.Start()
		require.NoError(t, err)
	}()
	time.Sleep(200 * time.Millisecond)
	server.Stop()
}

func TestServerStartAtAnotherPath(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	cfg, _ := config.LoadConfig(getTestConfigFilePath())
	chainPath, _ := util.MakeFileAbs("src/madledger/orderer/server/.data1/blocks", gopath)
	dbPath, _ := util.MakeFileAbs("src/madledger/orderer/server/.data1/leveldb", gopath)
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Dir = dbPath
	var err error
	server, err = NewServer(cfg)
	require.NoError(t, err)

	go func() {
		server.Start()
	}()
	time.Sleep(300 * time.Millisecond)
	client, err := getClient()
	require.NoError(t, err)
	// compare global genesis block
	globalGenesisBlock, _ := client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: types.GLOBALCHANNELID,
		Number:    0,
	})
	typesGlobalGenesisBlock, _ := globalGenesisBlock.ConvertToTypes()
	if !reflect.DeepEqual(typesGlobalGenesisBlock.Hash().Bytes(), genesisBlocksHash[types.GLOBALCHANNELID].Bytes()) {
		t.Fatal()
	}
	// compare config genesis block
	configGenesisBlock, _ := client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: types.CONFIGCHANNELID,
		Number:    0,
	})
	typesConfigGenesisBlock, _ := configGenesisBlock.ConvertToTypes()
	if !reflect.DeepEqual(typesConfigGenesisBlock.Hash().Bytes(), genesisBlocksHash[types.CONFIGCHANNELID].Bytes()) {
		t.Fatal()
	}
	server.Stop()
	// initTestEnvironment(".data1")
}

func TestCreateChannel(t *testing.T) {
	var err error
	server, err = NewServer(getTestConfig())
	require.NoError(t, err)
	go func() {
		server.Start()
	}()
	time.Sleep(300 * time.Millisecond)
	// then try to create a channel
	client, err := getClient()
	require.NoError(t, err)

	pbTx := getCreateChannelTx("test")
	_, err = client.CreateChannel(context.Background(), &pb.CreateChannelRequest{
		Tx: pbTx,
	})
	require.NoError(t, err)
	// then stop
	server.Stop()
}

func TestServerRestartWithUserChannel(t *testing.T) {
	var err error
	server, err = NewServer(getTestConfig())
	require.NoError(t, err)

	go func() {
		server.Start()
	}()
	time.Sleep(300 * time.Millisecond)
	client, _ := getClient()
	// Then try to send a tx to test channel
	// then add a tx into test channel
	privKey, _ := crypto.NewPrivateKey(rawPrivKey)
	typesTx, err := types.NewTx("test", common.ZeroAddress, []byte("Just for test"), privKey)
	require.NoError(t, err)

	pbTx, err := pb.NewTx(typesTx)
	require.NoError(t, err)

	_, err = client.AddTx(context.Background(), &pb.AddTxRequest{
		Tx: pbTx,
	})
	require.NoError(t, err)

	server.Stop()
}

func TestFetchBlockAsync(t *testing.T) {
	var err error
	server, err = NewServer(getTestConfig())
	require.NoError(t, err)

	go func() {
		server.Start()
	}()
	time.Sleep(300 * time.Millisecond)
	client, _ := getClient()
	channelInfos, _ := client.ListChannels(context.Background(), &pb.ListChannelsRequest{
		System: true,
		PK:     pubKeyBytes,
	})
	require.Len(t, channelInfos.Channels, 3)
	var globalInfo *pb.ChannelInfo
	for _, channelInfo := range channelInfos.Channels {
		if channelInfo.ChannelID == types.GLOBALCHANNELID {
			globalInfo = channelInfo
			break
		}
		switch channelInfo.ChannelID {
		case types.GLOBALCHANNELID:
			globalInfo = channelInfo
			require.Equal(t, globalInfo.Identity, pb.Identity_MEMBER)
		case types.CONFIGCHANNELID:
			require.Equal(t, channelInfo.Identity, pb.Identity_MEMBER)
		case "test":
			require.Equal(t, channelInfo.Identity, pb.Identity_ADMIN)
		default:
			t.Fatalf("Unexcepted channel:%s", channelInfo.ChannelID)
		}
	}
	exceptNum := globalInfo.BlockSize
	// try to fecth the block sync which is not exist
	_, err = client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: types.GLOBALCHANNELID,
		Number:    exceptNum,
		Behavior:  pb.Behavior_FAIL_IF_NOT_READY,
	})
	require.Error(t, err)

	// Then async fetch block
	// first here is a block which is exist
	_, err = client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: types.GLOBALCHANNELID,
		Number:    0,
		Behavior:  pb.Behavior_RETURN_UNTIL_READY,
	})
	require.NoError(t, err)
	// then fetch the except block
	var wg sync.WaitGroup
	wg.Add(1)
	// fetch a block, of course this can not be done until a new block is generated.
	go func() {
		defer wg.Done()
		block, err := client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
			ChannelID: types.GLOBALCHANNELID,
			Number:    exceptNum,
			Behavior:  pb.Behavior_RETURN_UNTIL_READY,
		})
		require.NoError(t, err)
		require.Equal(t, block.Header.ChannelID, types.GLOBALCHANNELID)
		require.Equal(t, block.Header.Number, exceptNum)
	}()
	wg.Add(1)
	// add a new channel, which will create a new global block
	go func() {
		defer wg.Done()
		time.Sleep(500 * time.Millisecond)
		pbTx := getCreateChannelTx("async")
		_, err = client.CreateChannel(context.Background(), &pb.CreateChannelRequest{
			Tx: pbTx,
		})
		require.NoError(t, err)
	}()
	wg.Wait()
	server.Stop()
}

func TestAddDuplicateTxs(t *testing.T) {
	var err error
	server, err = NewServer(getTestConfig())
	require.NoError(t, err)

	go func() {
		server.Start()
	}()
	time.Sleep(300 * time.Millisecond)
	client, _ := getClient()
	// Then try to send a tx to test channel
	// then add a tx into test channel
	privKey, _ := crypto.NewPrivateKey(rawPrivKey)
	typesTx, err := types.NewTx("test", common.ZeroAddress, []byte("Duplicate"), privKey)
	require.NoError(t, err)

	pbTx, err := pb.NewTx(typesTx)
	require.NoError(t, err)

	_, err = client.AddTx(context.Background(), &pb.AddTxRequest{
		Tx: pbTx,
	})
	require.NoError(t, err)
	// then add the tx again
	_, err = client.AddTx(context.Background(), &pb.AddTxRequest{
		Tx: pbTx,
	})
	if !strings.Contains(err.Error(), "The tx exist in the blockchain aleardy") {
		t.Error(err)
	}
	server.Stop()
}

func TestEnd(t *testing.T) {
	initTestEnvironment(".data")
	initTestEnvironment(".data1")
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

func getCreateChannelTx(channelID string) *pb.Tx {
	admin, _ := types.NewMember(privKey.PubKey(), "admin")
	payload, _ := json.Marshal(cc.Payload{
		ChannelID: channelID,
		Profile: &cc.Profile{
			Public: true,
			Admins: []*types.Member{admin},
		},
		Version: 1,
	})
	privKey, _ := crypto.NewPrivateKey(rawPrivKey)
	typesTx, _ := types.NewTx(types.CONFIGCHANNELID, types.CreateChannelContractAddress, payload, privKey)

	pbTx, _ := pb.NewTx(typesTx)
	return pbTx
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

func initTestEnvironment(path string) {
	gopath := os.Getenv("GOPATH")
	dataPath, _ := util.MakeFileAbs(fmt.Sprintf("src/madledger/orderer/server/%s", path), gopath)
	os.RemoveAll(dataPath)
}
