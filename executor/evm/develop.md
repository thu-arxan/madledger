# MadEVM 实现方案

## RunBlock

`RunBlock`时初始化`Contex`, 对每一个Tx初始化一个`EVM`, 其中Context作为所有Tx的共享上下文

EVM执行Tx时需要有两层Cache,

### evm

```go
func New(bc Blockchain, db DB, ctx *Context) *EVM
```

evm内部会对传入的DB增加一层封装(Cache), Cache会对本Tx内部对数据的修改进行缓存

一个交易处理成功之后，evm.Cache会进行Sync，此时，会将相关数据写入到db.WriteBatch中

db.WriteBatch的内容需要在后续的Tx处理过程中可读，即其写入的内容需要通过ctx进行在同一个block的tx之间顺序传递

一个block处理完成后，最后调用ctx.BlockFinalize写入到持久化DB的WriteBatch中，并最终落盘

## `vendor/evm`

`vendor/evm`包主要对外暴露以下接口:

```go
// New is the constructor of EVM
func New(bc Blockchain, db DB, ctx *Context) *EVM {
    return &EVM{
        bc:             bc,
        cache:          NewCache(db),
        memoryProvider: DefaultDynamicMemoryProvider,
        ctx:            ctx,
    }
}

// Create create a contract account, and return an error if there exist a contract on the address
func (evm *EVM) Create(caller Address) ([]byte, Address, error) {}

// Call run code on evm, and it will sync change to db if error is nil
func (evm *EVM) Call(caller, callee Address, code []byte) ([]byte, error) {}


// DB describe what function that db should provide to support the evm
type DB interface {
    // Exist return if the account exist
    // Note: if account is suicided, return true
    Exist(address Address) bool
    // GetStorage return a default account if unexist
    GetAccount(address Address) Account
    // Note: GetStorage return nil if key is not exist
    GetStorage(address Address, key []byte) (value []byte)
    NewWriteBatch() WriteBatch
}

// WriteBatch define a batch which support some write operations
type WriteBatch interface {
    SetStorage(address Address, key []byte, value []byte)
    // Note: db should delete all storages if an account suicide
    UpdateAccount(account Account) error
    AddLog(log *Log)
}

// Blockchain describe what function that blockchain system shoudld provide to support the evm
type Blockchain interface {
    // GetBlockHash return ZeroWord256 if num > 256 or num > max block height
    GetBlockHash(num uint64) []byte
    // CreateAddress will be called by CREATE Opcode
    CreateAddress(caller Address, nonce uint64) Address
    // Create2Address will be called by CREATE2 Opcode
    Create2Address(caller Address, salt, code []byte) Address
    // Note: NewAccount will create a default account in Blockchain service,
    // but please do not append the account into db here
    NewAccount(address Address) Account
    // BytesToAddress provide a way convert bytes(normally [32]byte) to Address
    BytesToAddress(bytes []byte) Address
}
```
