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

func TestInitEnv2(t *testing.T) {
	require.NoError(t, initRAFTEnvironment())
}

func TestRaftOrdererStart2(t *testing.T) {
	// then we can run orderers
	for i := range raftOrderers {
		pid := startOrderer(i)
		raftOrderers[i] = pid
	}
}

func TestRaftPeersStart2(t *testing.T) {
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

func TestLoadClients2(t *testing.T) {
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

func TestRaftCreateChannels2(t *testing.T) {
	// client 0 and client 1 create channels concurrently
	client0 := raftClients[0]
	client1 := raftClients[1]
	var channels []string
	for i := 0; i < 8; i++ {
		if i == 4 {
			fmt.Println("Stop peer 0 ...")
			raftPeers[0].Stop()
			require.NoError(t, os.RemoveAll(getRAFTPeerDataPath(0)))
		}
		if i == 6 {
			require.NoError(t, initPeer(0))

			go func(t *testing.T) {
				fmt.Println("Restart peer 0")
				err := raftPeers[0].Start()
				require.NoError(t, err)
			}(t)
		}
		channel := "test0"
		if i != 0 {
			channel = "test0" + strconv.Itoa(i)
		}
		channels = append(channels, channel)
		fmt.Printf("Create channel %s ...\n", channel)
		err := client0.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)

		// client 1 create channel
		channel = "test1"
		if i != 1 {
			channel = "test1" + strconv.Itoa(i)
		}
		channels = append(channels, channel)
		fmt.Printf("Create channel %s ...\n", channel)
		err = client1.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
	}
	time.Sleep(2 * time.Second)

	// then we will check if channels are create successful
	require.NoError(t, compareChannelName(channels))
	// to avoid block num is not consistent, we should check it
	require.NoError(t, compareChannelBlocks())
}

func TestRaftCreateTx2(t *testing.T) {
	client0 := raftClients[0]
	client1 := raftClients[1]
	for m := 0; m < 8; m++ {
		if m == 3 { // stop peer0
			fmt.Println("Stop peer 0 ...")
			raftPeers[0].Stop()
			require.NoError(t, os.RemoveAll(getRAFTPeerDataPath(0)))
		}
		if m == 6 {
			require.NoError(t, initPeer(0))

			go func(t *testing.T) {
				fmt.Println("Restart peer 0")
				err := raftPeers[0].Start()
				require.NoError(t, err)
			}(t)
		}
		// client 0 create contract
		contractCodes, err := readCodes(getRAFTClientPath(0) + "/MyTest.bin")
		require.NoError(t, err)
		channel := "test0"
		if m != 0 {
			channel = "test0" + strconv.Itoa(m)
		}
		fmt.Printf("Create contract %d on channel %s ...\n", m, channel)
		tx, err := types.NewTx(channel, common.ZeroAddress, contractCodes, client0.GetPrivKey(),types.NORMAL)
		require.NoError(t, err)

		_, err = client0.AddTx(tx)
		require.NoError(t, err)

		// client 1 create contract
		contractCodes, err = readCodes(getRAFTClientPath(1) + "/MyTest.bin")
		require.NoError(t, err)
		channel = "test1"
		if m != 1 {
			channel = "test1" + strconv.Itoa(m)
		}
		fmt.Printf("Create contract %d on channel %s ...\n", m, channel)
		tx, err = types.NewTx(channel, common.ZeroAddress, contractCodes, client1.GetPrivKey(),types.NORMAL)
		require.NoError(t, err)

		_, err = client1.AddTx(tx)
		require.NoError(t, err)
	}
	time.Sleep(2 * time.Second)

	require.NoError(t, compareChannelBlocks())
}

func TestRaftCallTx2(t *testing.T) {

	for m := 1; m <= 8; m++ {
		if m == 4 { // stop peer0
			fmt.Println("Stop peer 0 ...")
			raftPeers[0].Stop()
			require.NoError(t, os.RemoveAll(getRAFTPeerDataPath(0)))
		}
		if m == 6 {
			require.NoError(t, initPeer(0))

			go func(t *testing.T) {
				fmt.Println("Restart peer 0")
				err := raftPeers[0].Start()
				require.NoError(t, err)
			}(t)
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

func TestRaftEnd2(t *testing.T) {
	for _, pid := range raftOrderers {
		stopOrderer(pid)
	}

	for i := range raftPeers {
		raftPeers[i].Stop()
	}
	time.Sleep(2 * time.Second)

	// copy orderers log to other directory
	require.NoError(t, backupMdFile2("./peer_tests/"))

	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/raft"))
}
