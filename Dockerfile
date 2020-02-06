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

# copy code
COPY . gopath/src/madledger

