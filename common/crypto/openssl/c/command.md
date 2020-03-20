# 常用命令

## 签名验签

```sh
openssl dgst -sign priv.pem -sha256 -out sign.txt in.txt
openssl dgst -verify pub.pem -sha256 -signature sign.txt in.txt
```
