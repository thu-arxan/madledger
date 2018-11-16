# test the order of blocks
rm -rf $GOPATH/src/madledger/peer/channel/log.txt
rm -rf $GOPATH/src/madledger/peer/channel/.data
go test -v madledger/peer/channel -count=1 -s 1 > $GOPATH/src/madledger/peer/channel/log.txt
go test madledger/peer/channel -s 2 -count=1 