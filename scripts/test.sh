# run all tests

# run common test
go test madledger/common/util -count=1 -cover
go test madledger/common/event -count=1 -cover
go test madledger/common/math -count=1 -cover
go test madledger/common/hexutil -count=1 -cover
go test madledger/common/crypto -count=1 -cover
go test madledger/common/abi -count=1 -cover

# run core test
go test madledger/core/types -count=1 -cover

# run protos test
go test madledger/protos -count=1 -cover

# run blockchain test
go test madledger/blockchain/config -count=1 -cover

# run executor test
go test madledger/executor/evm -count=1 -cover

#run consensus test
go test madledger/consensus/solo -count=1 -cover
go test madledger/consensus/tendermint -count=1 -cover

# run orderer test
# rm -rf $GOPATH/src/madledger/orderer/config/.orderer.yaml
# orderer init -c $GOPATH/src/madledger/orderer/config/.orderer.yaml -p $GOPATH/src/madledger/orderer/config
go test madledger/orderer/config -count=1 -cover
# rm -rf $GOPATH/src/madledger/orderer/config/.tendermint
go test madledger/orderer/db -count=1 -cover
go test madledger/orderer/server -count=1 -cover

# run peer test
go test madledger/peer/db -count=1 -cover
go test madledger/peer/config -count=1 -cover
. $GOPATH/src/madledger/peer/channel/test.sh

# run all test
echo "Next test may cost 2 minutes ..."
go test madledger/tests -count=1 -cover