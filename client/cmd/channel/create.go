// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package channel

import (
	"errors"
	"madledger/client/lib"

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
	createCmd.Flags().StringP("config", "c", "client.yaml", "The config file of client")
	createViper.BindPFlag("config", createCmd.Flags().Lookup("config"))
	createCmd.Flags().StringP("gasPrice", "g", "", "Numbers of token spent for one gas")
	createViper.BindPFlag("gasPrice", createCmd.Flags().Lookup("gasPrice"))
	createCmd.Flags().StringP("maxGas", "m", "", "max gas spent for transaction execution")
	createViper.BindPFlag("maxGas", createCmd.Flags().Lookup("maxGas"))
	createCmd.Flags().StringP("ratio", "r", "", "Numbers of token exchanged from one asset")
	createViper.BindPFlag("ratio", createCmd.Flags().Lookup("ratio"))

}

func runCreate(cmd *cobra.Command, args []string) error {
	cfgFile := createViper.GetString("config")
	if cfgFile == "" {
		return errors.New("The config file of client can not be nil")
	}
	name := createViper.GetString("name")
	if name == "" {
		return errors.New("The name of channel should be [a-z0-9]{1,32} such as test, test01 and etc")
	}

	gasPrice := createViper.GetInt("gasPrice")
	if gasPrice < 0 {
		return errors.New("gasPrice cannot be negative")
	}

	maxGas := createViper.GetInt("maxGas")
	if maxGas < 0 {
		return errors.New("maxGas cannot be negative")
	}

	ratio := createViper.GetInt("ratio")
	if ratio < 0 {
		return errors.New("gasPrice cannot be negative")
	}

	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}
	return client.CreateChannel(name, true, nil, nil, uint64(gasPrice), uint64(ratio), uint64(maxGas))
}
