# Tool commands
GOCMD		= go
DOCKER_CMD	= docker

# MadLedger versions used in Makefile
MADLEDGER_VERSION		:=v0.0.1

# Test flags
CONSENSUS	=solo

# Build flags (overridable)
GO_LDFLAGS				?= -X madledger/version.GitCommit=`git rev-parse --short=8 HEAD` -X madledger/version.Version=$(MADLEDGER_VERSION)
GO_TEST_FLAGS			?= $(GO_LDFLAGS) -X madledger/tests/performance.consensus=$(CONSENSUS)
GO_TEST_COUNT			?= 1
GO_TEST_TIMEOUT			?= 20m
GO_SYMBOL				?= 					# eg:GO_SYMBOL="-v -race"
# database build tag, use rocksdb or leveldb
DB_TAG					?=leveldb
BUILD_TAGS				=

# check db tag
ifeq ($(DB_TAG), rocksdb)
BUILD_TAGS+=rocksdb
else ifeq ($(DB_TAG), leveldb)
BUILD_TAGS+=leveldb
else
$(error "invalid DB_TAG: {DB_TAG=rocksdb|leveldb}")
endif

# Go tools
GO_TEST 		= $(GOCMD) test -tags "$(BUILD_TAGS)" -parallel=1 -count=$(GO_TEST_COUNT) -timeout=$(GO_TEST_TIMEOUT) $(GO_SYMBOL) -ldflags "$(GO_TEST_FLAGS)"
GO_TEST_UNIT	= $(GO_TEST) -cover -race
GO_BUILD		= $(GOCMD) build -tags "$(BUILD_TAGS)"
GO_INSTALL		= $(GOCMD) install -tags "$(BUILD_TAGS)"

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
# 					madledger/core \
# 					madledger/protos \
# 					madledger/blockchain/config \

PACKAGES=$(shell go list ./...)

all: vet install

# go vet:format check, bug check
vet:
	@go vet `go list ./...`

# The below include contains tests(quick start, setup, client tx, etc)
# include tests.mk

unittest:
	@$(GO_TEST) $(PACKAGES)

install:
	@echo "install orderer..."
	@$(GO_INSTALL) madledger/orderer

	@echo "install peer..."
	@$(GO_INSTALL) madledger/peer

	@echo "install client..."
	@$(GO_INSTALL) madledger/client

proto:
	@ cd protos && protoc --go_out=plugins=grpc:. *.proto
	@ cd consensus/raft/protos && protoc --go_out=plugins=grpc:. *.proto

# test:
test:
	# @$(GO_TEST_UNIT) madledger/common/util
	# @$(GO_TEST_UNIT) madledger/common/event
	# @$(GO_TEST_UNIT) madledger/common/math
	# @$(GO_TEST_UNIT) madledger/common/crypto

	# @$(GO_TEST_UNIT) madledger/core

	# @$(GO_TEST_UNIT) madledger/protos

	# @$(GO_TEST_UNIT) madledger/blockchain/config

	# @$(GO_TEST_UNIT) madledger/consensus/solo
	# @$(GO_TEST_UNIT) madledger/consensus/raft
	# @$(GO_TEST_UNIT) madledger/consensus/raft/eraft
	# @$(GO_TEST_UNIT) madledger/consensus/tendermint

	# @$(GO_TEST_UNIT) madledger/orderer/config
	# @$(GO_TEST_UNIT) madledger/orderer/db
	# @$(GO_TEST_UNIT) madledger/orderer/server

	# @$(GO_TEST_UNIT) madledger/peer/db
	# @$(GO_TEST_UNIT) madledger/peer/config

	@echo "Next test may cost 1 minutes ..."
	@$(GO_TEST_UNIT) madledger/tests -v

performance:
	@$(GO_TEST) madledger/tests/performance
	@cat tests/performance/performance.out
	@rm -rf tests/performance/performance.out

docker:
	@docker build -t madledger:alpha .

httptest:
	@-kill -9 `pidof orderer`
	@-kill -9 `pidof peer`
	@go test madledger/tests -v -count=1 > log.out
	@tail log.out
	
clean:
	@rm -rf tests/.bft
	@cd tests/performance/raft && rm -rf .clients .orderer .peer
	@cd tests/performance/solo && rm -rf .clients .orderer .peer
	@cd tests/performance/bft && rm -rf .clients .orderer .peer

raft:
	@$(GO_TEST_UNIT) madledger/consensus/raft
