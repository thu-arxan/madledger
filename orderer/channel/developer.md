## 算法(目前)思路
### manager.AddBlock
channel manager负责自己通道的区块
orderer的channel manager的AddBlock负责执行系统通道块或对应用通道块进行存储收费。
peer的channel manager的AddBlock执行系统/用户通道块。
流程(orderer)
* 判断是否应用通道非初始块（不对初始块收费），若是，判断通道账户中是否有足够的钱进行块的存储，若无，返回错误。
* 存储交易到数据库，区块到blockchain
* 存储成功则区块一定会被执行（？没看懂重启的时候数据库里的区块从哪里来的），进行存储收费
* 等待global channel对全局区块排序后，按顺序执行系统块

### manager.AddAssetBlock
Asset通道负责管理全局的账户余额，每个应用通道可以按一定比例让用户用asset里的余额交换自己生成的token，来进行应用通道的evm运算。
Asset目前支持的内置合约有三种issue, transfer和tokenExchange，用tx的recipient变量来指定用户想执行的合约。
三者共用`blockchain/asset/payload.go`中定义的payload，若payload中`Address == common.ZeroAddress`，该合约向`channelID`指定通道执行。
即receiver = (payload.Address == common.ZeroAddress) ? common.BytesToAddress(payload.ChannelID) : payload.Address。

* issue
执行对象为用户地址/通道，即向用户/通道发钱。
首先判断sender是否为asset通道admin，目前admin设置为第一个调用issue的用户的公钥。若是，向receiver发钱。
* transfer
sender向receiver转钱。
* tokenExchange
receiver必须为channel。
首先调用transfer，sender向receiver转账。
receiver收到钱后，按比例向receiver发放自己通道的token。

## 实现问题
执行时需要保存很多KVS，比如asset admin的公钥，token Exchange的ratio等等。这些值在数据库中需要用什么key去查询在peer/orderer中需要统一。
如果只是在代码中像magic number一样出现感觉很难维护，但目前没发现好方法统一写在哪个地方。
