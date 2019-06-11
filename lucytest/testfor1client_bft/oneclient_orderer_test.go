package testfor1client_bft

import (
	"fmt"
	cc "madledger/client/config"
	client "madledger/client/lib"
	cliu "madledger/client/util"
	"madledger/common"
	"madledger/common/abi"
	"madledger/core/types"
	peer "madledger/peer/server"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// change the package
func TestInitEnv1(t *testing.T) {
	require.NoError(t, initBFTEnvironment())
}

func TestBFTOrdererStart1(t *testing.T) {
	// then we can run orderers
	for i := range bftOrderers {
		pid := startOrderer(i)
		bftOrderers[i] = pid
	}
}

func TestBFTPeersStart1(t *testing.T) {
	for i := 0; i < 4; i++ {
		cfg := getPeerConfig(i)
		server, err := peer.NewServer(cfg)
		require.NoError(t, err)
		bftPeers[i] = server
	}

	for i := range bftPeers {
		go func(t *testing.T, i int) {
			err := bftPeers[i].Start()
			require.NoError(t, err)
		}(t, i)
	}

	time.Sleep(2 * time.Second)
}

func TestBFTLoadClients1(t *testing.T) {
	for i := range bftClients {
		clientPath := getBFTClientPath(i)
		cfgPath := getBFTClientConfigPath(i)
		cfg, err := cc.LoadConfig(cfgPath)
		require.NoError(t, err)
		re, _ := regexp.Compile("^.*[.]keystore")
		for i := range cfg.KeyStore.Keys {
			cfg.KeyStore.Keys[i] = clientPath + "/.keystore" + re.ReplaceAllString(cfg.KeyStore.Keys[i], "")
		}
		client, err := client.NewClientFromConfig(cfg)
		require.NoError(t, err)
		bftClients[i] = client
	}
}

func TestBFTCreateChannels1(t *testing.T) {
	// client-0 create 4 channels
	client0 := bftClients[0]
	for i := 0; i <= 2; i++ {
		channel := "test" + strconv.Itoa(i)
		err := client0.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
	}
	time.Sleep(2 * time.Second)

	// then we will check if channels created by client-0 are create successful
	// query by client-1
	require.NoError(t, listChannel(1))
}

func listChannel(node int) error {
	client := bftClients[node]
	infos, err := client.ListChannel(true)
	if err != nil {
		return err
	}
	table := cliu.NewTable()
	table.SetHeader("Name", "System", "BlockSize", "Identity")
	for _, info := range infos {
		table.AddRow(info.Name, info.System, info.BlockSize, info.Identity)
	}
	table.Render()

	return nil
}

// 关闭orderer 1，关闭期间通过client 0创建test3通道，然后重启orderer 1，查询数据
func TestBFTNodeRestart(t *testing.T) {
	stopOrderer(bftOrderers[1])
	os.RemoveAll(getBFTOrdererDataPath(1))

	//client 0创建test3通道
	client0 := bftClients[0]
	channel := "test3"
	err := client0.CreateChannel(channel, true, nil, nil)
	require.NoError(t, err)
	time.Sleep(2 * time.Second)

	fmt.Println("Restart orderer 1 ...")
	bftOrderers[1] = startOrderer(1)
	time.Sleep(2 * time.Second)

	// query by client 1
	require.NoError(t, listChannel(1))

}

func TestBFTCreateChannelAfterRestart(t *testing.T) {
	//client 1创建通道test4
	client1 := bftClients[1]
	channel := "test4"
	err := client1.CreateChannel(channel, true, nil, nil)
	require.NoError(t, err)
	time.Sleep(2 * time.Second)

	// query by client 1
	require.NoError(t, listChannel(1))
}

func TestBFTCreateTxAfterRestart(t *testing.T) {
	//client 1创建智能合约
	contractCodes, err := readCodes(getBFTClientPath(1) + "/MyTest.bin")
	require.NoError(t, err)
	client := bftClients[1]
	tx, err := types.NewTx("test4", common.ZeroAddress, contractCodes, client.GetPrivKey())
	require.NoError(t, err)

	status, err := client.AddTx(tx)
	require.NoError(t, err)

	// Then print the status
	table := cliu.NewTable()
	table.SetHeader("BlockNumber", "BlockIndex", "ContractAddress")
	if status.Err != "" {
		table.AddRow(status.BlockNumber, status.BlockIndex, status.Err)
	} else {
		table.AddRow(status.BlockNumber, status.BlockIndex, status.ContractAddress)
	}
	table.Render()
}

func TestBFTCallTxAfterRestart(t *testing.T) {
	//client 1调用智能合约
	abiPath := fmt.Sprintf(getBFTClientPath(1) + "/MyTest.abi")
	funcName := "getNum"
	var inputs []string = make([]string, 0)
	payloadBytes, err := abi.GetPayloadBytes(abiPath, funcName, inputs)
	require.NoError(t, err)

	client := bftClients[1]
	tx, err := types.NewTx("test4", common.HexToAddress("0x0619e2393802cc99e90cf892b92a113f19af5887"), payloadBytes, client.GetPrivKey())
	require.NoError(t, err)

	status, err := client.AddTx(tx)
	require.NoError(t, err)

	// Then print the status
	table := cliu.NewTable()
	table.SetHeader("BlockNumber", "BlockIndex", "Output")
	if status.Err != "" {
		table.AddRow(status.BlockNumber, status.BlockIndex, status.Err)
	} else {
		values, err := abi.Unpacker(abiPath, funcName, status.Output)
		require.NoError(t, err)
		var output []string
		for _, value := range values {
			output = append(output, value.Value)
		}
		table.AddRow(status.BlockNumber, status.BlockIndex, output)
	}
	table.Render()
}

func TestBFTEnd1(t *testing.T) {
	for _, pid := range bftOrderers {
		stopOrderer(pid)
	}

	for i := range bftPeers {
		bftPeers[i].Stop()
	}
	time.Sleep(2 * time.Second)
	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/bft"))
}
