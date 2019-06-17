package testfor2clients_bft

import (
	"fmt"
	"github.com/stretchr/testify/require"
	cc "madledger/client/config"
	client "madledger/client/lib"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"
)

// change the package
func TestInitEnv1(t *testing.T) {
	require.NoError(t, initRAFTEnvironment())
}

func TestBFTOrdererStart1(t *testing.T) {
	// then we can run orderers
	for i := range raftOrderers {
		pid := startOrderer(i)
		raftOrderers[i] = pid
	}
}

func TestRAFTPeersStart1(t *testing.T) {
	for i := range raftPeers {
		require.NoError(t, initPeer(i))
	}

	for i := range raftPeers {
		go func(t *testing.T, i int) {
			err := raftPeers[i].Start()
			require.NoError(t, err)
		}(t, i)
	}
	time.Sleep(2 * time.Second)
}

func TestLoadClients1(t *testing.T) {
	for i := range raftClients {
		clientPath := getRAFTClientPath(i)
		cfgPath := getRAFTClientConfigPath(i)
		cfg, err := cc.LoadConfig(cfgPath)
		require.NoError(t, err)
		re, _ := regexp.Compile("^.*[.]keystore")
		for i := range cfg.KeyStore.Keys {
			cfg.KeyStore.Keys[i] = clientPath + "/.keystore" + re.ReplaceAllString(cfg.KeyStore.Keys[i], "")
		}
		client, err := client.NewClientFromConfig(cfg)
		require.NoError(t, err)
		raftClients[i] = client
	}
}

func TestRaftCreateChannels1(t *testing.T) {
	// client-0 create 4 channels
	client := raftClients[0]
	var channels []string
	for i := 0; i < 8; i++ {
		if i == 4 {
			fmt.Println("Stop Orderer 0 ...")
			stopOrderer(raftOrderers[0])
			// restart orderer0 by RestartNode
			require.NoError(t, os.RemoveAll(getRAFTOrdererDataPath(0)))
		}
		if i == 6 {
			fmt.Println("Restart Orderer 0 ...")
			raftOrderers[0] = startOrderer(0)
		}
		// client 0 create channel
		channel := "test" + strconv.Itoa(i)
		channels = append(channels, channel)
		fmt.Printf("Create channel %s ...\n", channel)
		err := client.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
	}
	time.Sleep(2 * time.Second)

	// then we will check if channels are create successful
	require.NoError(t, compareChannels(channels))
}

/*func TestRaftCreateTx1(t *testing.T) {
	client := raftClients[0]
	for i := 0; i < 8; i++ {
		if i == 3 {
			fmt.Println("Stop Orderer 0 ...")
			stopOrderer(raftOrderers[0])
			// restart orderer0 by RestartNode
			require.NoError(t, os.RemoveAll(getRAFTOrdererDataPath(0)))
		}
		if i == 5 { // restart orderer0
			fmt.Println("Restart Orderer 0 ...")
			raftOrderers[0] = startOrderer(0)
		}
		// client 0 create contract
		contractCodes, err := readCodes(getRAFTClientPath(0) + "/MyTest.bin")
		require.NoError(t, err)
		channel := "test" + strconv.Itoa(i)
		fmt.Printf("Create contract %d on channel %s ...\n", i, channel)
		tx, err := types.NewTx(channel, common.ZeroAddress, contractCodes, client.GetPrivKey())
		require.NoError(t, err)

		_, err = client.AddTx(tx)
		require.NoError(t, err)
	}
	time.Sleep(2 * time.Second)
}*/

/*func TestRaftCallTx1(t *testing.T) {
	for i := 1; i <= 6; i++ {
		if i == 3 {
			fmt.Println("Stop Orderer 0 ...")
			stopOrderer(raftOrderers[0])
			// restart orderer0 by RestartNode
			require.NoError(t, os.RemoveAll(getRAFTOrdererDataPath(0)))
		}
		if i == 4 {
			fmt.Println("Restart Orderer 0 ...")
			raftOrderers[0] = startOrderer(0)
		}

		// client0调用合约的setNum
		fmt.Printf("Call contract %d times on channel test0 ...\n", i)
		if i%2 == 0 {
			num := "1" + strconv.Itoa(i-1)
			require.NoError(t, getNumForCallTx(num))
		} else {
			num := "1" + strconv.Itoa(i)
			require.NoError(t, setNumForCallTx(num))
		}
	}
}*/

/*
func TestBFTEnd1(t *testing.T) {
	for _, pid := range bftOrderers {
		stopOrderer(pid)
	}

	for i := range bftPeers {
		bftPeers[i].Stop()
	}
	time.Sleep(2 * time.Second)
	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/raft"))
}*/

/*
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
*/
