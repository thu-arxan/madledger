package consensus

// Config is the config of consensus
type Config struct {
	Timeout int
	MaxSize int
	Resume  bool
	Number  uint64
}
