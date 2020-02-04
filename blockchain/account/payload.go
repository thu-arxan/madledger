package account

import "madledger/common"

type Payload struct {
	ChannelID string
	account common.Account
}