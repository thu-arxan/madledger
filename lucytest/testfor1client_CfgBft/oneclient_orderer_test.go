package testfor1client_CfgBft

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
	clientPath := getBFTClientPath("0")
	cfgPath := getBFTClientConfigPath("0")
	cfg, err := cc.LoadConfig(cfgPath)
	require.NoError(t, err)
	re, _ := regexp.Compile("^.*[.]keystore")
	for i := range cfg.KeyStore.Keys {
		cfg.KeyStore.Keys[i] = clientPath + "/.keystore" + re.ReplaceAllString(cfg.KeyStore.Keys[i], "")
	}
	client, err := client.NewClientFromConfig(cfg)
	require.NoError(t, err)
	bftClient = client
}

func TestBFTLoadAdmin1(t *testing.T) {
	clientPath := getBFTClientPath("admin")
	cfgPath := getBFTClientConfigPath("admin")
	cfg, err := cc.LoadConfig(cfgPath)
	require.NoError(t, err)
	re, _ := regexp.Compile("^.*[.]keystore")
	for i := range cfg.KeyStore.Keys {
		cfg.KeyStore.Keys[i] = clientPath + "/.keystore" + re.ReplaceAllString(cfg.KeyStore.Keys[i], "")
	}
	client, err := client.NewClientFromConfig(cfg)
	require.NoError(t, err)
	bftAdmin = client
}

// create channel and create contract on channel
// add orderer 3 randomly
func TestBFTAddNode1(t *testing.T) {
	order := rand.Intn(3)
	for i := 0; i <= 4; i++ {
		// create channel
		channel := "test" + strconv.Itoa(i)
		err := bftClient.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
		// add orderer 3
		if i == order {
			require.NoError(t, addOrRemoveNode("J0juvEKWlK6FAjw4b/oTMM7EYIH1NpHeNOcZ65tIHP8=", 10, "test0"))
		}
		// create contract
		require.NoError(t, createContractForCallTx(strconv.Itoa(i)))
	}

	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels())
}

// remove orderer 0 and then stop orderer 3
func TestBFTRemoveNode1(t *testing.T) {
	num := 12
	for i := 0; i <= 4; i++ {
		if i == 1 {
			require.NoError(t, addOrRemoveNode("eWdg85+iQWQzasBP8x/wOovhhUVk8yAQefW56OCQ6d4=", 0, "test1"))
		}
		if i == 3 { // 2 orderers left can't success
			fmt.Println("Stop orderer 3, two orderers left can't achieve consensus")
			stopOrderer(bftOrderers[3])
		}
		// call contract,even call getNum, odd call setNum
		if i%2 != 0 {
			fmt.Printf("%d: getNumForCallTx ...\n", i)
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
			fmt.Printf("%d: setNumForCallTx ...\n", i)
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
	// compare channel in differnt orderer
	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannels())
}

func TestBFTEND1(t *testing.T) {
	for _, pid := range bftOrderers {
		stopOrderer(pid)
	}
	for _, pid := range bftPeers {
		stopPeer(pid)
	}

	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/bft"))
}
