# env

一些配置好的环境以方便本地开发和使用。支持`solo、raft`共识

## 1. solo

待完善。

## 2. raft

该文件夹初始化了可以运行raft共识的一些环境，其中orderers下面分别由0、1、2、3四个文件夹，通过分别在文件夹目录下运行start.sh可以运行4个orderer节点，如下所示。

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

## 3. raft

TODO
