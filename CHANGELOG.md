# ChangeLog

## Unreleased

## Feature

## TODO

* [x] EVM移植
* [x] Account suicide之后，相关数据需要从数据库中删除
* [ ] 测试Account suicide之后，相关数据是否删除
* [ ] Raft部分代码重构
* [ ] Raft部分，多链共识
* [ ] 添加AddNode, RemoveNode接口
* [ ] 数据库操作的error处理
* [ ] db.ChainNum的处理
* [ ] raft测试
* [ ] 修复重启读入Channel问题

## Fix Me

* [ ] 修改增加Tx.Data.Value后部分代码未传入Value的问题
* [ ] 初始化时，Channel 的BlockNum需要从数据库中读入
* [ ] Peer manager fetchBlock由orderer主动通知
