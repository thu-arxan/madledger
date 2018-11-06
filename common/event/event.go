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

// WatchConfig is the config of Watch
type WatchConfig struct {
	Single bool
}

// NewWatchConfig is the constructor of WatchConfig
func NewWatchConfig(single bool) *WatchConfig {
	return &WatchConfig{
		Single: single,
	}
}

// DefaultWatchConfig return the default WatchConfig
func DefaultWatchConfig() *WatchConfig {
	return NewWatchConfig(false)
}
