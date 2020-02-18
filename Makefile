# Tool commands
GOCMD		= go
DOCKER_CMD	= docker

# MadLedger versions used in Makefile
MADLEDGER_VERSION		:= v0.0.1

# Build flags (overridable)
GO_LDFLAGS				?= -X madledger/version.GitCommit=`git rev-parse --short=8 HEAD` -X madledger/version.Version=$(MADLEDGER_VERSION)
GO_TEST_FLAGS			?= $(GO_LDFLAGS)
GO_TEST_COUNT			?= 1
GO_TEST_TIMEOUT			?= 20m
GO_SYMBOL				?= 					# eg:GO_SYMBOL="-v -race"

# Go tools
GO_TEST 		= $(GOCMD) test -parallel=1 -count=$(GO_TEST_COUNT) -timeout=$(GO_TEST_TIMEOUT) $(GO_SYMBOL)
GO_TEST_UNIT	= $(GO_TEST) -cover
GO_BUILD		= $(GOCMD) build

# Local variables used by makefile
PROJECT_NAME           := madledger
ARCH                   := $(shell uname -m)
OS_NAME                := $(shell uname -s)

# Test Packages
# UNIT_PACKAGES	=	madledger/common/util \
# 					madledger/common/event \
# 					madledger/common/math \
# 					madledger/common/hexutil \
# 					madledger/common/crypto \
# 					madledger/common/abi \
# 					madledger/core \
# 					madledger/protos \
# 					madledger/blockchain/config \

PACKAGES=$(shell go list ./...)

all: vet install

# go vet:format check, bug check
vet:
	@$(GOCMD) vet `go list ./...`

# The below include contains tests(quick start, setup, client tx, etc)
# include tests.mk

unittest:
	@$(GO_TEST) $(PACKAGES)

install:
	@echo "install orderer..."
	@$(GOCMD) install madledger/orderer

	@echo "install peer..."
	@$(GOCMD) install madledger/peer

	@echo "install client..."
	@$(GOCMD) install madledger/client

proto:
	@ cd protos && protoc --go_out=plugins=grpc:. *.proto
	@ cd consensus/raft/protos && protoc --go_out=plugins=grpc:. *.proto

# test:
test:
	@$(GO_TEST_UNIT) madledger/common/util
	@$(GO_TEST_UNIT) madledger/common/event
	@$(GO_TEST_UNIT) madledger/common/math
	@$(GO_TEST_UNIT) madledger/common/crypto
	@$(GO_TEST_UNIT) madledger/common/abi

	@$(GO_TEST_UNIT) madledger/core

	@$(GO_TEST_UNIT) madledger/protos

	@$(GO_TEST_UNIT) madledger/blockchain/config

	@$(GO_TEST_UNIT) madledger/consensus/solo
	@$(GO_TEST_UNIT) madledger/consensus/raft
	@$(GO_TEST_UNIT) madledger/consensus/raft/eraft
	@$(GO_TEST_UNIT) madledger/consensus/tendermint

	@$(GO_TEST_UNIT) madledger/orderer/config
	@$(GO_TEST_UNIT) madledger/orderer/db
	@$(GO_TEST_UNIT) madledger/orderer/server

	@$(GO_TEST_UNIT) madledger/peer/db
	@$(GO_TEST_UNIT) madledger/peer/config

	@echo "Next test may cost 1 minutes ..."
	@$(GO_TEST_UNIT) madledger/tests

performance:
	@$(GO_TEST) madledger/tests/performance
	@cat tests/performance/performance.out
	@rm -rf tests/performance/performance.out

docker:
	@docker build -t madledger:alpha .

clean:
	@rm -rf tests/.bft
	@cd tests/performance/raft && rm -rf .clients .orderer .peer
	@cd tests/performance/solo && rm -rf .clients .orderer .peer
	@cd tests/performance/bft && rm -rf .clients .orderer .peer

syncevm:
	@rm -rf vendor/evm
	@cd ../evm && zip evm.zip $$(git ls-files) && unzip -d ../madledger/vendor/evm evm.zip && rm evm.zip

raft:
	@$(GO_TEST_UNIT) madledger/consensus/raft -race
