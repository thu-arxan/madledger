package cmd

import (
	"errors"
	"madledger/common/util"
	"madledger/peer/config"
	"madledger/peer/server"

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
	startCmd.Flags().StringP("config", "c", "peer.yaml", "The config file of blockchain")
	startViper.BindPFlag("config", startCmd.Flags().Lookup("config"))
	rootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) error {
	cfgFile := startViper.GetString("config")
	if cfgFile == "" {
		return errors.New("Please provide cfgfile")
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
	s, err := server.NewServer(cfg)
	if err != nil {
		return err
	}
	err = s.Start()
	if err != nil {
		return err
	}
	var finish = make(chan bool, 1)
	go registerStop(s, finish)
	<-finish

	return nil
}
