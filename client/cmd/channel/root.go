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
	"os"

	"github.com/spf13/cobra"
)

var (
	homeDir, _ = os.Getwd()
	channelCmd = &cobra.Command{
		Use: "channel",
	}
)

// Cmd return the channel command
func Cmd() *cobra.Command {
	channelCmd.AddCommand(createCmd)
	channelCmd.AddCommand(listCmd)
	return channelCmd
}
