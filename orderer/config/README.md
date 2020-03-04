# Orderer配置说明

```yaml
#############################################################################
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
  Enable: false
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
  Path: /home/liuyihua/gopath/src/madledger/orderer/config/data/blocks
  # If verify the rightness of blocks (default: false)
  Verify: false

# Consensus mechanism configuration
Consensus:
  # will support solo, raft, bft. Only support solo yet and bft is constructed now.
  Type: solo
  # Tendermint is the bft consensus.
  Tendermint:
    # The path of tendermint (default: orderer/.tendermint)
    Path: /home/liuyihua/gopath/src/madledger/orderer/config/.tendermint
    # Some ports
    Port:
      P2P: 26656
      RPC: 26657
      APP: 26658
    # ID means to identity in p2p connections
    ID: 200fd1ad575274a457a69af7baf9b974fea4b2ff
    # P2P Persistent Address, like c395828cc2baaa6f6af2bd13ce62d1e9484919c8@localhost:36656
    P2PAddress:
      -
  # Raft is the raft consensus
  Raft:
    # The path of raft
    Path: /home/liuyihua/gopath/src/madledger/orderer/config/.raft
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
    Path: /home/liuyihua/gopath/src/madledger/orderer/config/data/leveldb
```

上述`yaml`配置文件中的注释已经足够详尽了，这里只对`Consensus`的配置进行一下说明：

- solo: solo模式下，只需要将Consensus.Type设置为solo即可，不需要做其他的配置。
- raft: raft模式下,各参数含义如下：
  - Path: raft日志、存储目录。
  - ID: 节点ID, 大于0，且不可重复。
  - Nodes: 启动时已知的orderer节点列表，格式为id@host:port,需要包含自身
  - Join: 如果系统已经启动并运行，需要临时增加一个节点，该节点的Join参数需要设置为true，其他情况设置为false
- bft: bft模式下，各参数含义如下：
  - Path: tendermint的日志、存储目录
  - Port:
    - P2P: p2p服务监听的端口号
    - RPC: RPC服务监听的端口号
    - APP: 应用服务监听的端口号
  - ID: 节点在p2p通信中的识别号，不可重复
  - P2PAddress: 网络中其他orderer节点的p2p通信地址，格式为ID@host:port
