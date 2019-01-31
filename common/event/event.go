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
	// 0 means no limit
	maxWatchSize int
}

// NewWatchConfig is the constructor of WatchConfig
func NewWatchConfig(maxWatchSize int) *WatchConfig {
	return &WatchConfig{
		maxWatchSize: maxWatchSize,
	}
}

// DefaultWatchConfig return the default WatchConfig
func DefaultWatchConfig() *WatchConfig {
	return NewWatchConfig(0)
}
