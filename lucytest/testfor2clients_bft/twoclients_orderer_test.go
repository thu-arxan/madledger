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
	// then we can run peers
	for i := range bftPeers {
		pid := startPeer(i)
		bftPeers[i] = pid
	}
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
	// client 0 and client 1 create channels concurrently
	for m := 0; m < 8; m++ {
		if m == 2 { // stop orderer0
			go func(t *testing.T) {
				fmt.Println("Begin to stop orderer 0 ...")
				stopOrderer(bftOrderers[0])
				require.NoError(t, os.RemoveAll(getBFTOrdererDataPath(0)))
			}(t)
		}
		if m == 5 { // restart orderer0
			go func() {
				fmt.Println("Restart orderer 0 ...")
				bftOrderers[0] = startOrderer(0)
			}()
		}
		// client 0 create channel
		channel := "test0"
		if m != 0 {
			channel = "test0" + strconv.Itoa(m)
		}
		fmt.Printf("Create channel %s by client0 ...\n", channel)
		err := bftClients[0].CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)

		// client 1 create channel
		channel = "test1"
		if m != 1 {
			channel = "test1" + strconv.Itoa(m)
		}
		fmt.Printf("Create channel %s by client1 ...\n", channel)
		err = bftClients[1].CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)

	}
	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels())
}

func TestBFTCreateTx1(t *testing.T) {
	for m := 0; m < 8; m++ {
		if m == 2 { // stop orderer0
			go func(t *testing.T) {
				fmt.Println("Begin to stop orderer 0 ...")
				stopOrderer(bftOrderers[0])
				require.NoError(t, os.RemoveAll(getBFTOrdererDataPath(0)))
			}(t)
		}
		if m == 5 { // restart orderer0
			go func() {
				fmt.Println("Restart orderer 0 ...")
				bftOrderers[0] = startOrderer(0)
			}()
		}
		// client 0 create contract
		channel := "test0"
		if m != 0 {
			channel = "test0" + strconv.Itoa(m)
		}
		contractCodes, err := readCodes(getBFTClientPath(0) + "/MyTest.bin")
		require.NoError(t, err)
		fmt.Printf("Create contract %d on channel %s by client0 ...\n", m, channel)
		tx, err := types.NewTx(channel, common.ZeroAddress, contractCodes, bftClients[0].GetPrivKey(), types.NORMAL)
		require.NoError(t, err)

		_, err = bftClients[0].AddTx(tx)
		require.NoError(t, err)

		// client 1 create channel
		channel = "test1"
		if m != 1 {
			channel = "test1" + strconv.Itoa(m)
		}
		contractCodes, err = readCodes(getBFTClientPath(1) + "/MyTest.bin")
		require.NoError(t, err)
		fmt.Printf("Create contract %d on channel %s by client1 ...\n", m, channel)
		tx, err = types.NewTx(channel, common.ZeroAddress, contractCodes, bftClients[1].GetPrivKey(), types.NORMAL)
		require.NoError(t, err)

		_, err = bftClients[1].AddTx(tx)
		require.NoError(t, err)
	}
	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels())
}

func TestBFTCallTx1(t *testing.T) {
	for m := 1; m <= 8; m++ {
		if m == 3 { // stop orderer0
			go func(t *testing.T) {
				fmt.Println("Begin to stop orderer 0 ...")
				stopOrderer(bftOrderers[0])
				require.NoError(t, os.RemoveAll(getBFTOrdererDataPath(0)))
			}(t)
		}
		if m == 6 { // restart orderer0
			go func() {
				fmt.Println("Restart orderer 0 ...")
				bftOrderers[0] = startOrderer(0)
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

		// client1 call setNum and getNum function
		fmt.Printf("Call contract %d on channel test1 ...\n", m)
		if m%2 == 0 {
			num := "2" + strconv.Itoa(m-1)
			require.NoError(t, getNumForCallTx(1, num))
		} else {
			num := "2" + strconv.Itoa(m)
			require.NoError(t, setNumForCallTx(1, num))
		}
	}
	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels())
}

func TestBFTEnd1(t *testing.T) {
	for _, pid := range bftOrderers {
		stopOrderer(pid)
	}
	for _, pid := range bftPeers {
		stopPeer(pid)
	}

	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/bft"))
}
