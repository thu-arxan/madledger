# Transaction

## 1. 结构

交易的定义位于core/types/tx.go中。

需要注意的主要有以下几个部分。

- 一般地，交易都需要被签名，但是global通道内的除外，因为没有一个主体可以代表通道签名交易并发送至global通道，global通道上链依赖于各个orderer生成交易并共识。
- tx的hash不包含时间等可能具有不确定的部分，如global的tx由于各个orderer节点生成无签名的tx，这里的time是各个orderer节点的本地时间，因此可能生成不完成一样的tx，但是tx的hash是一样的，所以仍然是同样的交易。
- tx的ID目前是简单的对tx的hash进行16进制编码。

## 2. 流程

client封装好将tx发送给orderer节点进行排序，orderer节点打包并确定无重复交易进入blockchain，然后peer节点执行交易并存储tx结果（包括执行输出、可能的执行错误等），client可以向peer节点查询交易结果。

## 3. 分类

用户通道的交易可以分为以下几类。

- 创建合约。
- 调用合约。
- 转账。

### 3.1. Create

创建合约的目标地址是0x0000000000000000000000000000000000000000，payload是创建合约的代码，目前仅支持EVM。

#### 3.1.1. 合约地址

合约地址的生成规则目前为将通道名、发送者地址和代码合并之后进行ripemd160哈希，其可以简要描述如下。

**将会很快发生修改**

```go
ripemd160.New().Write(channelID, sender, code)
```

如果创建的合约和以前已存在的合约地址重复，则会发生Duplicate Address错误。

### 3.2. Call

调用目标合约的地址并给予所需的输入，目前仅支持EVM。

### 3.3. Transfer

转账正在支持中。