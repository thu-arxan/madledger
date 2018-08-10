# run all tests

# run orderer test
orderer init -c orderer/config/.orderer.yaml
go test madledger/orderer/config -count=1

# run peer test
go test madledger/peer/config -count=1