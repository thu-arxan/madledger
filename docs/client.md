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

#### 1.3.1 info

查看当前账户信息。

```bash
client account info
```

会看到如下所示的结果。

Address |
 ----   |
0xb65127172831e8e66a3f8310ea83d2eda1fcefc5

### 1.4 Tx

Tx命令负责client的交易部分。

#### 1.4.1 send

发送交易并获取交易结果。

```bash
client tx send -n $name -r $receiver -p $payload
```

待完善。