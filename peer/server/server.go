package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"madledger/peer/config"
	"madledger/peer/orderer"
	"net"

	"google.golang.org/grpc/credentials"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "madledger/protos"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "peer", "package": "server"})
)

// Server provide the serve of peer
type Server struct {
	cfg       *config.ServerConfig
	rpcServer *grpc.Server
	cm        *ChannelManager
}

// NewServer is the constructor of server
func NewServer(cfg *config.Config) (*Server, error) {
	server := new(Server)
	var err error
	// set config of server
	server.cfg, err = cfg.GetServerConfig()
	if err != nil {
		return nil, err
	}
	// set channel manager
	server.cm, err = NewChannelManager(cfg)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func getOrdererClients(cfg *config.Config) ([]*orderer.Client, error) {
	// load orderer config
	ordererCfg, err := cfg.GetOrdererConfig()
	if err != nil {
		return nil, err
	}
	var clients = make([]*orderer.Client, len(ordererCfg.Address))
	for i := range ordererCfg.Address {
		clients[i], err = orderer.NewClient(ordererCfg.Address[i], cfg)
		if err != nil {
			return nil, err
		}
	}

	return clients, nil
}

// Start starts the server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Address, s.cfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.New("Failed to start the peer server")
	}
	log.Infof("Start the peer server at %s", addr)
	err = s.cm.start()
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	if s.cfg.TLS.Enable {
		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{*(s.cfg.TLS.Cert)},
			ClientCAs:    s.cfg.TLS.Pool,
		})
		opts = append(opts, grpc.Creds(creds))
	}
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
	s.cm.stop()
	log.Info("Succeed to stop the peer service")
	return nil
}
