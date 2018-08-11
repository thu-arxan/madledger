# Orderer

Orderer负责进行排序。

## 启动方法

### Init

该过程生成配置文件。如果已经有了配置文件可以跳过此过程。否则，执行下面命令。

```bash
orderer init
```

默认在当前目录下生成配置文件orderer.yaml，如果想要指定配置文件名称则执行如下命令。

```bash
orderer init -c $filepath.yaml
```

### Start

该过程根据配置文件启动Orderer节点，默认使用当前目录下orderer.yaml文件。

```bash
orderer start
```