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

# Configure for the TLS
TLS:
  # Should be true of false (default: true)
  Enable: true
  # The path of CA cert, it should not be empty if Enable is true
  CA: 
  # Cert of the peer, it should not be empty if Enable is true
  Cert: 
  # Key of the peer, it should not be empty if Enable is true
  Key: 

# Configure for the peer
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

# KeyStore manage some private keys
KeyStore:
  Key: <<<KEYFILE>>>
`
)
