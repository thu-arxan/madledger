// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package cmd

import (
	"fmt"
	"io/ioutil"
	"madledger/common/util"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tc "github.com/tendermint/tendermint/config"
	tlc "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
	tt "github.com/tendermint/tendermint/types/time"
)

var (
	initCmd = &cobra.Command{
		Use: "init",
	}
	initViper = viper.New()
	cfg       string
)

func init() {
	initCmd.RunE = runInit
	initCmd.Flags().StringP("type", "t", "bft", "The consensus type")
	initViper.BindPFlag("type", initCmd.Flags().Lookup("type"))
	initCmd.Flags().StringP("config", "c", "orderer.yaml", "The config file")
	initViper.BindPFlag("config", initCmd.Flags().Lookup("config"))
	initCmd.Flags().StringP("path", "p", "", "The path of orderer")
	initViper.BindPFlag("path", initCmd.Flags().Lookup("path"))
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	var err error
	consensusType := initViper.GetString("type")
	if consensusType == "" {
		consensusType = "bft"
	}
	cfgFile := initViper.GetString("config")
	if cfgFile == "" {
		cfgFile = "orderer.yaml"
	}
	ordererPath := initViper.GetString("path")
	if ordererPath == "" {
		ordererPath = homeDir
	} else {
		if ordererPath, err = util.MakeFileAbs(ordererPath, homeDir); err != nil {
			return err
		}
	}

	var tendermintP2PID string
	if tendermintP2PID, err = initTendermintEnv(ordererPath); err != nil {
		return nil
	}
	if err = createConfigFile(cfgFile, ordererPath, consensusType, tendermintP2PID); err != nil {
		return err
	}
	return nil
}

func createConfigFile(cfgFile, path, consensusType string, tendermintP2PID string) error {
	cfgAbsPath, err := util.MakeFileAbs(cfgFile, homeDir)
	if err != nil {
		return err
	}
	if util.FileExists(cfgAbsPath) {
		return nil
	}
	cfg = cfgTemplate
	blockChainPath, _ := util.MakeFileAbs("data/blocks", path)
	tendermintPath, _ := util.MakeFileAbs(".tendermint", path)
	raftPath, _ := util.MakeFileAbs(".raft", path)
	levelDBPath, _ := util.MakeFileAbs("data/leveldb", path)
	cfg = strings.Replace(cfg, "<<<BlockChainPath>>>", blockChainPath, 1)
	cfg = strings.Replace(cfg, "<<<ConsensusType>>>", consensusType, 1)
	cfg = strings.Replace(cfg, "<<<TendermintPath>>>", tendermintPath, 1)
	cfg = strings.Replace(cfg, "<<<RaftPath>>>", raftPath, 1)
	cfg = strings.Replace(cfg, "<<<LevelDBPath>>>", levelDBPath, 1)
	cfg = strings.Replace(cfg, "<<<TendermintP2PID>>>", tendermintP2PID, 1)

	return ioutil.WriteFile(cfgAbsPath, []byte(cfg), 0755)
}

// initTendermintEnv will create all necessary things that tendermint needs
func initTendermintEnv(path string) (string, error) {
	tendermintPath, _ := util.MakeFileAbs(".tendermint", path)
	os.MkdirAll(tendermintPath+"/config", 0777)
	os.MkdirAll(tendermintPath+"/data", 0777)
	var conf = tc.DefaultConfig()
	privValKeyFile := tendermintPath + "/" + conf.PrivValidatorKeyFile()
	privValStateFile := tendermintPath + "/" + conf.PrivValidatorStateFile()
	var pv *privval.FilePV
	if tlc.FileExists(privValKeyFile) {
		pv = privval.LoadFilePV(privValKeyFile, privValStateFile)
	} else {
		pv = privval.GenFilePV(privValKeyFile, privValStateFile)
		pv.Save()
	}
	nodeKeyFile := tendermintPath + "/" + conf.NodeKeyFile()
	if !tlc.FileExists(nodeKeyFile) {
		if _, err := p2p.LoadOrGenNodeKey(nodeKeyFile); err != nil {
			return "", err
		}
	}

	// genesis file
	genFile := tendermintPath + "/" + conf.GenesisFile()
	if !tlc.FileExists(genFile) {
		genDoc := types.GenesisDoc{
			ChainID:         "madledger",
			GenesisTime:     tt.Now(),
			ConsensusParams: types.DefaultConsensusParams(),
		}
		genDoc.Validators = []types.GenesisValidator{{
			Address: pv.GetPubKey().Address(),
			PubKey:  pv.GetPubKey(),
			Power:   10,
		}}

		if err := genDoc.SaveAs(genFile); err != nil {
			return "", err
		}
	}

	// load node key
	nodeKey, err := p2p.LoadNodeKey(nodeKeyFile)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", nodeKey.ID()), nil
}
