# TPS 变更记录

## Raft

- dev:c349946e11bb0d520d2d749d413f4dc1b22fdb3f, 2020-02-13, 3700
- feature/raft:415fdd4957a5fb6e6c0838f2707534eb5e56e6f1, 2020-02-17, 6000
  - 重构了Raft, 由Raft原生支持多通道的打包，每个通道的区块单独共识
  - 在peer.Manager中的轮询orderer新交易的产生增加了Sleep(50ms)==> 需进一步优化
