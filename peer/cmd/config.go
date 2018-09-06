package cmd

const (
	cfgTemplate = `#############################################################################
#   This is a configuration file for the MadLedger peer.
#############################################################################

# Port should be an integer (default: 23456)
Port: 23456

# Bind address for the server (default: localhost)
Address: localhost

# Should be false or true (default: true)
Debug: true

# Address of orderers
Orderer:
  Address:
    - localhost:12345
`
)
