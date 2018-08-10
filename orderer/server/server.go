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
	config    *config.ServerConfig
	rpcServer *grpc.Server
}

// NewServer is the constructor of server
func NewServer(config *config.ServerConfig) (*Server, error) {
	server := new(Server)
	server.config = config
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
	var opts []grpc.ServerOption
	s.rpcServer = grpc.NewServer(opts...)
	pb.RegisterOrdererServer(s.rpcServer, s)
	err = s.rpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

// Stop is used only for testing now,
func (s *Server) Stop() {
	s.rpcServer.Stop()
	log.Info().Msg("Succeed to stop the orderer service")
}
