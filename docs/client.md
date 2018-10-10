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
client tx call -a $abi -n $name -f $func -i $input -r $receiver
```

注意，-i参数可以有多个或者没有。

```bash
client tx call -a Balance.abi -n test -f add -i 1314 -r 0x16987f7117b6f0f0a8d55f6f15d6d8cb82fec58a
client tx call -a Balance.abi -n test -f get -r 0x16987f7117b6f0f0a8d55f6f15d6d8cb82fec58a
client tx call -a Balance.abi -n test -f mock -i 1314 -i 520 -r 0x16987f7117b6f0f0a8d55f6f15d6d8cb82fec58a
```

返回结果类似如下所示。

BlockNumber | BlockIndex | Output
---- | --- | ---
4 | 0 | [1314]

#### 1.4.3 history

查看当前账户的交易历史。

```bash
client tx history
```

返回的结果类似如下所示。

Channel | TxID
 ------ | ---
 test   | bb92201e40a00d643fae32ab5bd95bec46242951e88040eb7bbf0a00f50a0626
 test   | 8ac8890bb62892721296f5d699476cf7f3961fa510a761252e1e3436818cf182
 test   | 47a4af918ecce1b1450984dd552f3d51e5fc4edb5d340eff7d05ab42a55fe025
 test   | 567be0d1c8a0181c7c90d0a046fe578354bd0389490500735da8ed0dd687b167
 test   | d5e35a0c850b3b6ef2ff028614ce49d6748b1d50e0b0dabad44f89f3d548c5e6
 _config| a77f2fca39eadf40d19dd9f13246b70b30c28e65d3726f4c15261bc512edcb70
 _config| a77f2fca39eadf40d19dd9f13246b70b30c28e65d3726f4c15261bc512edcb70

## 2 分布式

### 2.1 Orderer

目前，尚不支持分布式的排序节点。

### 2.2 Peer

目前，客户端可以支持分布式的Peer节点。客户端独立的从各个节点获取结果并进行比较，如果能够获取足够多的相同结果则可认为得到了正确的执行结果，否则报错（注意，如果发生了非网络原因的报错说明系统出现了严重的问题，因为理论上所有的诚实Peer节点应当得到相同的执行结果）。