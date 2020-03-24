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
	client, err := lib.NewClient(cfgFile)
	if err != nil {
		return err
	}
	return client.CreateChannel(name, true, nil, nil, 0, 1, 1000000)
}
