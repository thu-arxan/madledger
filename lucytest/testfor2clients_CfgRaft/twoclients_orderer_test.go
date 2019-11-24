package testfor2clients_CfgRaft

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

// change the package
func TestInitEnv2RC(t *testing.T) {
	require.NoError(t, initRaftEnvironment())
}

func TestRaftOrdererStart2RC(t *testing.T) {
	// then we can run orderers
	for i := range raftOrderers {
		pid := startOrderer(i)
		raftOrderers[i] = pid
	}
}

func TestRaftPeersStart12RC(t *testing.T) {
	// then we can run peers
	for i := range raftPeers {
		pid := startPeer(i)
		raftPeers[i] = pid
	}
}

func TestRaftLoadClients2RC(t *testing.T) {
	time.Sleep(1*time.Second)
	require.NoError(t, loadClient("0",0))
	require.NoError(t,loadClient("1",1))
	require.NoError(t,loadClient("3",2))
}

func TestRaftLoadAdmin2RC(t *testing.T) {
	clientPath := getRaftClientPath("admin")
	cfgPath := getRaftClientConfigPath("admin")
	cfg, err := cc.LoadConfig(cfgPath)
	require.NoError(t, err)
	re, _ := regexp.Compile("^.*[.]keystore")
	for i := range cfg.KeyStore.Keys {
		cfg.KeyStore.Keys[i] = clientPath + "/.keystore" + re.ReplaceAllString(cfg.KeyStore.Keys[i], "")
	}
	client, err := client.NewClientFromConfig(cfg)
	require.NoError(t, err)
	raftAdmin = client
}

// create channel and create contract on channel
// add orderer 3 randomly
func TestRaftAddNode2RC(t *testing.T) {
	for i := 1; i <= 8; i++ {
		// create channel
		channel := "test0" + strconv.Itoa(i)
		err := raftClient[0].CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
		channel = "test1" + strconv.Itoa(i)
		err = raftClient[1].CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)

		// add orderer 3
		if i == 3 {
			go func() {
				err := addNode(4, "127.0.0.1:45680", "test01")
				if err != nil {
					panic(fmt.Sprintf("Add node failed, because %s!", err.Error()))
				}
			}()
		}

	}

	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels(raftClient[0], raftClient[2]))
}

// remove node 1, then stop node 2 and 4, node 3 left
func TestRaftRemoveNode2RC(t *testing.T) {
	for i := 1; i <= 8; i++ {
		if i < 7 {
			channel := "test0" + strconv.Itoa(i)
			require.NoError(t, createContractForCallTx(channel, "0", raftClient[0]))
			channel = "test1" + strconv.Itoa(i)
			require.NoError(t, createContractForCallTx(channel, "1", raftClient[1]))
		} else {
			channel := "test0" + strconv.Itoa(i)
			go createContractForCallTx(channel, "0", raftClient[0])
			select {
			case <-time.After(5 * time.Second):
				fmt.Println("run too long, execute another tx ...")
			}
			channel = "test1" + strconv.Itoa(i)
			go createContractForCallTx(channel, "1", raftClient[1])
			select {
			case <-time.After(5 * time.Second):
				fmt.Println("run too long, execute another tx ...")
			}
		}
		if i == 1 {
			go func() {
				err := removeNode(1, "test01")
				if err != nil {
					panic(fmt.Sprintf("Remove node failed, because %s!", err.Error()))
				}
			}()
		}
		if i == 6 {
			fmt.Println("Stop orderer 1 and orderer 3, one orerder left can't achieve consensus")
			stopOrderer(raftOrderers[1])
			stopOrderer(raftOrderers[3])
		}
	}
	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels(raftClient[0], raftClient[2]))
}

func TestRaftEND2RC(t *testing.T) {
	for _, pid := range raftOrderers {
		stopOrderer(pid)
	}
	for _, pid := range raftPeers {
		stopPeer(pid)
	}

	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/raft"))
}
