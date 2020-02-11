# Local variables used by makefile
# 共识协议类型, raft/solo/bft
CONSENSUS			=raft
ORIGIN_DIR 			=env_local/raft
SAMPLE_DIR			=samples

ifeq ($(CONSENSUS), raft)
ORIGIN_DIR=env_local/raft
else ifeq ($(CONSENSUS), solo)
ORIGIN_DIR=env_local/solo
else ifeq ($(CONSENSUS), bft)
ORIGIN_DIR=env_local/bft
$(error "bft is not supported now")
else
$(error "invalid CONSENSUS: {CONSENSUS=raft|solo|bft}")
endif

all: install

install:
	@make -s install

init: clean setup
	@cp -a -r $(ORIGIN_DIR)/. $(SAMPLE_DIR)
	@mv $(SAMPLE_DIR)/.clients $(SAMPLE_DIR)/clients
	@mv $(SAMPLE_DIR)/.peers $(SAMPLE_DIR)/peers
	@mv $(SAMPLE_DIR)/.orderers $(SAMPLE_DIR)/orderers

start:
	@cd $(SAMPLE_DIR) && bash start.sh

stop:
	@-kill `pidof orderer`
	@-kill `pidof peer`

setup:
	@mkdir -p $(SAMPLE_DIR)

clean:
	@rm -rf $(SAMPLE_DIR)

test:
	@cd $(SAMPLE_DIR) && bash test.sh
