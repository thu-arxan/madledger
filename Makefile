# Tool commands
GOCMD		= go
DOCKER_CMD	= docker

# MadLedger versions used in Makefile
MADLEDGER_VERSION		:= 0.0.1

# Build flags (overridable)
GO_LDFLAGS				?=
GO_TEST_FLAGS			?= $(GO_LDFLAGS)
GO_TEST_COUNT			?= 1
GO_TEST_TIMEOUT			?= 20m
GO_SYMBOL				?= 					# eg:GO_SYMBOL="-v"

# Go tools
GO_TEST 		= $(GOCMD) test -count=$(GO_TESTCOUNT) -timeout=$(GO_TEST_TIMEOUT) $(GO_SYMBOL)
GO_BUILD		= $(GOCMD) build

all: vet build

# go vet:format check, bug check
vet:
	@$(GOCMD) vet `go list ./...`

# The below include contains tests(quick start, setup, client tx, etc)
# include tests.mk

build:
	@echo "building orderer..."
	@$(GOCMD) install madledger/orderer

	@echo "building peer..."
	@$(GOCMD) install madledger/peer

	@echo "building client..."
	@$(GOCMD) install madledger/client

proto:
	@ cd protos && protoc --go_out=plugins=grpc:. *.proto

# test:
test:
	@$(GOCMD) test madledger/common/util -count=1 -cover
	@$(GOCMD) test madledger/common/event -count=1 -cover
	@$(GOCMD) test madledger/common/math -count=1 -cover
	@$(GOCMD) test madledger/common/hexutil -count=1 -cover
	@$(GOCMD) test madledger/common/crypto -count=1 -cover
	@$(GOCMD) test madledger/common/abi -count=1 -cover

	@$(GOCMD) test madledger/core -count=1 -cover

	@$(GOCMD) test madledger/protos -count=1 -cover

	@$(GOCMD) test madledger/blockchain/config -count=1 -cover

	@$(GOCMD) test madledger/executor/evm/wildevm -count=1 -cover

	@$(GOCMD) test madledger/consensus/solo -count=1 -cover
	@$(GOCMD) test madledger/consensus/raft -count=1 -cover
	@$(GOCMD) test madledger/consensus/tendermint -count=1 -cover

	@$(GOCMD) test madledger/orderer/config -count=1 -cover
	@$(GOCMD) test madledger/orderer/db -count=1 -cover
	@$(GOCMD) test madledger/orderer/server -count=1 -cover

	@$(GOCMD) test madledger/peer/db -count=1 -cover
	@$(GOCMD) test madledger/peer/config -count=1 -cover

	@echo "Next test may cost 1 minutes ..."
	@$(GOCMD) test madledger/tests -count=1 -cover

performance:
	@$(GOCMD) test madledger/tests/performance -count=1
	@cat tests/performance/performance.out
	@rm -rf tests/performance/performance.out

clean:
	@rm -rf tests/.bft	

syncevm:
	@rm -rf vendor/evm
	@cd ../evm && zip evm.zip $$(git ls-files) && unzip -d ../madledger/vendor/evm evm.zip && rm evm.zip
