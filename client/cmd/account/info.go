package account

import (
	"errors"
	"fmt"
	"madledger/client/lib"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	infoCmd = &cobra.Command{
		Use: "info",
	}
	infoViper = viper.New()
)

type info struct {
	Address string
}

func init() {
	infoCmd.RunE = runInfo
	infoCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	infoViper.BindPFlag("config", infoCmd.Flags().Lookup("config"))
}

func runInfo(cmd *cobra.Command, args []string) error {
	cfgFile := infoViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}
	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}

	var infos []info
	address, err := client.GetPrivKey().PubKey().Address()
	if err != nil {
		return err
	}

	infos = append(infos, info{
		Address: address.String(),
	})

	if len(infos) == 0 {
		fmt.Println("No results!")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Address"})
		table.Append([]string{address.String()})
		table.Render()
	}
	return nil
}
