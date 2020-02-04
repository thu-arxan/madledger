package testfor1client_raft

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

func TestInitEnv1OR(t *testing.T) {
	require.NoError(t, initRAFTEnvironment())
}

func TestRaftOrdererStart1OR(t *testing.T) {
	// then we can run orderers
	for i := range raftOrderers {
		pid := startOrderer(i)
		raftOrderers[i] = pid
	}
}

func TestRaftPeersStart1OR(t *testing.T) {
	// then we can run peers
	for i := range raftPeers {
		pid := startPeer(i)
		raftPeers[i] = pid
	}
}

func TestLoadClients1OR(t *testing.T) {
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

func TestRaftCreateChannels1OR(t *testing.T) {
	// client0 create 8 channels
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
		// client0 create channel
		channel := "test" + strconv.Itoa(i)
		fmt.Printf("Create channel %s ...\n", channel)
		err := raftClients[0].CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
	}
	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels())
}

func TestRaftCreateTx1OR(t *testing.T) {
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
		channel := "test" + strconv.Itoa(m)
		fmt.Printf("Create contract %d on channel %s ...\n", m, channel)
		tx, err := types.NewTx(channel, common.ZeroAddress, contractCodes, raftClients[0].GetPrivKey())
		require.NoError(t, err)

		_, err = raftClients[0].AddTx(tx)
		require.NoError(t, err)
	}
	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels())
}

func TestRaftCallTx1OR(t *testing.T) {
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
			require.NoError(t, getNumForCallTx(num))
		} else {
			num := "1" + strconv.Itoa(i)
			require.NoError(t, setNumForCallTx(num))
		}
	}
	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels())
}

func TestBFTEnd1OR(t *testing.T) {
	for _, pid := range raftOrderers {
		stopOrderer(pid)
	}
	for _, pid := range raftPeers {
		stopPeer(pid)
	}

	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/raft"))
}