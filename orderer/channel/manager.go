package channel

// Manager is the manager of channel
type Manager struct {
	// ID is the id of channel
	ID string
}

// NewManager is the constructor of Manager
// TODO: many things is not done yet
func NewManager(id string) (*Manager, error) {
	return &Manager{
		ID: id,
	}, nil
}
