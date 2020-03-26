// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package server

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"madledger/blockchain/asset"
	cc "madledger/blockchain/config"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/core"
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
	privKey, _           = crypto.NewPrivateKey(rawPrivKey, crypto.KeyAlgoSecp256k1)
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
		require.NoError(t, server.Start())
	}()
	time.Sleep(500 * time.Millisecond)
}

func TestListChannelsAtNil(t *testing.T) {
	client, err := getClient()
	require.NoError(t, err)

	infos, err := client.ListChannels(context.Background(), &pb.ListChannelsRequest{
		System: true,
		PK:     pubKeyBytes,
		Algo:   privKey.PubKey().Algo(),
	})
	require.NoError(t, err)
	require.Len(t, infos.Channels, 3)

	for _, channel := range infos.Channels {
		switch channel.ChannelID {
		case core.GLOBALCHANNELID, core.CONFIGCHANNELID, core.ASSETCHANNELID:
			require.Equal(t, channel.BlockSize, uint64(1))
		default:
			t.Fatal(fmt.Errorf("Unknown channel %s", channel.ChannelID))
		}
	}
	infos, err = client.ListChannels(context.Background(), &pb.ListChannelsRequest{
		System: false,
		PK:     pubKeyBytes,
		Algo:   privKey.PubKey().Algo(),
	})
	require.NoError(t, err)
	require.Len(t, infos.Channels, 0)
}

func TestFetchBlockAtNil(t *testing.T) {
	client, err := getClient()
	require.NoError(t, err)

	globalGenesisBlock, err := client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: core.GLOBALCHANNELID,
		Number:    0,
	})
	require.NoError(t, err)
	require.Equal(t, globalGenesisBlock.Header.Number, uint64(0))
	// set global genesis block hash
	typesGlobalGenesisBlock, err := globalGenesisBlock.ToCore()
	require.NoError(t, err)

	genesisBlocksHash[core.GLOBALCHANNELID] = typesGlobalGenesisBlock.Hash()
	// test bigger number
	_, err = client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: core.GLOBALCHANNELID,
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
		ChannelID: core.CONFIGCHANNELID,
		Number:    0,
	})
	require.NoError(t, err)
	require.Equal(t, configGenesisBlock.Header.Number, uint64(0))
	// set config genesis block hash
	typesConfigGenesisBlock, err := configGenesisBlock.ToCore()
	require.NoError(t, err)

	genesisBlocksHash[core.CONFIGCHANNELID] = typesConfigGenesisBlock.Hash()

	// get genesis block of asset
	assetGenesisBlock, err := client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: core.ASSETCHANNELID,
		Number:    0,
	})
	require.NoError(t, err)
	require.Equal(t, assetGenesisBlock.Header.Number, uint64(0))
	// set asset genesis block hash
	typesAssetGenesisBlock, err := assetGenesisBlock.ToCore()
	require.NoError(t, err)

	genesisBlocksHash[core.ASSETCHANNELID] = typesAssetGenesisBlock.Hash()
	server.Stop()
}

func TestServerRestart(t *testing.T) {
	var err error
	server, err = NewServer(getTestConfig())
	require.NoError(t, err)

	go func() {
		require.NoError(t, server.Start())
	}()
	time.Sleep(500 * time.Millisecond)
	server.Stop()
}

func TestServerStartAtAnotherPath(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	cfg, _ := config.LoadConfig(getTestConfigFilePath())
	chainPath, _ := util.MakeFileAbs("src/madledger/orderer/server/.data1/blocks", gopath)
	dbPath, _ := util.MakeFileAbs("src/madledger/orderer/server/.data1/leveldb", gopath)
	cfg.BlockChain.Path = chainPath
	cfg.DB.LevelDB.Path = dbPath
	var err error
	server, err = NewServer(cfg)
	require.NoError(t, err)

	go func(t *testing.T) {
		require.NoError(t, server.Start())
	}(t)
	time.Sleep(500 * time.Millisecond)
	client, err := getClient()
	require.NoError(t, err)
	// compare global genesis block
	globalGenesisBlock, _ := client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: core.GLOBALCHANNELID,
		Number:    0,
	})
	typesGlobalGenesisBlock, _ := globalGenesisBlock.ToCore()
	if !reflect.DeepEqual(typesGlobalGenesisBlock.Hash().Bytes(), genesisBlocksHash[core.GLOBALCHANNELID].Bytes()) {
		t.Fatal()
	}
	// compare config genesis block
	configGenesisBlock, _ := client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: core.CONFIGCHANNELID,
		Number:    0,
	})
	typesConfigGenesisBlock, _ := configGenesisBlock.ToCore()
	if !reflect.DeepEqual(typesConfigGenesisBlock.Hash().Bytes(), genesisBlocksHash[core.CONFIGCHANNELID].Bytes()) {
		t.Fatal()
	}
	// compare asset genesis block
	assetGenesisBlock, _ := client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: core.ASSETCHANNELID,
		Number:    0,
	})
	typesAssetGenesisBlock, _ := assetGenesisBlock.ToCore()
	if !reflect.DeepEqual(typesAssetGenesisBlock.Hash().Bytes(), genesisBlocksHash[core.ASSETCHANNELID].Bytes()) {
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
		require.NoError(t, server.Start())
	}()
	time.Sleep(500 * time.Millisecond)
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
		require.NoError(t, server.Start())
	}()
	time.Sleep(500 * time.Millisecond)
	client, _ := getClient()
	// Then try to send a tx to test channel
	// then add a tx into test channel
	privKey, _ := crypto.NewPrivateKey(rawPrivKey, crypto.KeyAlgoSecp256k1)
	coreTx, err := core.NewTx("test", common.ZeroAddress, []byte("Just for test"), 0, "", privKey)
	require.NoError(t, err)

	pbTx, err := pb.NewTx(coreTx)
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
		require.NoError(t, server.Start())
	}()
	time.Sleep(500 * time.Millisecond)
	client, _ := getClient()
	channelInfos, _ := client.ListChannels(context.Background(), &pb.ListChannelsRequest{
		System: true,
		PK:     pubKeyBytes,
		Algo:   privKey.PubKey().Algo(),
	})
	require.Len(t, channelInfos.Channels, 4)
	var globalInfo *pb.ChannelInfo
	for _, channelInfo := range channelInfos.Channels {
		if channelInfo.ChannelID == core.GLOBALCHANNELID {
			globalInfo = channelInfo
			break
		}
		switch channelInfo.ChannelID {
		case core.GLOBALCHANNELID:
			globalInfo = channelInfo
			require.Equal(t, globalInfo.Identity, pb.Identity_MEMBER)
		case core.CONFIGCHANNELID:
			require.Equal(t, channelInfo.Identity, pb.Identity_MEMBER)
		case core.ASSETCHANNELID:
			require.Equal(t, channelInfo.Identity, pb.Identity_MEMBER)
		case "test":
			require.Equal(t, channelInfo.Identity, pb.Identity_ADMIN)
		default:
			t.Fatalf("Unexpected channel:%s", channelInfo.ChannelID)
		}
	}
	expectNum := globalInfo.BlockSize
	// try to fecth the block sync which is not exist
	_, err = client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: core.GLOBALCHANNELID,
		Number:    expectNum,
		Behavior:  pb.Behavior_FAIL_IF_NOT_READY,
	})
	require.Error(t, err)

	// Then async fetch block
	// first here is a block which is exist
	_, err = client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: core.GLOBALCHANNELID,
		Number:    0,
		Behavior:  pb.Behavior_RETURN_UNTIL_READY,
	})
	require.NoError(t, err)
	// then fetch the expect block
	var wg sync.WaitGroup
	wg.Add(1)
	// fetch a block, of course this can not be done until a new block is generated.
	go func() {
		defer wg.Done()
		block, err := client.FetchBlock(context.Background(), &pb.FetchBlockRequest{
			ChannelID: core.GLOBALCHANNELID,
			Number:    expectNum,
			Behavior:  pb.Behavior_RETURN_UNTIL_READY,
		})
		require.NoError(t, err)
		require.Equal(t, block.Header.ChannelID, core.GLOBALCHANNELID)
		require.Equal(t, block.Header.Number, expectNum)
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
		require.NoError(t, server.Start())
	}()
	time.Sleep(500 * time.Millisecond)
	client, _ := getClient()
	// Then try to send a tx to test channel
	// then add a tx into test channel
	privKey, _ := crypto.NewPrivateKey(rawPrivKey, crypto.KeyAlgoSecp256k1)
	coreTx, err := core.NewTx("test", common.ZeroAddress, []byte("Duplicate"), 0, "", privKey)
	require.NoError(t, err)

	pbTx, err := pb.NewTx(coreTx)
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

func TestAsset(t *testing.T) {
	//1.test init
	var err error
	server, err = NewServer(getTestConfig())
	require.NoError(t, err)

	go func() {
		require.NoError(t, server.Start())
	}()
	time.Sleep(500 * time.Millisecond)

	client, _ := getClient()

	//2.test issue
	algo := crypto.KeyAlgoSecp256k1

	issuerKey, err := crypto.GeneratePrivateKey(algo)
	require.NoError(t, err)
	falseIssuerKey, err := crypto.GeneratePrivateKey(algo)
	require.NoError(t, err)
	require.NotEqual(t, issuerKey, falseIssuerKey)
	receiverKey, err := crypto.GeneratePrivateKey(algo)
	require.NoError(t, err)

	issuer, err := issuerKey.PubKey().Address()
	require.NoError(t, err)
	falseIssuer, err := falseIssuerKey.PubKey().Address()
	require.NoError(t, err)
	receiver, err := receiverKey.PubKey().Address()

	//issue to issuer itself
	pbTx := getAssetChannelTx(core.IssueContractAddress, issuer, "", uint64(10), issuerKey)
	_, err = client.AddTx(context.Background(), &pb.AddTxRequest{
		Tx: pbTx,
	})
	require.NoError(t, err)

	acc, err := client.GetAccountInfo(context.Background(), &pb.GetAccountInfoRequest{
		Address: issuer.Bytes(),
	})
	require.NoError(t, err)
	require.Equal(t, uint64(10), acc.GetBalance())

	//falseissuer issue fail
	pbTx = getAssetChannelTx(core.IssueContractAddress, falseIssuer, "", uint64(10), falseIssuerKey)
	_, err = client.AddTx(context.Background(), &pb.AddTxRequest{
		Tx: pbTx,
	})
	require.NoError(t, err)
	acc, err = client.GetAccountInfo(context.Background(), &pb.GetAccountInfoRequest{
		Address: falseIssuer.Bytes(),
	})
	require.NoError(t, err)
	require.Equal(t, uint64(0), acc.GetBalance())

	//test issue to channel
	pbTx = getAssetChannelTx(core.IssueContractAddress, common.ZeroAddress, "test", uint64(10), issuerKey)
	_, err = client.AddTx(context.Background(), &pb.AddTxRequest{
		Tx: pbTx,
	})
	require.NoError(t, err)
	acc, err = client.GetAccountInfo(context.Background(), &pb.GetAccountInfoRequest{
		Address: common.BytesToAddress([]byte("test")).Bytes(),
	})
	require.NoError(t, err)
	require.Equal(t, uint64(10), acc.GetBalance())

	//3.test transfer
	pbTx = getAssetChannelTx(core.TransferContractrAddress, receiver, "", uint64(5), issuerKey)
	_, err = client.AddTx(context.Background(), &pb.AddTxRequest{
		Tx: pbTx,
	})
	require.NoError(t, err)
	acc, err = client.GetAccountInfo(context.Background(), &pb.GetAccountInfoRequest{
		Address: receiver.Bytes(),
	})
	require.NoError(t, err)
	require.Equal(t, uint64(5), acc.GetBalance())

	//test transfer fail
	pbTx = getAssetChannelTx(core.TransferContractrAddress, receiver, "", uint64(5), falseIssuerKey)
	_, err = client.AddTx(context.Background(), &pb.AddTxRequest{
		Tx: pbTx,
	})
	require.NoError(t, err)
	acc, err = client.GetAccountInfo(context.Background(), &pb.GetAccountInfoRequest{
		Address: receiver.Bytes(),
	})
	require.NoError(t, err)
	require.Equal(t, uint64(5), acc.GetBalance())

	//4.test exchangeToken a.k.a transfer to channel in orderer execution
	pbTx = getAssetChannelTx(core.TransferContractrAddress, common.ZeroAddress, "test", uint64(5), receiverKey)
	_, err = client.AddTx(context.Background(), &pb.AddTxRequest{
		Tx: pbTx,
	})
	require.NoError(t, err)
	acc, err = client.GetAccountInfo(context.Background(), &pb.GetAccountInfoRequest{
		Address: common.BytesToAddress([]byte("test")).Bytes(),
	})
	require.NoError(t, err)
	require.Equal(t, uint64(15), acc.GetBalance())
	server.Stop()
}

// func TestChargeBlock(t *testing.T) {
// 	var err error
// 	server, err = NewServer(getTestConfig())
// 	require.NoError(t, err)

// 	go func() {
// 		require.NoError(t, server.Start())
// 	}()
// 	time.Sleep(500 * time.Millisecond)
// 	client, _ := getClient()
// 	// Then try to send a tx to test channel
// 	// then add a tx into test channel
// 	privKey, _ := crypto.NewPrivateKey(rawPrivKey, crypto.KeyAlgoSecp256k1)

// 	coreTx, err := core.NewTx("test", common.ZeroAddress, []byte("success"), 0, "", privKey)
// 	require.NoError(t, err)
// 	pbTx, err := pb.NewTx(coreTx)
// 	require.NoError(t, err)
// 	_, err = client.AddTx(context.Background(), &pb.AddTxRequest{
// 		Tx: pbTx,
// 	})
// 	require.NoError(t, err)

// 	//change BlockPrice of test channel's
// 	payload, err := json.Marshal(cc.Payload{
// 		ChannelID:  "test",
// 		BlockPrice: uint64(100000),
// 	})
// 	require.NoError(t, err)
// 	coreTx, err = core.NewTx(core.CONFIGCHANNELID, common.ZeroAddress, payload, 0, "", privKey)
// 	pbTx, err = pb.NewTx(coreTx)
// 	require.NoError(t, err)
// 	_, err = client.AddTx(context.Background(), &pb.AddTxRequest{
// 		Tx: pbTx,
// 	})
// 	require.NoError(t, err)

// 	//now add block should fail
// 	coreTx, err = core.NewTx("test", common.ZeroAddress, []byte("fail"), 0, "", privKey)
// 	require.NoError(t, err)
// 	pbTx, err = pb.NewTx(coreTx)
// 	require.NoError(t, err)
// 	_, err = client.AddTx(context.Background(), &pb.AddTxRequest{
// 		Tx: pbTx,
// 	})
// 	require.NoError(t, err)
// 	server.Stop()
// }

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
	admin, _ := core.NewMember(privKey.PubKey(), "admin")
	payload, _ := json.Marshal(cc.Payload{
		ChannelID: channelID,
		Profile: &cc.Profile{
			Public: true,
			Admins: []*core.Member{admin},
		},
		Version:         1,
		GasPrice:        0,
		AssetTokenRatio: 1,
		MaxGas:          10000000,
	})
	privKey, _ := crypto.NewPrivateKey(rawPrivKey, crypto.KeyAlgoSecp256k1)
	coreTx, _ := core.NewTx(core.CONFIGCHANNELID, core.CreateChannelContractAddress, payload, 0, "", privKey)

	pbTx, _ := pb.NewTx(coreTx)
	return pbTx
}

func getAssetChannelTx(contract, addressInPayload common.Address, channelInPayload string, value uint64, privKey crypto.PrivateKey) *pb.Tx {
	payload, _ := json.Marshal(asset.Payload{
		Address:   addressInPayload,
		ChannelID: channelInPayload,
	})
	coreTx, _ := core.NewTx(core.ASSETCHANNELID, contract, payload, value, "", privKey)
	pbTx, _ := pb.NewTx(coreTx)
	return pbTx
}

func getTestConfig() *config.Config {
	cfg, _ := config.LoadConfig(getTestConfigFilePath())
	cfg.BlockChain.Path = getTestChainPath()
	cfg.DB.LevelDB.Path = getTestDBPath()
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
