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

// Consensus defines functions that a consensus implementation should provide
type Consensus interface {
	// Start start the consensus service
	Start() error
	// Stop stop the consensus service
	Stop() error
	// Add a tx, it should return after the tx has been included in a block
	AddTx(tx *core.Tx) error
	// AddChannel will add a consensus of a channel, maybe we can support SetChannel in fucture
	// FUTURE: May support set channel
	AddChannel(channelID string, cfg Config) error
	// GetBlock return the block or error right away if async is false, else return block until the block is created.
	GetBlock(channelID string, num uint64, async bool) (Block, error)
}
