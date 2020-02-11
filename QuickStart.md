# Quick start 教程

## 1. Requirement

- Go 环境
- solcjs: [Solidity编译器](https://github.com/ethereum/solc-js), 用于编译用户自己编写的智能合约(测试文件中给出了部分示例，可以先直接使用该示例)
  - solcjs --bin *.sol
  - solcjs --abi *.sol

## 2. 利用Makefile进行Quick start

为了便于运行，提供了[Makefile](start.mk)来帮助运行。

```bash
# 安装
make -f start.mk install
# 初始化, CONSENSUS指定共识协议类型，注意, bft暂不支持
make -f start.mk init CONSENSUS=raft/solo/bft
# 启动
make -f start.mk start
# 测试
make -f start.mk test
# 关闭
make -f start.mk stop
# 清楚相关测试数据，配置文件
make -f start.mk clean
```

## 3. Install

一共3个服务模块，`orderer、peer、client`，可以使用Makefile进行安装：

```bash
make install
```

也可以手动用`go install`进行安装：

```bash
# install orderer
go install madledger/orderer
# install peer
go install madledger/peer
# install client
go install madledger/client
```

安装成功后可以运行下列命令，如果有正常输出则代表安装成功:

```bash
orderer version # Orderer version v0.0.1
peer version    # Peer version v0.0.1
client version  # Client version v0.0.1
```

可以使用`[basename] -h`的方式查看每个模块的命令和参数

## 4. 配置文件

可以使用服务提供的`init`命令初始化生成配置文件

```bash
orderer init -c orderer.yaml
client init -c client.yaml
peer init -c peer.yaml
```

各模块配置文件中参数含义请查阅对应模块`config`包的文档。

## 5. Start

安装后，就可以按照相应规则生成配置文件并启动相关服务了。`env_local`文件夹中准备了一些可用于本机启动相关服务并运行测试的配置文件和测试脚本，开发者可利用该脚本和配置文件启动服务进行测试，也可自行生成。

系统目前支持三种共识协议，每一种共识协议下，`client, peer`模块的配置大体一致，`orderer`模块配置文件略有不同。具体区别可查看`orderer/config`目录下的[文档](orderer/conifig/README.md)。

### 5.1. Solo

`solo`模式的示例配置文件在，`env_local/solo`。其中:

- `orderer`服务只部署一个。
- `client`客户端不受限制，这里提供了1个客户端。
- `peer`服务也不受限制，这里部署3个服务，每个`peer`服务均与同一个`orderer`服务通信。

### 5.2. Raft

`raft`模式的示例配置文件在，`env_local/raft`。其中：

- `ordere`服务部署4个。
- `client`客户端不受限制，这里提供了6个客户端。
- `peer`服务不受限制，这里部署4个服务，每个`peer`服务均与所有`orderer`服务通信。

### 5.3. BFT

TODO
