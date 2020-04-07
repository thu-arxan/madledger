# 证书生成脚本使用说明

## 指令

* `gen_cert.sh ca`，在脚本所在目录的`CA`子目录，生成根证书
* `gen_cert.sh`， 在当下的工作目录生成证书，证书名为`$NAME.crt`，`$NAME.key`，`CA.pem`

## 使用方式

首先将`gen_cert.sh`所在的目录加入PATH环境变量中，以方便在其他地方调用。

进入某个节点(如orderer)的cert目录，执行下面语句，（该语句告知脚本生成的cert的名字应该是`orderer.crt`）

```bash
export NAME=orderer
```

进入`cert`目录，执行

```bash
gen_cert.sh
```

此时会将CA的证书软链接到cert文件夹，并且依照CA的证书生成子证书。