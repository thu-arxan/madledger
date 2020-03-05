# Peer

执行节点。从Orderer节点获取区块并执行，修改更新global state，并且client可以从Peer获取交易的执行结果以及查询global state。

## 1. 启动方法

### 1.1. Init

该过程生成配置文件。如果已经有了配置文件可以跳过此过程。否则，执行下面命令。

```bash
peer init
```

默认在当前目录下生成配置文件peer.yaml，如果想要指定配置文件名称则执行如下命令。

```bash
peer init -c $filepath.yaml
```

### 1.2. Start

该过程根据配置文件启动Peer节点，默认使用当前目录下peer.yaml文件。

```bash
peer start
```

## 2. 配置文件说明

关于Peer配置文件的具体描述，详见[Peer配置文件](../peer/config/README.md)。

## 3. 分布式部署

目前，peer节点之间不需要任何通信，其完全依赖于排序节点所提供的区块数据并执行即可得到最终结果。