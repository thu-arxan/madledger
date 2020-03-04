package cmd

const (
	cfgTemplate = `#############################################################################
#   This is a configuration file for the MadLedger client.
#############################################################################

# Should be false or true (default: true)
Debug: true

# Configure for the TLS
TLS:
  # Should be true of false (default: false)
  Enable: false
  # The path of CA cert, it should not be empty if Enable is true
  CA: 
  # Cert of the client, it should not be empty if Enable is true
  Cert: 
  # Key of the client, it should not be empty if Enable is true
  Key: 

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
