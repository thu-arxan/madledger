# 快速上手

## 1. Requirement

### 1.1. Go 环境

Go语言的安装请自行搜索，版本不小于1.10。

### 1.2. OpenSSL1.1.1

```sh
sudo apt install build-essential
wget https://www.openssl.org/source/openssl-1.1.1.tar.gz
tar -xzf openssl-1.1.1.tar.gz
cd openssl-1.1.1
./config --prefix=/usr/local/openssl --openssldir=/usr/local/openssl --shared
make
sudo make install
```

然后设置环境变量。

```bash
LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/openssl/lib
```

### 1.3. solcjs

- solcjs: [Solidity编译器](https://github.com/ethereum/solc-js), 用于编译用户自己编写的智能合约(测试文件中给出了部分示例，可以先直接使用该示例)
  - solcjs --bin *.sol
  - solcjs --abi *.sol

## 2. Install

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

可以使用`[basename] -h`的方式查看每个模块的命令和参数。

## 3. Start

这里描述了几种启动的方式，分别适用于下列场景。

- 通过预先写好的脚本快速在单机上启动一个配置好的MadLedger网络。
- 利用预先构建的环境在单机上手动启动一个配置好的MadLedger网络。
- 根据实际情况一步步部署搭建一个MadLedger网络。

### 3.1. 利用Makefile进行Quick start

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

### 3.2. Env_local启动

安装后，就可以按照相应规则生成配置文件并启动相关服务了。`env_local`文件夹中准备了一些可用于本机启动相关服务并运行测试的配置文件和测试脚本，开发者可利用该脚本和配置文件启动服务进行测试，也可自行生成。

系统目前支持三种共识协议，每一种共识协议下，`client, peer`模块的配置大体一致，`orderer`模块配置文件略有不同。具体区别可查看`orderer/config`目录下的[文档](orderer/conifig/README.md)。

#### 3.2.1. Solo

`solo`模式的示例配置文件在，`env_local/solo`。其中:

- `orderer`服务只部署一个。
- `client`客户端不受限制，这里提供了1个客户端。
- `peer`服务也不受限制，这里部署3个服务，每个`peer`服务均与同一个`orderer`服务通信。

#### 3.2.2. Raft

`raft`模式的示例配置文件在，`env_local/raft`。其中：

- `ordere`服务部署4个。
- `client`客户端不受限制，这里提供了6个客户端。
- `peer`服务不受限制，这里部署4个服务，每个`peer`服务均与所有`orderer`服务通信。

#### 3.2.3. BFT

TODO

### 3.3. 手动构建

手动构建环境需要依次为orderer、peer及client准备工作目录并配置相关文件。

#### 3.3.1. Orderer

关于Orderer的详细配置见[Orderer](orderer.md)，也可参考env_local中的相关配置。

#### 3.3.2. Peer

关于Peer的详细配置见[Peer](peer.md)，也可参考env_local中的相关配置。

#### 3.3.3. Client

关于Client的详细配置见[Client](client.md)，也可参考env_local中的相关配置。

