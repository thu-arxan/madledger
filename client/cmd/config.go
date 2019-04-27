package cmd

const (
	cfgTemplate = `#############################################################################
#   This is a configuration file for the MadLedger client.
#############################################################################

# Should be false or true (default: true)
Debug: true

# Address of orderers
Orderer:
  Address:
    - localhost:12345
  
# Address of peers
Peer:
  Address:
    - localhost:23456

# KeyStore manage some private keys
KeyStore:
  Keys:
    - <<<KEYFILE>>>
`
)
