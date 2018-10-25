FROM ubuntu:18.04

RUN apt-get update \
    && apt install -y wget \
    && wget https://dl.google.com/go/go1.11.linux-amd64.tar.gz \
    && tar -xzf go1.11.linux-amd64.tar.gz -C /usr/local \
    && rm -rf go1.11.linux-amd64.tar.gz \
    && mkdir gopath \
    && mkdir gopath/src \
    && mkdir gopath/bin \
    && mkdir gopath/pkg \
    && export GOROOT=/usr/local/go \
    && export PATH=$PATH:$GOROOT/bin \
    && export GOPATH=/home/gopath 