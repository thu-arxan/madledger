# building all modules

echo "building orderer..."
go install madledger/orderer

echo "building peer..."
go install madledger/peer