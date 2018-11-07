package consensus

// Config is the config of consensus
type Config struct {
	Timeout int
	MaxSize int
	Resume  bool
	Number  uint64
}

// DefaultConfig is the DefaultConfig
func DefaultConfig() Config {
	return Config{
		Timeout: 1000,
		MaxSize: 10,
		Number:  1,
		Resume:  false,
	}
}
