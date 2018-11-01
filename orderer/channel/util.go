package channel

import (
	"regexp"
)

func isLegalChannelName(channelID string) bool {
	if m, _ := regexp.MatchString("^[a-z0-9]{1,32}$", channelID); !m {
		return false
	}
	return true
}
