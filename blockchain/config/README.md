# Config

定义了Config Channel的相关结构。

## Payload

支持以下几种配置属性。

- ChannelID: 所操作Channel的ID。
- Profile: Channel的相关属性。包含Open（是否公开）等。
- Version: 版本号，当前为１。

// 以下三个为应用通道专用
// 创建应用通道时务必设置，默认值分别为1/1/一百万（在文件`peer/channel/config.go: setChannelConfig.go`

- GasPrice: 每个gas的单价
- AssetTokenRatio: 一个asset能兑换多少token
- MaxGas: 通道中每次evm运行时可以消耗的最大数值，默认为？

- BlockPrice: 存储区块时系统对用户进行收费，默认为0	
