# 开发指南

## １ 准备工作

### 1.1 protobuf

如果还未安装protobuf,请点此[下载protobuf](https://github.com/google/protobuf/releases)，选择系统平台对应的Protocol Buffers v3.6.0压缩包并解压,然后将解压后bin/protoc文件夹复制到/usr/bin下面，以及把include/google文件夹复制到/usr/include里面。

学习文档请参见如下地址。

- [英语官方文档](https://developers.google.com/protocol-buffers/docs/proto3)
- [中文官方文档翻译](http://colobu.com/2017/03/16/Protobuf3-language-guide/)

## 2 目录说明

### 2.1 common

定义了一系列常用的函数以及密码学相关的函数，其不应当依赖于MadLedger的任何包。

### 2.2 core

定义了一系列基本的数据结构以及其相关操作，其仅依赖于common。

### 2.3 protos

定义了一系列传输结构以及服务。可能需要提供到core中数据结构的转换方式。

- compile.sh: 提供了快速的编译脚本。

### 2.4 orderer

Orderer相关的代码。

### 2.5 peer

Peer相关的代码。

### 2.6 executor

一系列执行器，目前仅支持EVM，且需要意识到EVM尚未经过完善的测试部分，还需要进一步的测试完善。

### 2.7 blockchain

Blockchain相关的代码。

### 2.8 docs

项目文档。

### 2.9 env

测试环境。

### 2.10 scripts

- build.sh:　编译MadLedger。
- test.sh: 运行所有测试。
- vendor.sh: 由于git不会提交包含.git的文件夹，所以需要删除vendor中的.git文件。

## 3 注意事项

### 3.1 测试

提交代码前请确保通过了所有的测试。

```bash
. build.sh
. test.sh
```