# 快速上手

## 1 编译

在项目scripts目录下执行编译脚本即可。

```bash
. build.sh
```

然后分别执行下述命令查看是否配置正确。

```bash
orderer version
peer version
client version
```

如果看到以下类似结果则说明配置正确。

```bash
Orderer version v0.0.1
Peer version v0.0.1
Client version v0.0.1
```

否则，可能是$PATH设置有误，请将$GOPATH/bin添加为环境变量。

## 2 环境

为了方便测试以及使用，初始化了环境于env文件夹。