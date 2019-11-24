package testfor2clients_raft

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

func TestInitEnv2PR(t *testing.T) {
	require.NoError(t, initRAFTEnvironment())
}

func TestRaftOrdererStart2PR(t *testing.T) {
	// then we can run orderers
	for i := range raftOrderers {
		pid := startOrderer(i)
		raftOrderers[i] = pid
	}
}

func TestRaftPeersStart2PR(t *testing.T) {
	// then we can run peers
	for i := range raftPeers {
		pid := startPeer(i)
		raftPeers[i] = pid
	}
}

func TestLoadClients2PR(t *testing.T) {
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

func TestRaftCreateChannels2PR(t *testing.T) {
	// client 0 and client 1 create channels concurrently
	for i := 0; i < 8; i++ {
		if i == 2 {
			go func(t *testing.T) {
				fmt.Println("Stop peer 0 ...")
				stopPeer(raftPeers[0])
				require.NoError(t, os.RemoveAll(getRAFTPeerDataPath(0)))
			}(t)
		}
		if i == 5 {
			go func() {
				fmt.Println("Restart peer 0 ...")
				raftPeers[0]=startPeer(0)
			}()
		}
		channel := "test0"
		if i != 0 {
			channel = "test0" + strconv.Itoa(i)
		}
		fmt.Printf("Create channel %s by client0 ...\n", channel)
		err := raftClients[0].CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)

		// client 1 create channel
		channel = "test1"
		if i != 1 {
			channel = "test1" + strconv.Itoa(i)
		}
		fmt.Printf("Create channel %s by client1 ...\n", channel)
		err = raftClients[1].CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
	}
	// compare tx, one is peer0 starting another is peer0 stopped
	time.Sleep(2 * time.Second)
	require.NoError(t, compareTxs())
}

func TestRaftCreateTx2PR(t *testing.T) {
	for m := 0; m < 8; m++ {
		if m == 3 { // stop peer0
			go func(t *testing.T) {
				fmt.Println("Stop peer 0 ...")
				stopPeer(raftPeers[0])
				require.NoError(t, os.RemoveAll(getRAFTPeerDataPath(0)))
			}(t)
		}
		if m == 6 {
			go func() {
				fmt.Println("Restart peer 0 ...")
				raftPeers[0]=startPeer(0)
			}()
		}
		// client 0 create contract
		contractCodes, err := readCodes(getRAFTClientPath(0) + "/MyTest.bin")
		require.NoError(t, err)
		channel := "test0"
		if m != 0 {
			channel = "test0" + strconv.Itoa(m)
		}
		fmt.Printf("Create contract %d on channel %s ...\n", m, channel)
		tx, err := types.NewTx(channel, common.ZeroAddress, contractCodes, raftClients[0].GetPrivKey(),types.NORMAL)
		require.NoError(t, err)

		_, err = raftClients[0].AddTx(tx)
		require.NoError(t, err)

		// client 1 create contract
		contractCodes, err = readCodes(getRAFTClientPath(1) + "/MyTest.bin")
		require.NoError(t, err)
		channel = "test1"
		if m != 1 {
			channel = "test1" + strconv.Itoa(m)
		}
		fmt.Printf("Create contract %d on channel %s ...\n", m, channel)
		tx, err = types.NewTx(channel, common.ZeroAddress, contractCodes, raftClients[1].GetPrivKey(),types.NORMAL)
		require.NoError(t, err)

		_, err = raftClients[1].AddTx(tx)
		require.NoError(t, err)
	}
	// compare tx, one is peer0 starting another is peer0 stopped
	time.Sleep(2 * time.Second)
	require.NoError(t, compareTxs())
}

func TestRaftCallTx2PR(t *testing.T) {
	for m := 1; m <= 8; m++ {
		if m == 4 { // stop peer0
			go func(t *testing.T) {
				fmt.Println("Stop peer 0 ...")
				stopPeer(raftPeers[0])
				require.NoError(t, os.RemoveAll(getRAFTPeerDataPath(0)))
			}(t)
		}
		if m == 6 {
			go func() {
				fmt.Println("Restart peer 0 ...")
				raftPeers[0]=startPeer(0)
			}()
		}

		// odd call setNum, even call GetNum
		fmt.Printf("Call contract %d times on channel test0 ...\n", m)
		if m%2 == 0 {
			num := "1" + strconv.Itoa(m-1)
			require.NoError(t, getNumForCallTx(0, num))
		} else {
			num := "1" + strconv.Itoa(m)
			require.NoError(t, setNumForCallTx(0, num))
		}

		// odd call setNum, even call GetNum
		fmt.Printf("Call contract %d times on channel test1 ...\n", m)
		if m%2 == 0 {
			num := "2" + strconv.Itoa(m-1)
			require.NoError(t, getNumForCallTx(1, num))
		} else {
			num := "2" + strconv.Itoa(m)
			require.NoError(t, setNumForCallTx(1, num))
		}
	}
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannelBlocks())
}

func TestRaftEnd2PR(t *testing.T) {
	for _, pid := range raftOrderers {
		stopOrderer(pid)
	}
	for _, pid := range raftPeers {
		stopPeer(pid)
	}

	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/raft"))
}
