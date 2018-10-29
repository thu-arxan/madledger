package cmd

import (
	"io/ioutil"
	"madledger/common/util"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	initCmd.Flags().StringP("config", "c", "orderer.yaml", "The config file")
	initViper.BindPFlag("config", initCmd.Flags().Lookup("config"))
	initCmd.Flags().StringP("path", "p", "", "The path of orderer")
	initViper.BindPFlag("path", initCmd.Flags().Lookup("path"))
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	cfgFile := initViper.GetString("config")
	if cfgFile == "" {
		cfgFile = "orderer.yaml"
	}
	ordererPath := initViper.GetString("path")
	if ordererPath == "" {
		ordererPath = homeDir
	}
	err := createConfigFile(cfgFile, ordererPath)
	if err != nil {
		return err
	}
	return nil
}

func createConfigFile(cfgFile, path string) error {
	cfgAbsPath, err := util.MakeFileAbs(cfgFile, homeDir)
	if err != nil {
		return err
	}
	cfg = cfgTemplate
	blockChainPath, _ := util.MakeFileAbs("orderer/data/blocks", path)
	tendermintPath, _ := util.MakeFileAbs("orderer/.tendermint", path)
	levelDBPath, _ := util.MakeFileAbs("orderer/data/leveldb", path)
	cfg = strings.Replace(cfg, "<<<BlockChainPath>>>", blockChainPath, 1)
	cfg = strings.Replace(cfg, "<<<TendermintPath>>>", tendermintPath, 1)
	cfg = strings.Replace(cfg, "<<<LevelDBPath>>>", levelDBPath, 1)

	return ioutil.WriteFile(cfgAbsPath, []byte(cfg), 0755)
}
