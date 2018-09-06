package cmd

import (
	"errors"

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

// TODO: fulfill the start
func runStart(cmd *cobra.Command, args []string) error {
	cfgFile := startViper.GetString("config")
	if cfgFile == "" {
		return errors.New("Please provide cfgfile")
	}
	return nil
}
