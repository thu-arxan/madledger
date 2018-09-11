package cmd

import (
	"errors"
	"io/ioutil"
	cutil "madledger/client/util"
	"madledger/common/util"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	initCmd = &cobra.Command{
		Use: "init",
	}
	initViper = viper.New()
)

func init() {
	initCmd.RunE = runInit
	initCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	initViper.BindPFlag("config", initCmd.Flags().Lookup("config"))
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	var err error
	cfgFile := initViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}
	keyStorePath, _ := util.MakeFileAbs(".keystore", homeDir)
	err = os.MkdirAll(keyStorePath, 0777)
	if err != nil {
		return err
	}
	keyPath, err := cutil.GeneratePrivateKey(keyStorePath)
	if err != nil {
		return err
	}
	var cfg = cfgTemplate
	cfg = strings.Replace(cfg, "<<<KEYFILE>>>", keyPath, 1)

	return ioutil.WriteFile(cfgFile, []byte(cfg), 0777)
}
