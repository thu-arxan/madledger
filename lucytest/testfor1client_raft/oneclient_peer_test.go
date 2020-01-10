package testfor1client_raft

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

func TestInitEnv1PR(t *testing.T) {
	require.NoError(t, initRAFTEnvironment())
}

func TestRaftOrdererStart1PR(t *testing.T) {
	// then we can run orderers
	for i := range raftOrderers {
		pid := startOrderer(i)
		raftOrderers[i] = pid
	}
}

func TestRaftPeersStart1PR(t *testing.T) {
	// then we can run peers
	for i := range raftPeers {
		pid := startPeer(i)
		raftPeers[i] = pid
	}
}

func TestLoadClients1PR(t *testing.T) {
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

func TestRaftCreateChannels1PR(t *testing.T) {
	client := raftClients[0]
	for m := 0; m < 8; m++ {
		if m == 2 {
			go func(t *testing.T) {
				fmt.Println("Stop peer 0 ...")
				stopPeer(raftPeers[0])
				require.NoError(t, os.RemoveAll(getRAFTPeerDataPath(0)))
			}(t)
		}
		if m == 6 {
			go func(t *testing.T) {
				fmt.Println("Restart peer 0 ..")
				raftPeers[0] = startPeer(0)
			}(t)
		}
		// client0 create channel
		channel := "test" + strconv.Itoa(m)
		fmt.Printf("Create channel %s ...\n", channel)
		err := client.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
	}
	// compare tx, one is peer0 starting another is peer0 stopped
	time.Sleep(2 * time.Second)
	require.NoError(t, compareTxs())
}

func TestRaftCreateTx1PR(t *testing.T) {
	for m := 0; m < 8; m++ {
		if m == 2 { // stop peer0
			go func(t *testing.T) {
				fmt.Println("Begin to stop peer 0 ...")
				stopPeer(raftPeers[0])
				require.NoError(t, os.RemoveAll(getRAFTPeerDataPath(0)))
			}(t)
		}
		if m == 5 { // restart peer0
			go func() {
				fmt.Println("Begin to restart peer 0 ...")
				raftPeers[0] = startPeer(0)
			}()
		}
		// client 0 create contract
		contractCodes, err := readCodes(getRAFTClientPath(0) + "/MyTest.bin")
		require.NoError(t, err)
		channel := "test" + strconv.Itoa(m)
		fmt.Printf("Create contract %d on channel %s ...\n", m, channel)
		tx, err := types.NewTx(channel, common.ZeroAddress, contractCodes, 0, "", raftClients[0].GetPrivKey())
		require.NoError(t, err)

		_, err = raftClients[0].AddTx(tx)
		require.NoError(t, err)
	}
	// compare tx, one is peer0 starting another is peer0 stopped
	time.Sleep(2 * time.Second)
	require.NoError(t, compareTxs())
}

func TestBFTCallTx1PR(t *testing.T) {
	for m := 1; m <= 8; m++ {
		if m == 3 {
			go func(t *testing.T) {
				fmt.Println("Begin to stop peer 0 ...")
				stopPeer(raftPeers[0])
				require.NoError(t, os.RemoveAll(getRAFTPeerDataPath(0)))
			}(t)
		}
		if m == 6 { // restart peer0
			go func() {
				fmt.Println("Begin to restart peer 0 ...")
				raftPeers[0] = startPeer(0)
			}()
		}

		// client0 call setNum and getNum function in smart contract
		fmt.Printf("Call contract %d times on channel test0 ...\n", m)
		if m%2 == 0 {
			num := "1" + strconv.Itoa(m-1)
			require.NoError(t, getNumForCallTx(num))
		} else {
			num := "1" + strconv.Itoa(m)
			require.NoError(t, setNumForCallTx(num))
		}
	}
	// compare tx, one is peer0 starting another is peer0 stopped
	time.Sleep(2 * time.Second)
	require.NoError(t, compareTxs())
}

func TestRaftEnd1PR(t *testing.T) {
	for _, pid := range raftOrderers {
		stopOrderer(pid)
	}
	for _, pid := range raftPeers {
		stopPeer(pid)
	}

	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/raft"))
}
