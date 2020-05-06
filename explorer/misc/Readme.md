# 证书生成脚本使用说明

## 指令

* `gen_cert.py ca`，在指定目录（如`./CA`)，生成根证书`localhostCA.key`,`localhostCA.pem`
* `gen_cert.py auth`， 在指定目录（如`./cert`)，生成子证书`localhost.key`, `localhost.crt`,并将`CA.pem`的*绝对路径*超链接过去
* `gen_all.py`，多次调用`gen_cert.py auth`，为整个项目更新证书

具体参数见help


