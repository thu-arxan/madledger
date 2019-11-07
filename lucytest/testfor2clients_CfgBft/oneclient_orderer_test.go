package testfor2celints_CfgBft

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

func TestBFTLoadClient1(t *testing.T) {
	client, err := loadClient("0")
	if err == nil {
		bftClient[0] = client
	} else {
		require.NoError(t, err)
	}
	client, err = loadClient("1")
	if err == nil {
		bftClient[1] = client
	} else {
		require.NoError(t, err)
	}
}

func TestBFTLoadAdmin1(t *testing.T) {
	client, err := loadClient("admin")
	if err == nil {
		bftAdmin = client
	} else {
		require.NoError(t, err)
	}
}

// create channel and create contract on channel
// add orderer 3 randomly
func TestBFTAddNode1(t *testing.T) {
	for i := 1; i <= 8; i++ {
		// create channel
		channel := "test0" + strconv.Itoa(i)
		err := bftClient[0].CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)

		channel = "test1" + strconv.Itoa(i)
		err = bftClient[1].CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
		// add orderer 3
		if i == 3 {
			go func() {
				err := addOrRemoveNode("J0juvEKWlK6FAjw4b/oTMM7EYIH1NpHeNOcZ65tIHP8=", 10, "test01")
				if err != nil {
					panic(fmt.Sprintf("Add node failed, because %s!", err.Error()))
				}
			}()
		}
	}

	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels(bftClient[0], bftAdmin))
}

// remove orderer 0 and then stop orderer 3
func TestBFTRemoveNode1(t *testing.T) {
	for i := 1; i <= 8; i++ {
		if i < 7 {
			channel := "test0" + strconv.Itoa(i)
			require.NoError(t, createContractForCallTx(channel, "0", bftClient[0]))
			channel = "test1" + strconv.Itoa(i)
			require.NoError(t, createContractForCallTx(channel, "1", bftClient[1]))
		} else {
			channel := "test0" + strconv.Itoa(i)
			go createContractForCallTx(channel, "0", bftClient[0])
			select {
			case <-time.After(5 * time.Second):
				fmt.Println("run too long, execute another tx ...")
			}
			channel = "test1" + strconv.Itoa(i)
			go createContractForCallTx(channel, "1", bftClient[1])
			select {
			case <-time.After(5 * time.Second):
				fmt.Println("run too long, execute another tx ...")
			}

		}

		if i == 1 {
			go func() {
				err := addOrRemoveNode("eWdg85+iQWQzasBP8x/wOovhhUVk8yAQefW56OCQ6d4=", 0, "test01")
				if err != nil {
					panic(fmt.Sprintf("Add node failed, because %s!", err.Error()))
				}
			}()
		}
		if i == 6 {
			fmt.Println("Stop orderer 1, two orderers left can't achieve consensus")
			stopOrderer(bftOrderers[1])
		}
	}
	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels(bftClient[0],bftClient[1]))
}

func TestBFTEND1(t *testing.T) {
	for _, pid := range bftOrderers {
		stopOrderer(pid)
	}

	time.Sleep(2 * time.Second)
	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/bft"))
}
