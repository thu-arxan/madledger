package testfor2clients_CfgRaft

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
	"time"
)

// change the package
func TestInitEnv1(t *testing.T) {
	require.NoError(t, initRaftEnvironment())
}

func TestRaftOrdererStart1(t *testing.T) {
	// then we can run orderers
	for i := range raftOrderers {
		pid := startOrderer(i)
		raftOrderers[i] = pid
	}
}

func TestRaftPeersStart1(t *testing.T) {
	// then we can run peers
	for i := range raftPeers {
		pid := startPeer(i)
		raftPeers[i] = pid
	}
}

func TestRaftLoadClients1(t *testing.T) {
	client, err := loadClient("0")
	require.NoError(t, err)
	raftClient[0] = client
	client, err = loadClient("1")
	require.NoError(t, err)
	raftClient[1] = client
	client, err = loadClient("3")
	require.NoError(t, err)
	raftClient[2] = client
}

func TestRaftLoadAdmin1(t *testing.T) {
	client, err := loadClient("admin")
	require.NoError(t, err)
	raftAdmin = client
}

// create channel and create contract on channel
// add orderer 3 randomly
func TestRaftAddNode1(t *testing.T) {
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
func TestRaftRemoveNode1(t *testing.T) {
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
}

func TestRaftEND1(t *testing.T) {
	for _, pid := range raftOrderers {
		stopOrderer(pid)
	}
	for _, pid := range raftPeers {
		stopPeer(pid)
	}

	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/raft"))
}
