package testfor2clients_bft

import (
	"fmt"
	"github.com/stretchr/testify/require"
	cc "madledger/client/config"
	client "madledger/client/lib"
	"madledger/common"
	"madledger/core/types"
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
	// client 0 and client 1 create channels concurrently
	client0 := bftClients[0]
	client1 := bftClients[1]
	var channels []string
	for m := 1; m <= 8; m++ {
		if m == 3 { // stop peer0
			go func(t *testing.T) {
				fmt.Println("Begin to stop peer 0 ... ")
				stopPeer(bftPeers[0])
				require.NoError(t, os.RemoveAll(getBFTPeerDataPath(0)))
			}(t)
		}
		if m == 6 { // restart peer0
			go func() {
				fmt.Println("Restart peer 0 ...")
				bftPeers[0]=startPeer(0)
			}()
		}
		// client 0 create channel
		channel := "test" + strconv.Itoa(m) + "0"
		channels = append(channels, channel)
		fmt.Printf("Create channel %s by client0 ...\n", channel)
		err := client0.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)

		// client 1 create channel
		channel = "test" + strconv.Itoa(m) + "1"
		channels = append(channels, channel)
		fmt.Printf("Create channel %s by client1 ...\n", channel)
		err = client1.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)

	}
	// compare tx, one is peer0 starting another is peer0 stopped
	time.Sleep(2 * time.Second)
	require.NoError(t, compareTxs())
}

func TestBFTCreateTx2(t *testing.T) {
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
				fmt.Println("Restart peer 0 ...")
				bftPeers[0]=startPeer(0)
			}()
		}
		// client 0 create contract
		contractCodes, err := readCodes(getBFTClientPath(0) + "/MyTest.bin")
		require.NoError(t, err)
		channel := "test" + strconv.Itoa(m) + "0"
		fmt.Printf("Create contract %d on channel %s by client0 ...\n", m, channel)
		tx, err := types.NewTx(channel, common.ZeroAddress, contractCodes, bftClients[0].GetPrivKey(), types.NORMAL)
		require.NoError(t, err)

		_, err = bftClients[0].AddTx(tx)
		require.NoError(t, err)

		// client 1 create channel
		contractCodes, err = readCodes(getBFTClientPath(1) + "/MyTest.bin")
		require.NoError(t, err)
		channel = "test" + strconv.Itoa(m) + "1"
		fmt.Printf("Create contract %d on channel %s by client1 ...\n", m, channel)
		tx, err = types.NewTx(channel, common.ZeroAddress, contractCodes, bftClients[1].GetPrivKey(), types.NORMAL)
		require.NoError(t, err)

		_, err = bftClients[1].AddTx(tx)
		require.NoError(t, err)
	}
	// compare tx, one is peer0 starting another is peer0 stopped
	time.Sleep(2 * time.Second)
	require.NoError(t, compareTxs())
}

func TestBFTCallTx2(t *testing.T) {
	// create channel test0 for client0 and test1 for client1
	require.NoError(t, createChannelForCallTx())

	// create smart contract on channel test0 and test1
	require.NoError(t, createContractForCallTx())

	for m := 1; m <= 8; m++ {
		if m == 2 { // stop peer0
			go func(t *testing.T) {
				fmt.Println("Begin to stop peer 0 ...")
				stopPeer(bftPeers[0])
				require.NoError(t, os.RemoveAll(getBFTPeerDataPath(0)))
			}(t)
		}
		if m == 5 { // restart peer0
			go func() {
				fmt.Println("Restart peer 0 ...")
				bftPeers[0]=startPeer(0)
			}()
		}

		// client0 call setNum and getNum function
		fmt.Printf("Call contract %d times on channel test0 ...\n", m)
		if m%2 == 0 {
			num := "1" + strconv.Itoa(m-1)
			require.NoError(t, getNumForCallTx(0, num))
		} else {
			num := "1" + strconv.Itoa(m)
			require.NoError(t, setNumForCallTx(0, num))
		}

		// client0 call setNum and getNum function
		fmt.Printf("Call contract %d on channel test1 ...\n", m)
		if m%2 == 0 {
			num := "2" + strconv.Itoa(m-1)
			require.NoError(t, getNumForCallTx(1, num))
		} else {
			num := "2" + strconv.Itoa(m)
			require.NoError(t, setNumForCallTx(1, num))
		}
	}
	// compare tx, one is peer0 starting another is peer0 stopped
	time.Sleep(2 * time.Second)
	require.NoError(t, compareTxs())
}

func TestBFTEnd2(t *testing.T) {
	for _, pid := range bftOrderers {
		stopOrderer(pid)
	}
	for _, pid := range bftPeers {
		stopPeer(pid)
	}

	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/bft"))
}
