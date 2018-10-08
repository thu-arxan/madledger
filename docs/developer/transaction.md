# Transaction

## 结构

待完成。

### 流程

client封装好将tx发送给orderer节点进行排序，orderer节点打包并确定无重复交易进入blockchain，然后peer节点执行交易并存储tx结果（包括执行输出、可能的执行错误等），client可以向peer节点查询交易结果。

## 分类

用户通道的交易可以分为以下几类。

- 创建合约。
- 调用合约。
- 转账(是否支持待定)。

## Create

创建合约的目标地址是0x0000000000000000000000000000000000000000，payload是创建合约的代码，目前仅支持EVM。

### 合约地址

合约地址的生成规则为将通道名、发送者地址和代码合并之后进行ripemd160哈希，其可以简要描述如下。

```go
ripemd160.New().Write(channelID, sender, code)
```

如果创建的合约和以前已存在的合约地址重复，则会发生Duplicate Address错误。

## Call

调用目标合约的地址并给予所需的输入，目前仅支持EVM。

## Transfer

单纯的转账是否需要支持待定。