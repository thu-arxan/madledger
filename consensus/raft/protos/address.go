package protos

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ChainToRaft return the raft service address according to the chain address
func ChainToRaft(chainAddress string) string {
	url, port, err := ParseChainAddress(chainAddress)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s:%d", url, port+1)
}

// ChainToERaft return the eraft address according to the chain address
func ChainToERaft(chainAddress string) string {
	url, port, err := ParseChainAddress(chainAddress)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s:%d", url, port+2)
}

// ERaftToRaft return the raft service address according to the etcd raft address
func ERaftToRaft(eraftAddress string) string {
	url, port, err := ParseERaftAddress(eraftAddress)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s:%d", url, port-1)
}

// ERaftToChain return the chain address according to the etcd raft address
func ERaftToChain(eraftAddress string) string {
	url, port, err := ParseERaftAddress(eraftAddress)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s:%d", url, port-2)
}

// ParseERaftAddress parse etcd raft address
func ParseERaftAddress(eraftAddress string) (string, int, error) {
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

// ParseChainAddress parse blockchain address
func ParseChainAddress(chainAddress string) (string, int, error) {
	return parseAddress(chainAddress, 1024)
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
