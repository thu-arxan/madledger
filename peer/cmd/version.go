package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	version = "Peer version v0.0.1"
)

var (
	versionCmd = &cobra.Command{
		Use: "version",
	}
	versionViper = viper.New()
)

func init() {
	versionCmd.RunE = runVersion
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) error {
	fmt.Println(version)
	return nil
}
