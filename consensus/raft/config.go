// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package raft

import (
	"madledger/consensus"
	"madledger/consensus/raft/eraft"
)

// Config is the config of consensus
type Config struct {
	id    uint64
	dir   string             // root dir for raft storage
	peers map[uint64]string  // id => grpc addr
	cc    consensus.Config   // consensus config
	ec    *eraft.EraftConfig // eraft config
}

// NewConfig is the constructor of Config
func NewConfig(dir string, id uint64, nodes map[uint64]string, join bool, cc consensus.Config) (*Config, error) {
	ec, err := eraft.NewEraftConfig(dir, id, nodes, join)
	if err != nil {
		return nil, err
	}

	return &Config{
		id:    id,
		dir:   dir,
		peers: nodes,
		cc:    cc,
		ec:    ec,
	}, nil
}
