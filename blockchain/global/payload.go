package global

// Payload is the payload of global chain
type Payload struct {
	ChannelID string
	Number    uint64
	Hash      []byte
}
