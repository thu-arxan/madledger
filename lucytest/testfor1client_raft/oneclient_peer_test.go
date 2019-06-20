package testfor1client_raft

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"regexp"
	"strconv"
	"testing"
	cc "madledger/client/config"
	client "madledger/client/lib"
	"time"
)

// change the package
func TestInitEnv2(t *testing.T) {
	require.NoError(t, initRAFTEnvironment())
}

func TestBFTOrdererStart2(t *testing.T) {
	// then we can run orderers
	for i := range raftOrderers {
		pid := startOrderer(i)
		raftOrderers[i] = pid
	}
}

func TestRAFTPeersStart2(t *testing.T) {
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

func TestBFTCreateChannels2(t *testing.T) {
	client := raftClients[0]
	var channels []string
	for m := 1; m <= 8; m++ {
		if m == 4 {
			go func(t *testing.T) {
				fmt.Println("Stop peer 0")
				raftPeers[0].Stop()
				require.NoError(t, os.RemoveAll(getRAFTPeerDataPath(0)))
			}(t)
		}
		if m == 6 {
			require.NoError(t, initPeer(0))

			go func(t *testing.T) {
				fmt.Println("Restart peer 0")
				err := raftPeers[0].Start()
				require.NoError(t, err)
			}(t)
		}
		// client 0 create channel
		channel := "test" + strconv.Itoa(m)
		channels = append(channels, channel)
		fmt.Printf("Create channel %s ...\n", channel)
		err := client.CreateChannel(channel, true, nil, nil)
		require.NoError(t, err)
	}
	time.Sleep(2 * time.Second)

	// then we will check if channels are create successful
	require.NoError(t, compareChannelName(channels))
}

func TestBFTEnd2(t *testing.T) {
	for _, pid := range raftOrderers {
		stopOrderer(pid)
	}

	for i := range raftPeers {
		raftPeers[i].Stop()
	}
	time.Sleep(2 * time.Second)

	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/raft"))
}
