/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"madledger/explorer/server"
)

// runCmd represents the run command
var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Start a blockchain server here.",
		Long: `Start a blockchain server. Two mode are supported:
server : This will start a client process in the server side. All operation was execute in the backend.
explorer : This will only provide a static website, All operation will be execute in user's explorer
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			host := runViper.GetString("host")
			port := runViper.GetInt("port")
			config := runViper.GetString("config")
			return server.RunServer(host, port, config)
		},
	}
	runViper = viper.New()
)
var cfgFile  = ""

func init() {
	rootCmd.AddCommand(runCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")
	runCmd.Flags().StringP("config", "c", "./config/explorer-client.yaml", "config file path")
	runViper.BindPFlag("config", runCmd.Flags().Lookup("config"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	runCmd.Flags().StringP("host", "", "127.0.0.1", "host to run server")
	runViper.BindPFlag("host", runCmd.Flags().Lookup("host"))

	runCmd.Flags().IntP("port", "", 8080, "port to run server")
	runViper.BindPFlag("port", runCmd.Flags().Lookup("port"))
}
