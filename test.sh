# run all tests

# run common test
go test madledger/common/util -count=1

# run blockchain test
go test madledger/blockchain/config -count=1

# run executor test
go test madledger/executor/evm -count=1

# run orderer test
orderer init -c orderer/config/.orderer.yaml
go test madledger/orderer/config -count=1
go test madledger/orderer/db -count=1

# run peer test
go test madledger/peer/config -count=1