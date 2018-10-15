package server

import (
	"errors"
	"fmt"
	"madledger/peer/config"
	"madledger/peer/orderer"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "madledger/protos"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "peer", "package": "server"})
)

// Server provide the serve of orderer
type Server struct {
	config         *config.ServerConfig
	rpcServer      *grpc.Server
	ChannelManager *ChannelManager
	ordererClient  *orderer.Client
}

// NewServer is the constructor of server
func NewServer(cfg *config.Config) (*Server, error) {
	server := new(Server)
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
	ordererCfg, err := cfg.GetOrdererConfig()
	if err != nil {
		return nil, err
	}
	ordererClient, err := orderer.NewClient(ordererCfg.Address[0])
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
	channelManager, err := NewChannelManager(dbCfg.LevelDB.Dir, identity, chainCfg, ordererClient)
	if err != nil {
		return nil, err
	}
	server.ChannelManager = channelManager
	server.ordererClient = ordererClient

	return server, nil
}

// Start starts the server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Address, s.config.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.New("Failed to start the orderer")
	}
	log.Infof("Start the orderer at %s", addr)
	err = s.ChannelManager.start()
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	s.rpcServer = grpc.NewServer(opts...)
	pb.RegisterPeerServer(s.rpcServer, s)
	err = s.rpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

// Stop stops the server
// TODO: The channel manager failed to stop
func (s *Server) Stop() error {
	s.rpcServer.Stop()
	// s.ChannelManager.stop()
	log.Info("Succeed to stop the orderer service")
	return nil
}
