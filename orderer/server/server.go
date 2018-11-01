package server

import (
	"fmt"
	"madledger/orderer/channel"
	"madledger/orderer/config"
	"net"

	pb "madledger/protos"

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
	cc        *channel.Coordinator
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
	return server, nil
}

// Start starts the server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Address, s.config.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("Failed to start the orderer server because %s", err.Error())
	}
	log.Infof("Start the orderer at %s", addr)
	err = s.cc.Start()
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	s.rpcServer = grpc.NewServer(opts...)
	pb.RegisterOrdererServer(s.rpcServer, s)
	err = s.rpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

// Stop will stop the rpc service and the consensus service
func (s *Server) Stop() {
	// if s.rpcServer != nil {
	s.rpcServer.Stop()
	// }
	s.cc.Stop()
	log.Info("Succeed to stop the orderer service")
}
