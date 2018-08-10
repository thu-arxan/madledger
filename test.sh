# run all tests

# run orderer test
go test madledger/orderer/config -count=1

# run peer test
go test madledger/peer/config -count=1