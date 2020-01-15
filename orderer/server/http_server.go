package server

import (
	"fmt"
	"madledger/orderer/channel"
	"madledger/orderer/config"

	"github.com/gin-gonic/gin"
)

// Here defines some consts
const (
	Version = "v1"
)

// These define api for gin
const (
	ActionFetchBlock    = "fetchblock"
	ActionListChannels  = "listchannels"
	ActionCreateChannel = "createchannel"
	ActionAddTx         = "addtx"
)

// HTTPServer provide the serve of orderer
type HTTPServer struct {
	config *config.ServerConfig
	engine *gin.Engine
	cc     *channel.Coordinator
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
	// load chain config
	chainCfg, err := cfg.GetBlockChainConfig()
	if err != nil {
		return nil, err
	}
	// load consensus config
	consensusCfg, err := cfg.GetConsensusConfig()
	if err != nil {
		return nil, err
	}
	// get channel coordinator
	cc, err := channel.NewCoordinator(dbCfg.LevelDB.Path, chainCfg, consensusCfg)
	if err != nil {
		return nil, err
	}
	server.cc = cc
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
		v1.GET(ActionFetchBlock, hs.FetchBlock)
		v1.GET(ActionListChannels, hs.ListChannels)
		v1.GET(ActionCreateChannel, hs.CreateChannel)
		v1.GET(ActionAddTx, hs.AddTx)
	}
	return nil
}

// Start starts the server
func (hs *HTTPServer) Start() error {
	err := hs.cc.Start()
	if err != nil {
		return err
	}
	// TODO: TLS support not implemented
	addr := fmt.Sprintf("%s:%d", hs.config.Address, hs.config.Port-1)
	log.Infof("Start the orderer at %s", addr)
	hs.engine.Run(addr)

	return nil
}

// Stop stops the server
func (hs *HTTPServer) Stop() error {
	return nil
}
