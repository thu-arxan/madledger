package config

// Payload config a channel
type Payload struct {
	ChannelID string
	Profile   Profile
	Version   int32
}

// Profile is the profile in payload
type Profile struct {
	Public bool
}
