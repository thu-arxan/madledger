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

func TestInitEnv1PB(t *testing.T) {
	require.NoError(t, initBFTEnvironment())
}

func TestBFTOrdererStart1PB(t *testing.T) {
	// then we can run orderers
	for i := range bftOrderers {
		pid := startOrderer(i)
		bftOrderers[i] = pid
	}
}

func TestBFTPeersStart1PB(t *testing.T) {
	// then we can run peers
	for i := range bftPeers {
		pid := startPeer(i)
		bftPeers[i] = pid
	}
}

func TestBFTLoadClients1PB(t *testing.T) {
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

func TestBFTCreateChannels1PB(t *testing.T) {
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
				bftPeers[0] = startPeer(0)
			}()
		}

		// client0 create channel
		channel := "test" + strconv.Itoa(m) + "0"
		fmt.Printf("Create channel %s ...\n", channel)
		err := bftClients[0].CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
	}

	// compare tx, one is peer0 starting another is peer0 stopped
	time.Sleep(2 * time.Second)
	require.NoError(t, compareTxs())
}

func TestBFTCallTx1PB(t *testing.T) {
	// client0 create channel test0
	require.NoError(t, createChannelForCallTx())
	// create smart contract on channel test0
	require.NoError(t, createContractForCallTx())

	for m := 1; m <= 6; m++ {
		if m == 2 { // stop peer0
			go func(t *testing.T) {
				fmt.Println("Begin to stop peer 0 ...")
				stopPeer(bftPeers[0])
				require.NoError(t, os.RemoveAll(getBFTPeerDataPath(0)))
			}(t)
		}
		if m == 4 { // restart peer0
			go func() {
				fmt.Println("Begin to restart peer 0 ...")
				bftPeers[0] = startPeer(0)
			}()
		}

		// client0 call setNum fun
		fmt.Printf("Call contract %d times on channel test0 ...\n", m)
		if m%2 == 0 {
			num := "1" + strconv.Itoa(m-1)
			require.NoError(t, getNumForCallTx(0, num))
		} else {
			num := "1" + strconv.Itoa(m)
			require.NoError(t, setNumForCallTx(0, num))
		}
	}
	// compare tx, one is peer0 starting another is peer0 stopped
	time.Sleep(2 * time.Second)
	require.NoError(t, compareTxs())
}

func TestBFTEnd1PB(t *testing.T) {
	for _, pid := range bftOrderers {
		stopOrderer(pid)
	}
	for _, pid := range bftPeers {
		stopPeer(pid)
	}

	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/bft"))
}
