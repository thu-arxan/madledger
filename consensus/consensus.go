// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package consensus

import "madledger/core"

// Consensus is the interface
type Consensus interface {
	// Start Consensus service
	Start() error
	// Stop stop the consensus
	Stop() error
	// Add a tx
	AddTx(tx *core.Tx) error
	// AddChannel will add a consensus of a channel
	AddChannel(channelID string, cfg Config) error
	// GetBlock return the block or error right away if async is false, else return block until the block is created.
	GetBlock(channelID string, num uint64, async bool) (Block, error)
}
