// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

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
	cfg       *config.Config
	rpcServer *grpc.Server
	cm        *ChannelManager
}

// NewServer is the constructor of server
func NewServer(cfg *config.Config) (*Server, error) {
	server := new(Server)
	server.cfg = cfg
	var err error
	// set channel manager
	server.cm, err = NewChannelManager(cfg)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func getOrdererClients(cfg *config.Config) ([]*orderer.Client, error) {
	var clients = make([]*orderer.Client, len(cfg.Orderer.Address))
	var err error
	for i := range cfg.Orderer.Address {
		clients[i], err = orderer.NewClient(cfg.Orderer.Address[i], cfg)
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
