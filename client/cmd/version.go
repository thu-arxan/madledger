package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	version = "Client version v0.0.1"
)

var (
	versionCmd = &cobra.Command{
		Use: "version",
	}
	versionViper = viper.New()
)

func init() {
	versionCmd.RunE = runversion
	rootCmd.AddCommand(versionCmd)
}

func runversion(cmd *cobra.Command, args []string) error {
	fmt.Println(version)
	return nil
}
