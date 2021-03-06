# Copyright (c) 2020 THU-Arxan
# Madledger is licensed under Mulan PSL v2.
# You can use this software according to the terms and conditions of the Mulan PSL v2.
# You may obtain a copy of Mulan PSL v2 at:
#          http://license.coscl.org.cn/MulanPSL2
# THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
# EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
# MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
# See the Mulan PSL v2 for more details.

FROM ubuntu:18.04

# go environment
RUN apt-get update \
    && apt-get install -y --no-install-recommends apt-utils \
    && apt-get install -y wget \
    && apt-get install -y lsof \
    && apt-get install -y sudo \
    && wget http://digiccy.liuyihua.com/go1.12.9.linux-amd64.tar.gz \
    && tar -xzf go1.12.9.linux-amd64.tar.gz -C /usr/local \
    && rm -rf go1.12.9.linux-amd64.tar.gz \
    && mkdir gopath \
    && mkdir gopath/src \
    && mkdir gopath/bin \
    && mkdir gopath/pkg


ENV GOROOT=/usr/local/go \
    PATH=$PATH:/usr/local/go/bin \
    GOPATH=/gopath
ENV PATH=$PATH:/gopath/bin

# gcc, solc
RUN apt install -y build-essential software-properties-common \
    && add-apt-repository ppa:ethereum/ethereum \
    && apt update \
    && apt install -y solc

# git
RUN apt install -y git

# install openssl
RUN wget http://digiccy.liuyihua.com/openssl-1.1.1.tar.gz
RUN tar -xzf openssl-1.1.1.tar.gz \
    && rm -rf openssl-1.1.1.tar.gz \
    && apt-get install build-essential -y \
    && cd openssl-1.1.1 \
    && ./config --prefix=/usr/local/openssl --openssldir=/usr/local/openssl --shared \
    && make \
    && make install \
    && cd .. \
    && rm -rf openssl-1.1.1

ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/openssl/lib

# copy code
COPY . gopath/src/madledger

# Enable GO MOD, set proxy to bypass GFW
# Not that if you are using private github repositories, you may need to set GOPRIVATE env var.
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn"

WORKDIR $GOPATH/src/madledger

# Download GO dependencies
RUN go mod download

RUN make
