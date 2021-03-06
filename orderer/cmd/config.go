// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package cmd

const (
	cfgTemplate = `#############################################################################
#   This is a configuration file for the MadLedger orderer.
#############################################################################

# Port should be an integer (default: 12345)
Port: 12345

# Bind address for the server (default: localhost)
Address: localhost

# Should be false or true (default: true)
Debug: true

# Configure for the TLS
TLS:
  # Should be true of false (default: true)
  Enable: true
  # The path of CA cert, it should not be empty if Enable is true
  CA: 
  # Cert of the orderer, it should not be empty if Enable is true
  Cert: 
  # Key of the orderer, it should not be empty if Enable is true
  Key: 

# Configure for the BlockChain
BlockChain:
  # Max time to create a block which unit is milliseconds (default: 1000)
  BatchTimeout: 1000
  # Max txs can be included in a block (defalut: 100)
  BatchSize: 100
  # Path to store the blocks (default: orderer/data/blocks)
  Path: <<<BlockChainPath>>>
  # If verify the rightness of blocks (default: false)
  Verify: false

# Consensus mechanism configuration
Consensus:
  # will support solo, raft, bft. Only support solo yet and bft is constructed now.
  Type: <<<ConsensusType>>>
  # Tendermint is the bft consensus.
  Tendermint:
    # The path of tendermint (default: orderer/.tendermint)
    Path: <<<TendermintPath>>>
    # Some ports
    Port:
      P2P: 26656
      RPC: 26657
      APP: 26658
    # ID means to identity in p2p connections
    ID: <<<TendermintP2PID>>>
    # P2P Persistent Address, like c395828cc2baaa6f6af2bd13ce62d1e9484919c8@localhost:36656
    P2PAddress:
      -
  # Raft is the raft consensus
  Raft:
    # The path of raft
    Path: <<<RaftPath>>>
    # ID should be int, and it should not be duplicate
    ID:
    # Node should be like 1@localhost:12345
    Nodes:
      -
    # Should be true of false (default: false)
    Join: false

# DB only support leveldb now
DB:
  Type: leveldb
  # LevelDB
  LevelDB:
    # The path of leveldb (default: orderer/data/leveldb)
    Path: <<<LevelDBPath>>>
`
)
