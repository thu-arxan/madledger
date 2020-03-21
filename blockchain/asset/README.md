# Account

定义了Account Channel相关的数据结构。

## Payload
```
// Payload specify contract receiver
type Payload struct {
	// if address == common.ZeroAddress
	// This is an op to channel

	Address   common.Address
	ChannelID string
}
```
若Address不为common.ZeroAddress，该合约向address执行。
否则该合约向channelID指定通道执行