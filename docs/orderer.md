# Orderer

Orderer主要负责进行排序，还涉及到少量系统通道如_config,_asset的交易执行。

## 1. 启动方法

### 1.1. Init

该过程生成配置文件。如果已经有了配置文件可以跳过此过程。否则，执行下面命令。

```bash
orderer init
```

默认在当前目录下生成配置文件orderer.yaml，如果想要指定配置文件名称则执行如下命令。

```bash
orderer init -c $filepath.yaml
```

### 1.2. Start

该过程根据配置文件启动Orderer节点，默认使用当前目录下orderer.yaml文件。

```bash
orderer start
```

## 2. 配置文件说明

关于Orderer配置文件的具体描述，详见[Orderer配置文件](../orderer/config/README.md)。

## 3. 分布式部署

对于Orderer，除了支持测试用的solo外，还支持Raft以及拜占庭容错。

### 3.1. Raft

关于Raft的部署，在Orderer配置文件中有相关描述，请详细查阅或参考env_local中的raft样例。

### 3.2. 拜占庭容错

针对拜占庭容错，目前还在测试阶段，等测试完成之后进一步完善文档。

