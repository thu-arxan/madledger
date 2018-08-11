package db

// DB is the interface of db
type DB interface {
	// ListChannel list all channels
	ListChannel() []string
	// AddChannel add a channel
	AddChannel(id string) error
}
