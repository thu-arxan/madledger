# HTTP Client

端口为gRPC端口-100

## ListChannel

查看当前所有的通道信息
Input: system bool
Output: []client.ChannelInfo, error
Usage:  

``` go
info, err := client.ListChannelByHTTP(true)
```

## CreateChannel

创建一个新的通道
Input: channelID string, public bool, admins, members []*core.Member, gasPrice uint64, ratio uint64, maxGas uint64
Output: error
Usage:  

``` go
err = client.CreateChannelByHTTP("private", false, nil, peers, 0, 1, 10000000)
```

## AddTx

用来进行call contract, create contract, issue, transfer, token exchange 等操作
Input: tx *core.Tx
Output: *pb.TxStatus, error
Usage:  

``` go
_, err = client.AddTxByHTTP(coreTx)
```

## GetAccountBalance

获取（地址为address的）用户的asset数量
Input: address common.Address
Output: uint64, error
Usage:  

``` go
balance, err := client.GetAccountBalanceByHTTP(receiverAddress)
```

## GetHistory

查看当前账户的交易历史
Input: address []byte
Output: *pb.TxHistory, error
Usage:  

``` go
history, err := client.GetHistoryByHTTP(addr)
```

## GetTokenInfo

获取（地址为address的）用户在（通道名为channelID的）通道的Token数量
Input: address common.Address, channelID []byte
Output: uint64, error
Usage:  

``` go
token, err := client.GetTokenInfoByHTTP(receiver, []byte("test"))
```
