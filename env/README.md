# env

一些配置好的环境以方便开发和使用。主要是单机solo和拜占庭容错bft。

## 1. solo

待完善。

## 2. bft

该文件夹初始化了可以运行bft共识的一些环境，其中orderers下面分别由0、1、2、3四个文件夹，通过分别在文件夹目录下运行start.sh可以运行4个orderer节点，如下所示。

```bash
cd bft/orderers/0
. start.sh
```

当所有节点都启动后进入clients/admin文件夹。

```bash
client channel list #查看所有通道
client channel create -n test #创建test通道
client channel list #查看所有通道
```

可以看到所有orderer节点都创建了所希望的通道，说明这些orderer节点之间进行了共识。

要注意的是，bft的实现基于Tendermint且很不完善，存在一些可能的bug，这个环境也不完善，比如缺少peers节点因此只能创建通道但是没办法支持交易，这都是需要完善的部分。

另外，如果需要将bft环境恢复初始状态，可以运行下面的脚本。

```bash
. init.sh
```

如果想要手动创建Tendermint的运行环境，可以参考下面的步骤。

> TODO

## 3. raft

该文件夹初始化了可以运行raft共识的一些环境，其中orderers下面分别有1,2,3三个文件夹。
