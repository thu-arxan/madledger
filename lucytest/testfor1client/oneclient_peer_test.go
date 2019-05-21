package testfor1client

import (
	"fmt"
	cc "madledger/client/config"
	client "madledger/client/lib"
	"madledger/common"
	"madledger/core/types"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

/*var (
	bftOrderers [4]string
	bftClients  [4]*client.Client
	bftPeers    [4]*peer.Server
)
*/

func TestInitEnv2(t *testing.T) {
	require.NoError(t, initBFTEnvironment())
}

func TestBFTOrdererStart2(t *testing.T) {
	// then we can run orderers
	for i := range bftOrderers {
		pid := startOrderer(i)
		bftOrderers[i] = pid
	}
}

func TestBFTPeersStart2(t *testing.T) {
	for i := 0; i < 4; i++ {
		require.NoError(t, initPeer(i))
	}

	for i := range bftPeers {
		go func(t *testing.T, i int) {
			err := bftPeers[i].Start()
			require.NoError(t, err)
		}(t, i)
	}

	time.Sleep(2 * time.Second)
}

func TestBFTLoadClients2(t *testing.T) {
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

func TestBFTCreateChannels2(t *testing.T) {
	// client 0 and client 1 create channels concurrently
	client := bftClients[0]
	var channels []string
	for m := 1; m <= 8; m++ {
		if m == 4 { // stop peer0
			go func(t *testing.T) {
				fmt.Println("Begin to stop peer 0")
				bftPeers[0].Stop()
				require.NoError(t, os.RemoveAll(getBFTPeerDataPath(0)))
			}(t)
		}
		if m == 6 { // restart peer0
			require.NoError(t, initPeer(0))

			go func(t *testing.T) {
				err := bftPeers[0].Start()
				require.NoError(t, err)
				fmt.Printf("Restart orderer0 successfully ...\n")
			}(t)
		}
		// client 0 create channel
		channel := "test" + strconv.Itoa(m) + "0"
		channels = append(channels, channel)
		fmt.Printf("Create channel %s ...\n", channel)
		err := client.CreateChannel(channel, true, nil, nil)
		fmt.Printf("Create channel %s done\n", channel)
		require.NoError(t, err)
	}
	time.Sleep(5 * time.Second)

	// then we will check if channels are create successful
	require.NoError(t, compareChannels(channels))
}

func TestBFTCreateTx(t *testing.T) {
	client := bftClients[0]
	for m := 1; m <= 6; m++ {
		if m == 2 { // stop peer0
			go func(t *testing.T) {
				fmt.Println("Begin to stop peer 0")
				bftPeers[0].Stop()
				require.NoError(t, os.RemoveAll(getBFTPeerDataPath(0)))
			}(t)
		}
		if m == 4 { // restart peer0
			require.NoError(t, initPeer(0))

			go func(t *testing.T) {
				err := bftPeers[0].Start()
				require.NoError(t, err)
				fmt.Printf("Restart peer0 successfully ...\n")
			}(t)
		}
		// client 0 create contract
		contractCodes, err := readCodes(getBFTClientPath(0) + "/MyTest.bin")
		require.NoError(t, err)
		channel := "test" + strconv.Itoa(m) + "0"
		fmt.Printf("Create contract %d on channel %s ...\n", m, channel)
		tx, err := types.NewTx(channel, common.ZeroAddress, contractCodes, client.GetPrivKey())
		require.NoError(t, err)

		_, err = client.AddTx(tx)
		require.NoError(t, err)
	}
	time.Sleep(5 * time.Second)
}

func TestBFTCallTx(t *testing.T) {
	// 为client0和client1分别创建test0
	require.NoError(t, createChannelForCallTx())

	// 在test0上创建合约
	require.NoError(t, createContractForCallTx())

	for m := 1; m <= 6; m++ {
		if m == 2 { // stop peer0
			go func(t *testing.T) {
				fmt.Println("Begin to stop peer 0")
				bftPeers[0].Stop()
				require.NoError(t, os.RemoveAll(getBFTPeerDataPath(0)))
			}(t)
		}
		if m == 4 { // restart peer0
			require.NoError(t, initPeer(0))

			go func(t *testing.T) {
				err := bftPeers[0].Start()
				require.NoError(t, err)
				fmt.Printf("Restart peer0 successfully ...\n")
			}(t)
		}

		// client0调用合约的setNum
		fmt.Printf("Call contract %d times on channel test0 ...\n", m)
		if m%2 == 0 {
			num := "1" + strconv.Itoa(m-1)
			require.NoError(t, getNumForCallTx(0, num))
		} else {
			num := "1" + strconv.Itoa(m)
			require.NoError(t, setNumForCallTx(0, num))
		}
	}

	time.Sleep(5 * time.Second)

}

func TestBFTEnd2(t *testing.T) {
	for _, pid := range bftOrderers {
		stopOrderer(pid)
	}

	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/bft"))
}
