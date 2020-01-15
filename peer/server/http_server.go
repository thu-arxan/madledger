package server

import (
	"fmt"
	"madledger/peer/config"
	"madledger/peer/orderer"

	"github.com/gin-gonic/gin"
)

// Here defines some consts
const (
	Version = "v1"
)

// These define api for gin
const (
	ActionGetTxStatus   = "gettxstatus"
	ActionListTxHistory = "listtxhistory"
)

// HTTPServer provide the serve of orderer
type HTTPServer struct {
	config         *config.ServerConfig
	engine         *gin.Engine
	ChannelManager *ChannelManager
	ordererClients []*orderer.Client
}

// NewHTTPServer is the constructor of HTTPServer
func NewHTTPServer(cfg *config.Config) (*HTTPServer, error) {
	server := new(HTTPServer)
	// set config of server
	serverCfg, err := cfg.GetServerConfig()
	if err != nil {
		return nil, err
	}
	server.config = serverCfg
	// load db config
	dbCfg, err := cfg.GetDBConfig()
	if err != nil {
		return nil, err
	}
	// load orderer config
	ordererClients, err := getOrdererClients(cfg)
	if err != nil {
		return nil, err
	}
	// load chain config
	chainCfg, err := cfg.GetBlockChainConfig()
	if err != nil {
		return nil, err
	}
	// load identity
	identity, err := cfg.GetIdentity()
	if err != nil {
		return nil, err
	}
	channelManager, err := NewChannelManager(dbCfg.LevelDB.Dir, identity, chainCfg, ordererClients)
	if err != nil {
		return nil, err
	}
	server.ChannelManager = channelManager
	server.ordererClients = ordererClients
	if err != nil {
		return nil, err
	}
	server.engine = gin.New()
	err = server.initServer()
	if err != nil {
		log.Error("Init router failed: ", err)
		return nil, err
	}

	return server, nil
}

func (hs *HTTPServer) initServer() error {
	v1 := hs.engine.Group(Version)
	{
		v1.GET(ActionGetTxStatus, hs.GetTxStatus)
		v1.GET(ActionListTxHistory, hs.ListTxHistory)
	}
	return nil
}

// Start starts the server
func (hs *HTTPServer) Start() error {
	err := hs.ChannelManager.start()
	if err != nil {
		return err
	}
	// TODO: TLS support not implemented
	addr := fmt.Sprintf("%s:%d", hs.config.Address, hs.config.Port-1)
	log.Infof("Start the peer server at %s", addr)
	log.Fatal(hs.engine.Run(addr))

	return nil
}

// Stop stops the server
func (hs *HTTPServer) Stop() error {
	return nil
}
