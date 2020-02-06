package testfor1client_CfgRaft

import (
	"fmt"
	cc "madledger/client/config"
	client "madledger/client/lib"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// change the package
func TestInitEnv1RC(t *testing.T) {
	require.NoError(t, initRaftEnvironment())
}

func TestRaftOrdererStart1RC(t *testing.T) {
	// then we can run orderers
	for i := range raftOrderers {
		pid := startOrderer(i)
		raftOrderers[i] = pid
	}
}

func TestRaftPeersStart1RC(t *testing.T) {
	// then we can run peers
	for i := range raftPeers {
		pid := startPeer(i)
		raftPeers[i] = pid
	}
}

func TestRaftLoadClients1RC(t *testing.T) {
	time.Sleep(1 * time.Second)
	require.NoError(t, loadClient("0", 0))
	require.NoError(t, loadClient("3", 1))
}

func TestRaftLoadAdmin1RC(t *testing.T) {
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
func TestRaftAddNode1RC(t *testing.T) {
	order := rand.Intn(3)
	for i := 0; i <= 4; i++ {
		// create channel
		channel := "test" + strconv.Itoa(i)
		fmt.Printf("Create channel %s ...\n", channel)
		err := raftClient[0].CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
		// add orderer 3
		if i == order {
			fmt.Println("Add raft node on channel test0 ...")
			require.NoError(t, addNode(4, "127.0.0.1:45680", "test0"))
		}
		// create contract
		fmt.Printf("Create contract on channel %s ...\n", channel)
		require.NoError(t, createContractForCallTx(channel))
	}

	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels())
}

// remove node 1, then stop node 2 and 4, node 3 left
func TestRaftRemoveNode1RC(t *testing.T) {
	num := 12
	for i := 0; i <= 4; i++ {
		if i == 1 {
			fmt.Println("Remove node 1 on channel test1 ...")
			require.NoError(t, removeNode(1, "test1"))
		}
		if i == 3 { // one node left can't success
			fmt.Println("Stop node 2 and 4 on channel test1 ...")
			stopOrderer(raftOrderers[1])
			stopOrderer(raftOrderers[3])
		}
		// call contract,even call getNum, odd call setNum
		if i%2 != 0 {
			fmt.Printf("%d: getNumForCallTx on channel test%d ...\n", i, i)
			if i == 3 {
				go getNumForCallTx(strconv.Itoa(i), strconv.Itoa(num))
				select {
				case <-time.After(5 * time.Second):
					fmt.Println("run too long, execute another tx ...")
				}
			} else {
				require.NoError(t, getNumForCallTx(strconv.Itoa(i), strconv.Itoa(num)))
			}
			num = num + 4
		}
		if i%2 == 0 {
			fmt.Printf("%d: setNumForCallTx on channel test%d...\n", i, i)
			if i == 4 {
				go setNumForCallTx(strconv.Itoa(i), strconv.Itoa(num))
				select {
				case <-time.After(5 * time.Second):
					fmt.Println("run too long, execute another tx ...")
				}
			} else {
				require.NoError(t, setNumForCallTx(strconv.Itoa(i), strconv.Itoa(num)))
			}
		}
	}
}

func TestRaftEND1RC(t *testing.T) {
	for _, pid := range raftOrderers {
		stopOrderer(pid)
	}
	for _, pid := range raftPeers {
		stopPeer(pid)
	}

	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/raft"))
}
