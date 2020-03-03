# 测试说明

## 1. global中出现重复tx的测试说明

### 1.1. 关于测试文件

- 在env/bft目录下有完整的client和orderer的配置，其中其中每个client里面有`test.sh`,`test1.sh`,`create_ca.sh`。
- `test.sh`用于创建多个通道，如client 0 创建通道`test10`、`test20`、`test30`等，数目可以自己通过修改shell脚本发生改变。目前是创建8个新通道。
- `test1.sh`与test.sh的作用一样，最开始只是为了在不初始化环境的情况下，重新测试通道的创建。
- `create_ca.sh`用于为指定通道创建智能合约，虽然peer节点没有启动会报错，但是在通道中会创建对应的block。在使用create_ca.sh之前，需要为每个通道client创建对应的通道。如`client 0`创建`test0`通道，`client 1`创建`test1`通道。
- **注意：** 如果使用`init.sh`重新初始化env后，再通过`start.sh`启动orderer后，发现所有的orderer会自动读取之前的信息。**解决办法：** 使用`rm -rf .tendermint`删除`.tendermint`，再通过`start.sh`重新启动。
- **PS：** 可以考虑直接在`start.sh`中添加`rm -rf .tendermint`。

### 1.2. 关于代码

- 主要改动集中在`orderer/channel/manager.go`中。
- `Start()`方法中，首先对`txs, _ := manager.getTxsFromConsensusBlock(cb)`中的`getTxsFromConsensusBlock()`进行了改动，在向legal中添加tx时，主动判断是否与legal中存在的tx重复，若重复则不添加，并打印相关信息。
- `Start()`方法中，接下来对非global通道新增block，向global通道发送tx做了改动，添加了打印tx的代码`log.Debugf("Channel %s add tx %s to global channel.", manager.ID, tx.ID)`。

### 1.3. 自己的想法

- 从非global通道新增block向global添加交易的log来看，每个orderer都只显示了一条添加信息。不太可能是重复添加？

```shell
# config中新增block
INFO[0068] Channel _config create new block 2, hash is 54c0a99995a0cdfce3f74a9c5615ec36c3a6c8298b78e42282d3d420f5c23003  app=orderer package=channel
# config通道向global通道添加tx
INFO[0068] Channel _config add tx f90b86eeddc6d4bb59e7f86c0a3b1e7158adce9549c91b32cd9216d79106e16a to global channel  app=orderer package=channel
```

- 若不是重复添加，那什么会出现重复的tx？会不会是具体的添加tx的函数`manager.coordinator.GM.AddTx(tx)`那里出现的问题？

## 2. client访问多个ordererClient，实现`channel list`的测试说明

### 2.1. 关于client.yaml

- 在`client.yaml`中修改了Orderer/Address，添加了两个不存在的orderer address
- 使用`client channel list`进行测试时，可以将所有的orderer address都配置为不存在的地址。

### 2.2. 关于代码

- 主要改动集中在`madledger/client/lib/client.go`中。
- 将`ordererClient`改成了数组，用于存储多个ordererCleint的信息。
- `getOrdererClient()改为`getOrdererClients()，用于获取所有的ordererClient。
- `ListChannel()`中，遍历ordererClient直到成功获取the info of channels。

## 3. client访问多个ordererClient，实现`channel create`的测试说明

### 3.1. 使用之前配置的`test.sh`实现`channel create`命令的测试

- 多个client并发创建channel，没有问题。
- 中途突然关闭某个orderer，发现所有的global倒是几乎同步了，但是不同orderer的config不是同步的，导致global上的block比config通道的block多。**这个bug待修复**

### 3.2. 关于代码

- 主要改动集中在`madledger/client/lib/client.go`的`CreateChannel()`方法中。
