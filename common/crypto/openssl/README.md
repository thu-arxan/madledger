# OpenSSL 安装

静态编译安装

下载`OpenSSL1.1.1`发行版解压

```sh
sudo apt install build-essential
wget https://www.openssl.org/source/openssl-1.1.1.tar.gz
tar -xzf openssl-1.1.1.tar.gz
cd openssl-1.1.1
./config --prefix=/usr/local/openssl --openssldir=/usr/local/openssl --shared
make
sudo make install
```

然互设置环境变量。

```bash
LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/openssl/lib
```
