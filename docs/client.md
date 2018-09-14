# Client

客户端。负责和Orderer和Peer交互。

## 1 使用方法

### 1.1 Init

该过程生成初始配置文件以及私钥。如果已经初始化则跳过当前步骤。

```bash
client init
```

其中，配置文件默认名为client.yaml，私钥生成于.keystore文件夹。

### 1.2 Channel

Channel命令负责client对通道的相关操作。可以通过下述命令查看所提供的功能。

```bash
client channel -h
```

#### 1.2.1 list

查看当前所有的通道信息。

```bash
client channel list
```

会看到如下所示的结果。

Name | System | BlockSize
---- | --- | ---
_global | true | 1
_config | true | 1

#### 1.2.2 create

创建新的通道。

```bash
client channel create -n $name
```

### 1.3 Account

Account命令负责client对账户相关的操作。

#### 1.3.1 list

查看当前账户信息。

```bash
client account list
```

会看到类似于如下所示的结果。

Address |
 ----   |
0xb65127172831e8e66a3f8310ea83d2eda1fcefc5|

### 1.4 Tx

Tx命令负责client的交易部分。

#### 1.4.1 create

根据一个sol文件编译生成的bin文件，创建一个合约。

```bash
client tx create -b $bin -n $name
```

返回结果类似如下所示。

BlockNumber | BlockIndex | ContractAddress  
---- | --- | ---
1 | 0 | 0x16987f7117b6f0f0a8d55f6f15d6d8cb82fec58a

#### 1.4.2 call

根据一个sol文件编译生成的abi文件，调用一个合约。

```bash
client tx call -a $abi -n $name -f $func -p $payload -r $receiver
```

返回结果类似如下所示。

BlockNumber | BlockIndex | Output
---- | --- | ---
4 | 0 | [1314]
