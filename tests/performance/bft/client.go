package bft

import (
	client "madledger/client/lib"
)

const (
	clientConfigTemplate = `#############################################################################
#   This is a configuration file for the MadLedger client.
#############################################################################

# Should be false or true (default: true)
Debug: true

# Address of orderers
Orderer:
  Address:
    - localhost:12345
    - localhost:23456
    - localhost:34567
  
# Address of peers
Peer:
  Address:
    <<<ADDRESS1>>>
    <<<ADDRESS2>>>
    <<<ADDRESS3>>>

# KeyStore manage some private keys
KeyStore:
  Keys:
    - <<<KEYFILE>>>
`
)

var (
	clients    = make([]*client.Client, 400)
	clientInit = false
)
