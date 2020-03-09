package testfor1client_bft

import (
	"fmt"
	cc "madledger/client/config"
	client "madledger/client/lib"
	cliu "madledger/client/util"
	"madledger/common"
	"madledger/common/abi"
	"madledger/core"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// change the package
func TestInitEnv1OB(t *testing.T) {
	require.NoError(t, initBFTEnvironment())
}

func TestBFTOrdererStart1OB(t *testing.T) {
	// then we can run orderers
	for i := range bftOrderers {
		pid := startOrderer(i)
		bftOrderers[i] = pid
	}
}

func TestBFTPeersStart1OB(t *testing.T) {
	// then we can run peers
	for i := range bftPeers {
		pid := startPeer(i)
		bftPeers[i] = pid
	}
}

func TestBFTLoadClients1OB(t *testing.T) {
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

func TestBFTCreateChannels1OB(t *testing.T) {
	// client0 create 3 channels
	for i := 0; i < 3; i++ {
		channel := "test" + strconv.Itoa(i)
		err := bftClients[0].CreateChannel(channel, true, nil, nil, 1, 1, 10000000)
		require.NoError(t, err)
	}

	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels())
}

// stop orderer1 while client0 create channel test3
// then restart orderer1 and compare channel infos
func TestBFTNodeRestart1OB(t *testing.T) {
	stopOrderer(bftOrderers[1])
	os.RemoveAll(getBFTOrdererDataPath(1))

	//client0 create channel test3
	fmt.Println("Create channel test3 ...")
	channel := "test3"
	err := bftClients[0].CreateChannel(channel, true, nil, nil, 1, 1, 10000000)
	require.NoError(t, err)

	fmt.Println("Restart orderer 1 ...")
	bftOrderers[1] = startOrderer(1)

	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels())

}

func TestBFTCreateChannelAfterRestart1OB(t *testing.T) {
	//client1 create channel test4
	channel := "test4"
	err := bftClients[1].CreateChannel(channel, true, nil, nil, 1, 1, 10000000)
	require.NoError(t, err)
}

func TestBFTCreateTxAfterRestart1OB(t *testing.T) {
	//client1 create smart contract
	contractCodes, err := readCodes(getBFTClientPath(1) + "/MyTest.bin")
	require.NoError(t, err)
	client := bftClients[1]
	tx, err := core.NewTx("test4", common.ZeroAddress, contractCodes, 0, "", client.GetPrivKey())
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

func TestBFTCallTxAfterRestart1OB(t *testing.T) {
	//client1 call smart contract
	abiPath := fmt.Sprintf(getBFTClientPath(1) + "/MyTest.abi")
	funcName := "getNum"
	var inputs = make([]string, 0)
	payloadBytes, err := abi.Pack(abiPath, funcName, inputs...)
	require.NoError(t, err)

	tx, err := core.NewTx("test4", common.HexToAddress("0x0619e2393802cc99e90cf892b92a113f19af5887"),
		payloadBytes, 0, "", bftClients[1].GetPrivKey())
	require.NoError(t, err)

	_, err = bftClients[1].AddTx(tx)
	require.NoError(t, err)

	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels())
}

func TestBFTEnd1OB(t *testing.T) {
	for _, pid := range bftOrderers {
		stopOrderer(pid)
	}
	for _, pid := range bftPeers {
		stopPeer(pid)
	}

	time.Sleep(2 * time.Second)
	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/bft"))
}
