package testfor2clients_raft

import (
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"os"
	cc "madledger/client/config"
	client "madledger/client/lib"
	"time"
)

func TestInitEnv1(t *testing.T) {
	require.NoError(t, initRAFTEnvironment())
}

func TestRaftOrdererStart1(t *testing.T) {
	// then we can run orderers
	for i := range raftOrderers {
		pid := startOrderer(i)
		raftOrderers[i] = pid
	}
}

func TestRaftPeersStart1(t *testing.T) {
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

func TestRaftEnd1(t *testing.T) {
	for _, pid := range raftOrderers {
		stopOrderer(pid)
	}

	for i := range raftPeers {
		raftPeers[i].Stop()
	}
	time.Sleep(2 * time.Second)

	// copy orderers log to other directory
	//require.NoError(t, backupMdFile1("./orderer_tests/"))

	// remove raft data
	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/raft"))
}

