package cmd

import (
	"io/ioutil"
	"madledger/common/util"

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
	initCmd.RunE = runinit
	initCmd.Flags().StringP("config", "c", "orderer.yaml", "The config file")
	initViper.BindPFlag("config", initCmd.Flags().Lookup("config"))
	rootCmd.AddCommand(initCmd)
}

func runinit(cmd *cobra.Command, args []string) error {
	cfgFile := initViper.GetString("config")
	if cfgFile == "" {
		cfgFile = "orderer.yaml"
	}
	err := createConfigFile(cfgFile)
	if err != nil {
		return err
	}
	return nil
}

func createConfigFile(cfgFile string) error {
	cfgAbsPath, err := util.MakeFileAbs(cfgFile, homeDir)
	if err != nil {
		return err
	}
	cfg = cfgTemplate
	return ioutil.WriteFile(cfgAbsPath, []byte(cfg), 0755)
}
