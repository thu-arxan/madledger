package server

import (
	"errors"
	"fmt"
	"madledger/orderer/config"
	"net"

	pb "madledger/protos"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// Server provide the serve of orderer
type Server struct {
	config         *config.ServerConfig
	rpcServer      *grpc.Server
	ChannelManager *ChannelManager
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
	dbCfg, err := cfg.GetDBConfig()
	if err != nil {
		return nil, err
	}
	// load chain config
	chainCfg, err := cfg.GetBlockChainConfig()
	if err != nil {
		return nil, err
	}
	// get channel manager
	channelManager, err := NewChannelManager(dbCfg.LevelDB.Dir, chainCfg)
	if err != nil {
		return nil, err
	}
	server.ChannelManager = channelManager
	return server, nil
}

// Start starts the server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Address, s.config.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.New("Failed to start the orderer")
	}
	log.Info().Msgf("Start the orderer at %s", addr)
	err = s.ChannelManager.start()
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
	s.ChannelManager.stop()
	log.Info().Msg("Succeed to stop the orderer service")
}
