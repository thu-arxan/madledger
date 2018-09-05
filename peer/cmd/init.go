package cmd

import (
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
	initCmd.RunE = runinit
	rootCmd.AddCommand(initCmd)
}

// TODO: fulfill the init
func runinit(cmd *cobra.Command, args []string) error {
	log.Info("init is not finished yet")
	return nil
}
