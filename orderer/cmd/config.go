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

# Configure for the BlockChain
BlockChain:
  # Max time to create a block which unit is milliseconds (default: 1000)
  BatchTimeout: 1000
  # Max txs can be included in a block (defalut: 100)
  BatchSize: 100
  # Path to store the blocks (default: $GOPATH/src/madledger/orderer/data/blocks)
  # But in the production environment, you must provide a path
  Path: 
  # If verify the rightness of blocks (default: false)
  Verify: false

# Consensus mechanism configuration
Consensus:
  # will support solo, raft, pbft. Only support solo yet.
  Type: solo
`
)