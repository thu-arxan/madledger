package asset

import (
	"encoding/json"
	"errors"
	"madledger/blockchain/asset"
	"madledger/client/lib"
	"madledger/client/util"
	"madledger/common"
	coreTypes "madledger/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		"receiver's hex address to be issued in asset channel")
	issueViper.BindPFlag("address", issueCmd.Flags().Lookup("address"))

	issueCmd.Flags().BoolP("self", "s", false, "issue to your self")
	issueViper.BindPFlag("self", issueCmd.Flags().Lookup("self"))
}

func runIssue(cmd *cobra.Command, args []string) error {
	cfgFile := issueViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}
	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}

	//channelID can be empty
	channelID := issueViper.GetString("channelID")
	action := "person"
	if channelID != "" {
		action = "channel"
	}
	value := issueViper.GetInt("value")
	if value < 0 {
		return errors.New("cannot issue negative value")
	}

	receiver := issueViper.GetString("address")
	var recipient common.Address
	if receiver != "" {
		recipient = common.HexToAddress(receiver)
	}

	self := issueViper.GetBool("self")
	if self {
		recipient, _ = client.GetPrivKey().PubKey().Address()
	}

	payload, err := json.Marshal(asset.Payload{
		Action:    action,
		ChannelID: channelID,
		Address:   recipient,
	})
	if err != nil {
		return err
	}

	tx, err := coreTypes.NewTx(coreTypes.ASSETCHANNELID, coreTypes.IssueContractAddress, payload, uint64(value), "", client.GetPrivKey())
	if err != nil {
		return err
	}

	status, err := client.AddTxInOrderer(tx)
	table := util.NewTable()
	table.SetHeader("Status", "Error")

	table.AddRow(status, err)

	table.Render()

	return err
}
