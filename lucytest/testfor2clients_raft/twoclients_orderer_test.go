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

func TestInitEnv2BR(t *testing.T) {
	require.NoError(t, initRAFTEnvironment())
}

func TestRaftOrdererStart2BR(t *testing.T) {
	// then we can run orderers
	for i := range raftOrderers {
		pid := startOrderer(i)
		raftOrderers[i] = pid
	}
}

func TestRaftPeersStart2BR(t *testing.T) {
	// then we can run peers
	for i := range raftPeers {
		pid := startPeer(i)
		raftPeers[i] = pid
	}
}

func TestLoadClients2BR(t *testing.T) {
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

func TestRaftCreateChannels2BR(t *testing.T) {
	for i := 0; i < 8; i++ {
		if i == 2 {
			go func(t *testing.T) {
				fmt.Println("Stop Orderer 0 ...")
				stopOrderer(raftOrderers[0])
				require.NoError(t, os.RemoveAll(getRAFTOrdererDataPath(0)))
			}(t)
		}
		if i == 5 {
			go func() {
				fmt.Println("Restart Orderer 0 ...")
				raftOrderers[0] = startOrderer(0)
			}()
		}
		// client 0 create channel
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
	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels())
}

func TestRaftCreateTx2BR(t *testing.T) {
	for m := 0; m < 8; m++ {
		if m == 2 {
			go func(t *testing.T) {
				fmt.Println("Stop Orderer 0 ...")
				stopOrderer(raftOrderers[0])
				require.NoError(t, os.RemoveAll(getRAFTOrdererDataPath(0)))
			}(t)
		}
		if m == 5 {
			go func() {
				fmt.Println("Restart Orderer 0 ...")
				raftOrderers[0] = startOrderer(0)
			}()
		}
		// client 0 create contract
		contractCodes, err := readCodes(getRAFTClientPath(0) + "/MyTest.bin")
		require.NoError(t, err)
		channel := "test0"
		if m != 0 {
			channel = "test0" + strconv.Itoa(m)
		}
		fmt.Printf("Create contract %d on channel %s by client0...\n", m, channel)
		tx, err := types.NewTx(channel, common.ZeroAddress, contractCodes, raftClients[0].GetPrivKey())
		require.NoError(t, err)

		_, err = raftClients[0].AddTx(tx)
		require.NoError(t, err)

		// client 0 create contract
		contractCodes, err = readCodes(getRAFTClientPath(1) + "/MyTest.bin")
		require.NoError(t, err)
		channel = "test1"
		if m != 1 {
			channel = "test1" + strconv.Itoa(m)
		}
		fmt.Printf("Create contract %d on channel %s by client1 ...\n", m, channel)
		tx, err = types.NewTx(channel, common.ZeroAddress, contractCodes, raftClients[1].GetPrivKey())
		require.NoError(t, err)

		_, err = raftClients[1].AddTx(tx)
		require.NoError(t, err)
	}
	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels())
}

func TestRaftCallTx12BR(t *testing.T) {
	for i := 1; i <= 8; i++ {
		if i == 3 {
			go func(t *testing.T) {
				fmt.Println("Stop Orderer 0 ...")
				stopOrderer(raftOrderers[0])
				require.NoError(t, os.RemoveAll(getRAFTOrdererDataPath(0)))
			}(t)
		}
		if i == 6 {
			go func() {
				fmt.Println("Restart Orderer 0 ...")
				raftOrderers[0] = startOrderer(0)
			}()
		}

		// odd call setNum, even call GetNum
		fmt.Printf("Call contract %d times on channel test0 ...\n", i)
		if i%2 == 0 {
			num := "1" + strconv.Itoa(i-1)
			require.NoError(t, getNumForCallTx(0, num))
		} else {
			num := "1" + strconv.Itoa(i)
			require.NoError(t, setNumForCallTx(0, num))
		}

		// odd call setNum, even call GetNum
		fmt.Printf("Call contract %d times on channel test1 ...\n", i)
		if i%2 == 0 {
			num := "2" + strconv.Itoa(i-1)
			require.NoError(t, getNumForCallTx(1, num))
		} else {
			num := "2" + strconv.Itoa(i)
			require.NoError(t, setNumForCallTx(1, num))
		}
	}
	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels())
}
func TestRaftEnd2BR(t *testing.T) {
	for _, pid := range raftOrderers {
		stopOrderer(pid)
	}
	for _, pid := range raftPeers {
		stopPeer(pid)
	}

	// remove raft data
	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/raft"))
}
