// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	bc "madledger/blockchain/config"
	cc "madledger/client/config"
	client "madledger/client/lib"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/common/util"
	"madledger/core"
	oc "madledger/orderer/config"
	pc "madledger/peer/config"
	peer "madledger/peer/server"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

/*
* This test will start from a empty environment and start some orderers support bft consensus.
* This test will include operations below.
 */

var (
	// bftOrderers store the pid of orderers
	bftOrderers [4]string
	bftClients  [4]*client.Client
	bftPeers    [4]*peer.Server
	bftChannels []string
)

var peerAddress []string

var (
	bftClientsSet []*core.Member
)

func TestBFT(t *testing.T) {
	require.NoError(t, initBFTEnvironment())
}

func TestBFTRun(t *testing.T) {
	for i := range bftOrderers {
		pid := startOrderer(i)
		bftOrderers[i] = pid
	}
	fmt.Println("We need 3sec to ensure Orderers has been run...")
	time.Sleep(3 * time.Second) // wait until the orderer start
	fmt.Println("Then continue")
}

func TestBFTPeersStart(t *testing.T) {
	peerAddress = make([]string, 4)
	peerAddress[0] = "localhost:20500"
	peerAddress[1] = "localhost:20501"
	peerAddress[2] = "localhost:20502"
	peerAddress[3] = "localhost:20503"
	for i := 0; i < 4; i++ {
		cfgPath := getBFTPeerConfigPath(i)
		cfg, err := pc.LoadConfig(cfgPath)
		require.NoError(t, err)
		server, err := peer.NewServer(cfg)

		require.NoError(t, err)
		bftPeers[i] = server
	}
	for i := range bftPeers {
		go func(t *testing.T, i int) {
			err := bftPeers[i].Start()
			require.NoError(t, err)
		}(t, i)
	}

	time.Sleep(2 * time.Second)
}
func TestBFTLoadClients(t *testing.T) {
	for i := range bftClients {
		clientPath := getBFTClientPath(i)
		cfgPath := getBFTClientConfigPath(i)
		cfg, err := cc.LoadConfig(cfgPath)
		require.NoError(t, err)
		re, _ := regexp.Compile("^.*[.]keystore")
		for i := range cfg.KeyStore.Keys {
			cfg.KeyStore.Keys[i] = clientPath + "/.keystore" + re.ReplaceAllString(cfg.KeyStore.Keys[i], "")
		}
		client, err := client.NewClientFromConfig(cfg)
		require.NoError(t, err)
		bftClients[i] = client

	}
	c0, _ := core.NewMember(bftClients[0].GetPrivKey().PubKey(), "admin")
	c1, _ := core.NewMember(bftClients[1].GetPrivKey().PubKey(), "admin")
	c2, _ := core.NewMember(bftClients[2].GetPrivKey().PubKey(), "admin")
	c3, _ := core.NewMember(bftClients[3].GetPrivKey().PubKey(), "admin")
	bftClientsSet = []*core.Member{c0, c1, c2, c3}
}

func TestBFTCreateChannels(t *testing.T) {
	var wg sync.WaitGroup
	var lock sync.RWMutex
	var channels []string

	for i := range bftClients {
		// each client will create 5 channels
		for m := 0; m < 5; m++ {
			fmt.Printf("BFT Create Channel %d-%d", i, m)
			wg.Add(1)
			go func(t *testing.T, i int) {

				client := bftClients[i]
				channel := strings.ToLower(util.RandomString(16))
				lock.Lock()
				channels = append(channels, channel)
				lock.Unlock()

				err := client.CreateChannel(channel, true, nil, nil, 0, 1, 10000000, peerAddress)
				require.NoError(t, err)
				defer wg.Done()
			}(t, i)
		}
	}
	wg.Wait()
	fmt.Printf("create done\n")
	// then we will check if all channels are create successful
	time.Sleep(2 * time.Second)
	for i := range bftClients {
		wg.Add(1)
		go func(t *testing.T, i int) {
			defer wg.Done()
			client := bftClients[i]
			infos, err := client.ListChannel(false)
			require.NoError(t, err)
			lock.RLock()
			require.Len(t, infos, len(channels))
			for i := range infos {
				require.True(t, util.Contain(channels, infos[i].Name))
			}
			lock.RUnlock()
		}(t, i)
	}
	wg.Wait()
	fmt.Printf("check done\n")
}

func TestBFTOrdererRestart(t *testing.T) {
	stopOrderer(bftOrderers[1])
	os.RemoveAll(getBFTOrdererDataPath(1))
	time.Sleep(2000 * time.Millisecond)
	bftOrderers[1] = startOrderer(1)
	time.Sleep(2000 * time.Millisecond)
	for i := range bftOrderers {
		require.True(t, util.IsDirSame(getBFTOrdererBlockPath(0), getBFTOrdererBlockPath(i)), fmt.Sprintf("Orderer %d is not same with 0", i))
	}
}

func TestBFTReCreateChannels(t *testing.T) {
	// Here we recreate 2 channels
	for i := 0; i < 2; i++ {
		channel := strings.ToLower(util.RandomString(16))

		err := bftClients[i].CreateChannel(channel, true, nil, nil, 0, 1, 10000000, peerAddress)
		require.NoError(t, err)
	}
	time.Sleep(2 * time.Second)
	var prevChannels []string
	// Then we will list channels
	for i := range bftClients {
		infos, err := bftClients[i].ListChannel(false)
		require.NoError(t, err)
		var channels []string
		for i := range infos {
			channels = append(channels, infos[i].Name)
		}
		if len(prevChannels) != 0 {
			require.Equal(t, prevChannels, channels)
		}
		prevChannels = channels
		bftChannels = channels
	}
}
func TestBFTCreateTx(t *testing.T) {
	client0 := bftClients[0]
	client1 := bftClients[1]
	for m := 1; m <= 6; m++ {
		if m == 3 { // stop orderer0
			stopOrderer(bftOrderers[0])
			require.NoError(t, os.RemoveAll(getBFTOrdererDataPath(0)))
		}
		if m == 4 { // restart orderer0
			bftOrderers[0] = startOrderer(0)
		}
		// client 0 create contract
		contractCodes, err := readCodes(getBFTClientPath(1) + "/MyTest.bin")
		require.NoError(t, err)
		channel := bftChannels[util.RandNum(len(bftChannels))]
		fmt.Printf("Create contract on channel %s(m=%d) ...\n", channel, m)
		tx, err := core.NewTx(channel, common.ZeroAddress, contractCodes, 0, "", client0.GetPrivKey())
		require.NoError(t, err)
		_, err = client0.AddTx(tx)
		require.NoError(t, err)

		// client 1 create channel
		contractCodes, err = readCodes(getBFTClientPath(1) + "/MyTest.bin")
		require.NoError(t, err)
		channel = bftChannels[util.RandNum(len(bftChannels))]
		fmt.Printf("Create contract on channel %s(m=%d) ...\n", channel, m)
		tx, err = core.NewTx(channel, common.ZeroAddress, contractCodes, 0, "", client1.GetPrivKey())
		require.NoError(t, err)
		_, err = client1.AddTx(tx)
		require.NoError(t, err)

	}
	time.Sleep(1000 * time.Millisecond)
	for i := range bftOrderers {
		require.True(t, util.IsDirSame(getBFTOrdererBlockPath(0), getBFTOrdererBlockPath(i)), fmt.Sprintf("Orderer %d is not same with 0", i))
	}
}

//func TestBFTNodeAdd(t *testing.T) {
//	// get system admin key
//	// TODO: Hard code in config/genesis.go, seems planning to remove
//	// get pubkey from string by base64 encoding
//	data, err := base64.StdEncoding.DecodeString("BGXcjZ3bhemsoLP4HgBwnQ5gsc8VM91b3y8bW0b6knkWu8x" +
//		"CSKO2qiJXARMHcbtZtvU7Jos2A5kFCD1haJ/hLdg=")
//	require.NoError(t, err)
//	privKey, err := crypto.NewPrivateKey(data, crypto.KeyAlgoSecp256k1)
//	require.NoError(t, err)
//
//	// add one orderer
//	ordererPrivKey, err := crypto.GeneratePrivateKey(crypto.KeyAlgoSecp256k1)
//	require.NoError(t, err)
//	ordererPKBytes, err := ordererPrivKey.PubKey().Bytes()
//	require.NoError(t, err)
//	payload, err := json.Marshal(types.ValidatorUpdate{
//		PubKey: types.PubKey{Data: ordererPKBytes,},
//		Power: 10,
//	})
//	tx, err := core.NewTx(core.CONFIGCHANNELID, core.CfgConsensusAddress, payload, 0, "", privKey)
//	require.NoError(t, err)
//	_, err = bftClients[0].AddTx(tx)
//	require.NoError(t, err)
//
//	// check if it's working
//
//	// add duplicated pubkey should fail?
//
//
//}

func TestBFTAsset(t *testing.T) {
	peers := peerAddress
	client := bftClients[0]
	algo := crypto.KeyAlgoSecp256k1

	issuerKey, err := crypto.GeneratePrivateKey(algo)
	require.NoError(t, err)
	falseIssuerKey, err := crypto.GeneratePrivateKey(algo)
	require.NoError(t, err)
	require.NotEqual(t, issuerKey, falseIssuerKey)
	receiverKey, err := crypto.GeneratePrivateKey(algo)
	require.NoError(t, err)

	issuer, err := issuerKey.PubKey().Address()
	require.NoError(t, err)
	falseIssuer, err := falseIssuerKey.PubKey().Address()
	require.NoError(t, err)
	receiver, err := receiverKey.PubKey().Address()
	require.NoError(t, err)

	err = client.CreateChannel("bftasset", true, nil, nil, 0, 1, 10000000, peers)
	require.NoError(t, err)

	//issue to issuer itself
	coreTx := getAssetChannelTx(core.IssueContractAddress, issuer, "", uint64(10), issuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)

	balance, err := client.GetAccountBalance(issuer)
	require.NoError(t, err)
	require.Equal(t, uint64(10), balance)

	//falseissuer issue fail
	coreTx = getAssetChannelTx(core.IssueContractAddress, falseIssuer, "", uint64(10), falseIssuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)

	balance, err = client.GetAccountBalance(falseIssuer)
	require.NoError(t, err)
	require.Equal(t, uint64(0), balance)

	//test issue to channel
	coreTx = getAssetChannelTx(core.IssueContractAddress, common.ZeroAddress, "bftasset", uint64(10), issuerKey)
	// question, what if test is not created?
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)
	balance, err = client.GetAccountBalance(common.AddressFromChannelID("bftasset"))
	require.NoError(t, err)
	require.Equal(t, uint64(10), balance)

	//test transfer
	coreTx = getAssetChannelTx(core.TransferContractrAddress, receiver, "", uint64(5), issuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)
	balance, err = client.GetAccountBalance(receiver)
	require.NoError(t, err)
	require.Equal(t, uint64(5), balance)

	//test transfer fail
	coreTx = getAssetChannelTx(core.TransferContractrAddress, receiver, "", uint64(5), falseIssuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)
	balance, err = client.GetAccountBalance(receiver)
	require.NoError(t, err)
	require.Equal(t, uint64(5), balance)

	//test transfer to oneself
	coreTx = getAssetChannelTx(core.TransferContractrAddress, receiver, "", uint64(5), receiverKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)
	balance, err = client.GetAccountBalance(receiver)
	require.NoError(t, err)
	require.Equal(t, uint64(5), balance)

	//4.test exchangeToken a.k.a transfer to channel in orderer execution
	coreTx = getAssetChannelTx(core.TokenExchangeAddress, common.ZeroAddress, "bftasset", uint64(5), receiverKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)

	balance, err = client.GetAccountBalance(common.AddressFromChannelID("bftasset"))
	require.NoError(t, err)
	require.Equal(t, uint64(15), balance)

	token, err := client.GetTokenInfo(receiver, []byte("bftasset"))
	require.NoError(t, err)
	require.Equal(t, uint64(5), token)

	//test Block Price
	coreTx, err = core.NewTx("bftasset", common.ZeroAddress, []byte("success"), 0, "", issuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)

	//change BlockPrice of test channel's

	payload, err := json.Marshal(bc.Payload{
		ChannelID: "bftasset",
		Profile: &bc.Profile{
			BlockPrice: 100,
		},
	})
	require.NoError(t, err)
	coreTx, err = core.NewTx(core.CONFIGCHANNELID, common.ZeroAddress, payload, 0, "", issuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)

	//now add tx that cause due
	coreTx, err = core.NewTx("bftasset", common.ZeroAddress, []byte("cause due but pass"), 0, "", issuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)

	// now add multiple txs to ensure that orderers have executed prev tx and stopped receiving tx
	for i := 0; i < 10; i++ {
		coreTx, err = core.NewTx("bftasset", common.ZeroAddress, []byte("multiple tx"), 0, fmt.Sprintln(i), issuerKey)
		_, _ = client.AddTx(coreTx)
	}

	//this one should fail
	coreTx, err = core.NewTx("bftasset", common.ZeroAddress, []byte("fail"), 0, "", issuerKey)
	_, err = client.AddTx(coreTx)
	require.Error(t, err)

	//now issue money to channel account to wake it
	coreTx = getAssetChannelTx(core.IssueContractAddress, common.ZeroAddress, "bftasset", uint64(1000000), issuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)

	coreTx, err = core.NewTx("bftasset", common.ZeroAddress, []byte("success again"), 0, "", issuerKey)
	_, err = client.AddTx(coreTx)
	require.NoError(t, err)
}

func TestBFTEnd(t *testing.T) {
	for i := range bftPeers {
		bftPeers[i].Stop()
	}
	time.Sleep(3 * time.Second)
	for i := range bftOrderers {
		stopOrderer(bftOrderers[i])
	}
	gopath := os.Getenv("GOPATH")
	require.NoError(t, os.RemoveAll(gopath+"/src/madledger/tests/.bft"))
}

// initBFTEnvironment will remove old test folders and copy necessary folders, also it will stop orderers on running
func initBFTEnvironment() error {
	// kill all Orderers
	pids := getProcessPid("orderer")
	fmt.Println("pidof orderer returns ", pids)
	for _, pid := range pids {
		stopOrderer(pid)
	}

	gopath := os.Getenv("GOPATH")
	if err := os.RemoveAll(gopath + "/src/madledger/tests/.bft"); err != nil {
		return err
	}
	if err := copy.Copy(gopath+"/src/madledger/tests/env/bft/orderers", gopath+"/src/madledger/tests/.bft/orderers"); err != nil {
		return err
	}
	if err := copy.Copy(gopath+"/src/madledger/tests/env/bft/clients", gopath+"/src/madledger/tests/.bft/clients"); err != nil {
		return err
	}
	if err := copy.Copy(gopath+"/src/madledger/tests/env/bft/peers", gopath+"/src/madledger/tests/.bft/peers"); err != nil {
		return err
	}

	for i := range bftOrderers {
		if err := absBFTOrdererConfig(i); err != nil {
			return err
		}
	}
	for i := range bftPeers {
		if err := absBFTPeerConfig(i); err != nil {
			return err
		}
	}
	for i := range bftClients {
		if err := absBFTClientConfig(i); err != nil {
			return err
		}
	}

	return nil
}

func absBFTClientConfig(node int) error {
	cfgPath := getBFTClientConfigPath(node)
	// load config
	cfg, err := loadClientConfig(cfgPath)
	if err != nil {
		return err
	}
	// change relative path into absolute path
	cfg.KeyStore.Keys[0] = getBFTClientPath(node) + "/" + cfg.KeyStore.Keys[0]
	cfg.TLS.CA = getBFTClientPath(node) + "/" + cfg.TLS.CA
	cfg.TLS.RawCert = getBFTClientPath(node) + "/" + cfg.TLS.RawCert
	cfg.TLS.Key = getBFTClientPath(node) + "/" + cfg.TLS.Key
	// rewrite peer config
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfgPath, data, os.ModePerm)
}

func loadClientConfig(cfgPath string) (*cc.Config, error) {
	cfgBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	var cfg cc.Config
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func absBFTOrdererConfig(node int) error {
	cfgPath := getBFTOrdererConfigPath(node)
	// load config
	cfg, err := loadOrdererConfig(cfgPath)
	if err != nil {
		return err
	}
	cfg.BlockChain.Path = getBFTOrdererPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Path = getBFTOrdererPath(node) + "/" + cfg.DB.LevelDB.Path
	cfg.Consensus.Tendermint.Path = getBFTOrdererPath(node) + "/" + cfg.Consensus.Tendermint.Path
	cfg.TLS.CA = getBFTOrdererPath(node) + "/" + cfg.TLS.CA
	cfg.TLS.RawCert = getBFTOrdererPath(node) + "/" + cfg.TLS.RawCert
	cfg.TLS.Key = getBFTOrdererPath(node) + "/" + cfg.TLS.Key

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfgPath, data, os.ModePerm)
}

func loadOrdererConfig(cfgPath string) (*oc.Config, error) {
	cfgBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	var cfg oc.Config
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// startOrderer run orderer and return pid
func startOrderer(node int) string {
	before := getProcessPid("orderer")
	go func() {
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("orderer start -c %s", getBFTOrdererConfigPath(node)))
		data, err := cmd.Output()
		if err != nil {
			panic(fmt.Sprintf("Run orderer %d failed because %v, %s\nend", node, err, string(data)))
		}
	}()

	for {
		after := getProcessPid("orderer")
		if len(after) != len(before) {
			for _, pid := range after {
				if !util.Contain(before, pid) {
					return pid
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func getProcessPid(process string) []string {
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("pidof %s", process))
	output, err := cmd.Output()
	if err != nil {
		return nil
	}
	strOutput := strings.TrimSpace(string(output))
	if len(strOutput) == 0 {
		return make([]string, 0)
	} else {
		pids := strings.Split(strOutput, " ")
		return pids
	}
}

// stopOrderer stop an orderer
func stopOrderer(pid string) {
	before := getProcessPid("orderer")
	if len(before) == 0 {
		panic("no orderer running")
	}
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("kill -TERM %s", pid))
	data, err := cmd.Output()
	if err != nil {
		panic(fmt.Sprintf("stop orderer %s failed: %v, %s", pid, err, string(data)))
	}

	i := 0
	for {
		if i > 100 {
			panic(fmt.Sprintf("timeout to wait for orderer %s stop", pid))
		}
		after := getProcessPid("orderer")
		if len(after) < len(before) {
			fmt.Printf("stopped orderer %s after %d attempts\n", pid, i)
			break
		}
		i++
		time.Sleep(100 * time.Millisecond)
	}
}

func getBFTOrdererPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/.bft/orderers/%d", gopath, node)
}

func getBFTOrdererConfigPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/.bft/orderers/%d/orderer.yaml", gopath, node)
}

func getBFTClientPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/.bft/clients/%d", gopath, node)
}

func getBFTClientConfigPath(node int) string {
	return getBFTClientPath(node) + "/client.yaml"
}

func getBFTOrdererDataPath(node int) string {
	return fmt.Sprintf("%s/data", getBFTOrdererPath(node))
}

func getBFTOrdererBlockPath(node int) string {
	return fmt.Sprintf("%s/data/blocks", getBFTOrdererPath(node))
}

func getBFTPeerPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/.bft/peers/%d", gopath, node)
}
func getBFTPeerConfigPath(node int) string {
	gopath := os.Getenv("GOPATH")
	return fmt.Sprintf("%s/src/madledger/tests/.bft/peers/%d/peer.yaml", gopath, node)
}

func absBFTPeerConfig(node int) error {
	cfgPath := getBFTPeerConfigPath(node)
	// load config
	cfg, err := loadPeerConfig(cfgPath)
	if err != nil {
		return err
	}
	// change relative path into absolute path
	cfg.BlockChain.Path = getBFTPeerPath(node) + "/" + cfg.BlockChain.Path
	cfg.DB.LevelDB.Dir = getBFTPeerPath(node) + "/" + cfg.DB.LevelDB.Dir
	cfg.KeyStore.Key = getBFTPeerPath(node) + "/" + cfg.KeyStore.Key
	cfg.TLS.CA = getBFTPeerPath(node) + "/" + cfg.TLS.CA
	cfg.TLS.RawCert = getBFTPeerPath(node) + "/" + cfg.TLS.RawCert
	cfg.TLS.Key = getBFTPeerPath(node) + "/" + cfg.TLS.Key
	// rewrite peer config
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfgPath, data, os.ModePerm)
}

func loadPeerConfig(cfgPath string) (*pc.Config, error) {
	cfgBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	var cfg pc.Config
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
