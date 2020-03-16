// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package node

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	homeDir, _ = os.Getwd()
	nodeCmd    = &cobra.Command{
		Use: "node",
	}
)

// Cmd return the channel command
func Cmd() *cobra.Command {
	nodeCmd.AddCommand(addCmd)
	nodeCmd.AddCommand(removeCmd)
	return nodeCmd
}
