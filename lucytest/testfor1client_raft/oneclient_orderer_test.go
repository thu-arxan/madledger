package testfor1client_raft

import (
	"fmt"
	"github.com/stretchr/testify/require"
	cc "madledger/client/config"
	client "madledger/client/lib"
	"os"
	"madledger/common"
	"madledger/core/types"
	"regexp"
	"strconv"
	"testing"
	"time"
)

// change the package
func TestInitEnv1(t *testing.T) {
	require.NoError(t, initRAFTEnvironment())
}

func TestBFTOrdererStart1(t *testing.T) {
	// then we can run orderers
	for i := range raftOrderers {
		pid := startOrderer(i)
		raftOrderers[i] = pid
	}
}

func TestRAFTPeersStart1(t *testing.T) {
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

func TestLoadClients1(t *testing.T) {
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

func TestRaftCreateChannels1(t *testing.T) {
	// client-0 create 4 channels
	client := raftClients[0]
	var channels []string
	for i := 0; i < 8; i++ {
		if i == 4 {
			fmt.Println("Stop Orderer 0 ...")
			stopOrderer(raftOrderers[0])
			// restart orderer0 by RestartNode
			require.NoError(t, os.RemoveAll(getRAFTOrdererDataPath(0)))
			_, err := copyFile("./0.md", "./00.md")
			require.NoError(t, err)
		}
		if i == 6 {
			fmt.Println("Restart Orderer 0 ...")
			raftOrderers[0] = startOrderer(0)
		}
		// client 0 create channel
		channel := "test" + strconv.Itoa(i)
		channels = append(channels, channel)
		fmt.Printf("Create channel %s ...\n", channel)
		err := client.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
	}
	time.Sleep(2 * time.Second)

	// then we will check if channels are create successful
	require.NoError(t, compareChannelName(channels))
}

func TestBFTCreateTx1(t *testing.T) {
	client := raftClients[0]
	for m := 1; m <= 8; m++ {
		if m == 4 {
			fmt.Println("Stop Orderer 0 ...")
			stopOrderer(raftOrderers[0])
			// restart orderer0 by RestartNode
			require.NoError(t, os.RemoveAll(getRAFTOrdererDataPath(0)))
			_, err := copyFile("./0.md", "./000.md")
			require.NoError(t, err)
		}
		if m == 6 {
			fmt.Println("Restart Orderer 0 ...")
			raftOrderers[0] = startOrderer(0)
		}
		// client 0 create contract
		contractCodes, err := readCodes(getRAFTClientPath(0) + "/MyTest.bin")
		require.NoError(t, err)
		fmt.Printf("Create contract %d on channel test0 ...\n", m)
		tx, err := types.NewTx("test0", common.ZeroAddress, contractCodes, client.GetPrivKey())
		require.NoError(t, err)

		_, err = client.AddTx(tx)
		require.NoError(t, err)
	}

	time.Sleep(2 * time.Second)
	require.NoError(t, compareChannelBlocks())
}

func TestRaftCallTx1(t *testing.T) {
	for i := 1; i <= 8; i++ {
		if i == 4 {
			fmt.Println("Stop Orderer 0 ...")
			stopOrderer(raftOrderers[0])
			// restart orderer0 by RestartNode
			require.NoError(t, os.RemoveAll(getRAFTOrdererDataPath(0)))
			_, err := copyFile("./0.md", "./0000.md")
			require.NoError(t, err)
		}
		if i == 6 {
			fmt.Println("Restart Orderer 0 ...")
			raftOrderers[0] = startOrderer(0)
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
	time.Sleep(2 * time.Second)

	require.NoError(t, compareChannelBlocks())
}

func TestBFTEnd1(t *testing.T) {
	for _, pid := range raftOrderers {
		stopOrderer(pid)
	}

	for i := range raftPeers {
		raftPeers[i].Stop()
	}
	time.Sleep(2 * time.Second)

	// copy orderers log to other directory
	_, err := copyFile("./0.md", "./orderer_tests/0.md")
	require.NoError(t, err)
	_, err = copyFile("./00.md", "./orderer_tests/00.md")
	require.NoError(t, err)
	_, err = copyFile("./1.md", "./orderer_tests/1.md")
	require.NoError(t, err)
	_, err = copyFile("./000.md", "./orderer_tests/000.md")
	require.NoError(t, err)
	_, err = copyFile("./2.md", "./orderer_tests/2.md")
	require.NoError(t, err)
	_, err = copyFile("./0000.md", "./orderer_tests/0000.md")
	require.NoError(t, err)
	_, err = copyFile("./3.md", "./orderer_tests/3.md")
	require.NoError(t, err)

	// remove raft data
	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/raft"))
}