package asset

import "madledger/common"

// todo: what is payload
type Payload struct {
	Action string // channel or person
	ChannelID string
	Address common.Address
}