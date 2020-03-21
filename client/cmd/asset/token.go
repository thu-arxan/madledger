package asset

import (
	"encoding/json"
	"errors"
	"madledger/blockchain/asset"
	"madledger/client/lib"
	"madledger/client/util"

	coreTypes "madledger/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	tokenCmd = &cobra.Command{
		Use: "token",
	}
	tokenViper = viper.New()
)

func init() {
	tokenCmd.RunE = runToken
	tokenCmd.Flags().StringP("channelID", "n", "", "The name of channel")
	tokenViper.BindPFlag("channelID", tokenCmd.Flags().Lookup("channelID"))

	tokenCmd.Flags().Int64P("value", "v", 0, "The amount of asset to exchange")
	tokenViper.BindPFlag("value", tokenCmd.Flags().Lookup("value"))

	tokenCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	tokenViper.BindPFlag("config", tokenCmd.Flags().Lookup("config"))
}

func runToken(cmd *cobra.Command, args []string) error {
	cfgFile := tokenViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}
	channelID := tokenViper.GetString("channelID")
	if channelID == "" {
		return errors.New("The name of channel should be [a-z0-9]{1,32} such as test, test01 and etc")
	}
	value := uint64(tokenViper.GetInt64("value"))
	if value <= 0 {
		return errors.New("the amount can not be less than or equal to 0")
	}
	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(asset.Payload{
		ChannelID: channelID,
	})
	if err != nil {
		return err
	}
	tx, err := coreTypes.NewTx(coreTypes.ASSETCHANNELID, coreTypes.TokenExchangeAddress, payload, value, "", client.GetPrivKey())
	if err != nil {
		return err
	}

	status, err := client.AddTx(tx)
	table := util.NewTable()
	table.SetHeader("Status", "Error")
	table.AddRow(status, err)
	table.Render()
	return err
}
