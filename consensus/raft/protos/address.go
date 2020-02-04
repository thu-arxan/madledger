package protos

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ERaftToRaft return the raft service address according to the etcd raft address
func ERaftToRaft(eraftAddress string) string {
	url, port, err := ParseRaftAddress(eraftAddress)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s:%d", url, port-1)
}

// RaftToERaft return the eraft service address according the raft address
func RaftToERaft(raftAddress string) string {
	url, port, err := ParseRaftAddress(raftAddress)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s:%d", url, port+1)
}

// ParseRaftAddress parse etcd raft address
func ParseRaftAddress(eraftAddress string) (string, int, error) {
	eraftAddress = strings.Replace(eraftAddress, " ", "", -1)
	s := strings.Split(eraftAddress, ":")
	if len(s) != 2 {
		return "", 0, errors.New("The length is not 2")
	}
	port, err := strconv.Atoi(s[1])
	if err != nil || port < 2 {
		return "", 0, errors.New("Failed to parse the port")
	}
	return s[0], port, nil
}

func parseAddress(address string, minPort int) (string, int, error) {
	address = strings.Replace(address, " ", "", -1)
	s := strings.Split(address, ":")
	if len(s) != 2 {
		return "", 0, errors.New("The length is not 2")
	}
	port, err := strconv.Atoi(s[1])
	if err != nil || port < minPort {
		return "", 0, errors.New("Failed to parse the port")
	}
	return s[0], port, nil

}
