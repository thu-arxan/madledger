package account

import (
	"encoding/json"
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"madledger/blockchain/asset"
	"madledger/client/lib"
	"madledger/common"
	coreTypes "madledger/core"
)

var (
	issueCmd = &cobra.Command{
		Use: "issue",
	}
	issueViper = viper.New()
)

func init() {
	issueCmd.RunE = runIssue
	issueCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	issueViper.BindPFlag("config", issueCmd.Flags().Lookup("config"))
	issueCmd.Flags().StringP("channelID", "n", "", "The channelID of the tx")
	issueViper.BindPFlag("channelID", issueCmd.Flags().Lookup("channelID"))
	issueCmd.Flags().StringP("value", "v", "0",
		"value to be issued")
	issueViper.BindPFlag("value", issueCmd.Flags().Lookup("value"))
	issueCmd.Flags().StringP("address", "a", "0",
		"hex address of the account issued")
	issueViper.BindPFlag("address", issueCmd.Flags().Lookup("address"))
}

func runIssue(cmd *cobra.Command, args []string) error {
	cfgFile := issueViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}

	channelID := issueViper.GetString("channelID")
	if channelID == "" {
		return errors.New("The channelID of tx can not be nil")
	}

	value := issueViper.GetInt("value")
	if value < 0 {
		return errors.New("cannot issue negative value")
	}

	receiver := issueViper.GetString("address")
	recipient := common.HexToAddress(receiver)

	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(asset.Payload{
		Action: "issue",
		ChannelID: channelID,
	})
	if err != nil {
		return err
	}

	tx, err := coreTypes.NewTx(coreTypes.ASSETCHANNELID, recipient, payload, value, "", client.GetPrivKey())
	if err != nil {
		return err
	}

	_, err = client.AddTx(tx)

	return err
}
