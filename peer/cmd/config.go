package cmd

const (
	cfgTemplate = `#############################################################################
#############################################################################
#   This is a configuration file for the MadLedger peer.
#############################################################################

# Port should be an integer (default: 23456)
Port: 23456

# Bind address for the server (default: localhost)
Address: localhost

# Should be false or true (default: true)
Debug: true

# Configure for the BlockChain
BlockChain:
  # default: $GOPATH/src/madledger/peer/data/blocks
  # But in the production environment, you must provide a path
  Path: 

# Address of orderers
Orderer:
  Address:
    - localhost:12345

# DB only support leveldb now
DB:
  Type: leveldb
  # LevelDB
  LevelDB:
    # default: $GOPATH/src/madledger/peer/data/leveldb
    # But in the production environment, you must provide a path
    Dir: 
`
)