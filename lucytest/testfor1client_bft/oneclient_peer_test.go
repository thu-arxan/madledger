package testfor1client_bft

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
	// then we can run peers
	for i := range bftPeers {
		pid := startPeer(i)
		bftPeers[i] = pid
	}
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
	for m := 1; m <= 8; m++ {
		if m == 3 { // stop peer0
			go func(t *testing.T) {
				fmt.Println("Begin to stop peer 0 ...")
				stopPeer(bftPeers[0])
				require.NoError(t, os.RemoveAll(getBFTPeerDataPath(0)))
			}(t)
		}
		if m == 6 { // restart peer0
			go func() {
				fmt.Println("Begin to restart peer 0 ...")
				bftPeers[0]=startPeer(0)
			}()
		}

		// client0 create channel
		channel := "test" + strconv.Itoa(m) + "0"
		fmt.Printf("Create channel %s ...\n", channel)
		err := bftClients[0].CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
	}

	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareTxs())
}






/*
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
				fmt.Println("Begin to restart peer 0")
				err := bftPeers[0].Start()
				require.NoError(t, err)
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
*/