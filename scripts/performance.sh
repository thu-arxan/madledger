go test madledger/tests/performance -count=1
cat $GOPATH/src/madledger/tests/performance/performance.out
rm -rf $GOPATH/src/madledger/tests/performance/performance.out