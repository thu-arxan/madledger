package raft

import (
	"strconv"
	"strings"
)

// Error is the type of raft error
type Error int

// Here defines some kinds of errors in raft
const (
	NotLeader Error = iota
	RemovedNode
	TxInPool
	Unknown
)

// Here defines error msg for check
const (
	NotLeaderMsg   = "Please send to leader"
	RemovedNodeMsg = "I've been removed from cluster"
	TxInPoolMsg    = "Transaction is aleardy in the pool"
)

// GetError returns error type of raft error
func GetError(err error) Error {
	if err != nil {
		return Unknown
	}
	e := err.Error()
	if strings.Contains(e, NotLeaderMsg) {
		return NotLeader
	}
	if strings.Contains(e, RemovedNodeMsg) {
		return RemovedNode
	}
	if strings.Contains(e, TxInPoolMsg) {
		return TxInPool
	}
	return Unknown
}

// GetLeader parses leader id from raft error
func GetLeader(e error) uint64 {
	if GetError(e) != NotLeader {
		return 0
	}
	idstr := strings.Replace(e.Error(), "rpc error: code = Unknown desc = Please send to leader ", "", -1)
	if len(idstr) == 0 {
		return 0
	}
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		log.Debugf("failed to parse leader id: %v", err)
	}
	return id
}
