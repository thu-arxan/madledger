# Client配置说明

```yaml
#############################################################################
#   This is a configuration file for the MadLedger client.
#############################################################################

# Should be false or true (default: true)
Debug: true

# Configure for the TLS
TLS:
  # Should be true of false (default: true)
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
    - localhost:34567

# KeyStore manage some private keys
KeyStore:
  Keys:
    - .key.pem
```

其中，Client的配置文件如上所示。对于部分选项，下文中会进行一些解释。