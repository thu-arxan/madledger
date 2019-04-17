FROM ubuntu:18.04

# go environment
RUN apt-get update \
    && apt install -y wget \
    && wget https://dl.google.com/go/go1.11.linux-amd64.tar.gz \
    && tar -xzf go1.11.linux-amd64.tar.gz -C /usr/local \
    && rm -rf go1.11.linux-amd64.tar.gz \
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

# build
RUN . gopath/src/madledger/scripts/build.sh

