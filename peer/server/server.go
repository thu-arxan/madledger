package server

import (
	"madledger/peer/config"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "orderer", "package": "server"})
)

// Server provide the serve of orderer
type Server struct {
	config    *config.ServerConfig
	rpcServer *grpc.Server
	// ChannelManager *ChannelManager
}

// NewServer is the constructor of server
// todo: many thing need to be done
func NewServer(cfg *config.Config) (*Server, error) {
	server := new(Server)
	// set config of server
	serverCfg, err := cfg.GetServerConfig()
	if err != nil {
		return nil, err
	}
	server.config = serverCfg
	// load db config
	// dbCfg, err := cfg.GetDBConfig()
	// if err != nil {
	// 	return nil, err
	// }
	// // load chain config
	// ordererCfg, err := cfg.GetOrdererConfig()
	// if err != nil {
	// 	return nil, err
	// }
	return server, nil
}
