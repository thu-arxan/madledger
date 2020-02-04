package consensus

// Status is the status of consensus service
type Status int

// Here defines some kind of status
const (
	Stopped Status = iota
	OnStopped
	Started
	OnStarted
)
