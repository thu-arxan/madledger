package tendermint

// Config obtain all necessary configurations
type Config struct {
	Port       Port
	Dir        string
	P2PAddress []string
}

// Port includes all ports
type Port struct {
	// P2P prot
	P2P int
	// rpc port(tendermint client port)
	RPC int
	// App port
	App int
}
