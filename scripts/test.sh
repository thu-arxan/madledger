# run all tests

# run common test
go test madledger/common/util -count=1
go test madledger/common/event -count=1
go test madledger/common/math -count=1
go test madledger/common/hexutil -count=1
go test madledger/common/crypto -count=1
go test madledger/common/abi -count=1

# run core test
go test madledger/core/types -count=1

# run protos test
go test madledger/protos -count=1

# run blockchain test
go test madledger/blockchain/config -count=1

# run executor test
go test madledger/executor/evm -count=1

#run consensus test
go test madledger/consensus/solo -count=1

# run orderer test
# rm -rf $GOPATH/src/madledger/orderer/config/.orderer.yaml
# orderer init -c $GOPATH/src/madledger/orderer/config/.orderer.yaml -p $GOPATH/src/madledger/orderer/config
go test madledger/orderer/config -count=1
# rm -rf $GOPATH/src/madledger/orderer/config/.tendermint
go test madledger/orderer/db -count=1
go test madledger/orderer/server -count=1

# run peer test
go test madledger/peer/db -count=1
go test madledger/peer/config -count=1
. $GOPATH/src/madledger/peer/channel/test.sh

# run all test
echo "Next test may need nearly 40 seconds..."
go test madledger/tests -count=1