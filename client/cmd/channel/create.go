package channel

import (
	"errors"
	"madledger/client/orderer"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	createCmd = &cobra.Command{
		Use: "create",
	}
	createViper = viper.New()
)

func init() {
	createCmd.RunE = runCreate
	createCmd.Flags().StringP("name", "n", "", "The name of channel")
	createViper.BindPFlag("name", createCmd.Flags().Lookup("name"))
}

func runCreate(cmd *cobra.Command, args []string) error {
	name := createViper.GetString("name")
	if name == "" {
		return errors.New("The name of channel should be [a-z0-9]{1,32} such as test, test01 and etc")
	}
	client, err := orderer.NewClient()
	if err != nil {
		return err
	}
	return client.CreateChannel(name)
}
