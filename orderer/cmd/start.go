package cmd

import (
	"errors"
	"madledger/orderer/config"
	"madledger/orderer/server"
	"madledger/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	startCmd = &cobra.Command{
		Use: "start",
	}
	startViper = viper.New()
)

func init() {
	startCmd.RunE = runStart
	startCmd.Flags().StringP("config", "c", "orderer.yaml", "The config file of blockchain")
	startViper.BindPFlag("config", startCmd.Flags().Lookup("config"))
	rootCmd.AddCommand(startCmd)
}

// TODO: fulfill the start
func runStart(cmd *cobra.Command, args []string) error {
	cfgFile := startViper.GetString("config")
	if cfgFile == "" {
		return errors.New("Please provide the config file")
	}
	cfgAbsPath, err := util.MakeFileAbs(cfgFile, homeDir)
	if err != nil {
		return err
	}
	cfg, err := config.LoadConfig(cfgAbsPath)
	if err != nil {
		return err
	}
	// set the log
	setLog(cfg.Debug)
	serverCfg, err := cfg.GetServerConfig()
	if err != nil {
		return err
	}
	dbCfg, err := cfg.GetDBConfig()
	if err != nil {
		return err
	}
	s, err := server.NewServer(serverCfg, dbCfg)
	if err != nil {
		return err
	}
	go registerStop(s)
	err = s.Start()
	if err != nil {
		return err
	}
	return nil
}
