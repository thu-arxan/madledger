package event

// Result is the result of event
type Result struct {
	Err error
}

// NewResult is the constructor of Result
func NewResult(err error) *Result {
	return &Result{
		Err: err,
	}
}

// Event is a interface
type Event interface {
	// ID provide a unique id to notify
	ID() string
}
