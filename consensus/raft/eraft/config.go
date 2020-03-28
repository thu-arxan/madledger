// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package eraft

import (
	"fmt"
	pb "madledger/consensus/raft/protos"
)

// EraftConfig is the config of eraft
// The raft will use some linear address to make sure the service can config and run simple
// If raft blockchain use 12346 then eraft service will use 12347
type EraftConfig struct {
	id      uint64
	dbDir   string
	walDir  string
	snapDir string
	// peers are eraft address
	peers map[uint64]string
	join  bool //node is joining an existing cluster
	// The url of node, maybe ip or domain
	url string
	// The port of eraft port, eraftPort = chainPort + 1
	eraftPort int

	snapshotInterval uint64
}

// NewEraftConfig is the constructor of EraftConfig
// works on dir and listen on address, id is the id of raft node, nodes is a url map of all nodes
func NewEraftConfig(dir string, id uint64, nodes map[uint64]string, join bool) (*EraftConfig, error) {
	url, chainPort, err := pb.ParseRaftAddress(nodes[id])
	if err != nil {
		return nil, err
	}

	return &EraftConfig{
		id:               id,
		dbDir:            fmt.Sprintf("%s/db", dir),
		walDir:           fmt.Sprintf("%s/wal", dir),
		snapDir:          fmt.Sprintf("%s/snap", dir),
		peers:            nodes,
		join:             join,
		url:              url,
		eraftPort:        chainPort + 1,
		snapshotInterval: 100,
	}, nil
}

// GetID return the id
func (c *EraftConfig) GetID() uint64 {
	return c.id
}

// GetPeers ...
func (c *EraftConfig) GetPeers() map[uint64]string {
	return c.peers
}

func (c *EraftConfig) getERaftAddress() string {
	return fmt.Sprintf("%s:%d", c.url, c.eraftPort)
}
