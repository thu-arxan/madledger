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
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"madledger/peer/config"
	"madledger/peer/orderer"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc/credentials"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "madledger/protos"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "peer", "package": "server"})
)

// Here defines some consts
const (
	Version = "v1"
)

// These define api for gin
const (
	ActionGetTxStatus   = "gettxstatus"
	ActionListTxHistory = "listtxhistory"
	ActionGetTokenInfo  = "gettokeninfo"
)

// Server provide the serve of peer
type Server struct {
	cfg       *config.Config
	rpcServer *grpc.Server
	cm        *ChannelManager
	srv       *http.Server
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
func (s *Server) initServer(engine *gin.Engine) error {
	v1 := engine.Group(Version)
	{
		v1.POST(ActionGetTxStatus, s.GetTxStatusByHTTP)
		v1.POST(ActionListTxHistory, s.ListTxHistoryByHTTP)
		v1.POST(ActionGetTokenInfo, s.GetTokenInfoByHTTP)
	}
	return nil
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

	go func() {
		log.Infof("Start the peer server at %s", addr)
		err = s.rpcServer.Serve(lis)
		if err != nil {
			log.Error("gRPC Serve error: ", err)
			return
		}
	}()

	haddr := fmt.Sprintf("%s:%d", s.cfg.Address, s.cfg.Port-100)
	router := gin.Default()
	err = s.initServer(router)
	if err != nil {
		log.Error("Init router failed: ", err)
		return err
	}
	s.srv = &http.Server{
		Addr:    haddr,
		Handler: router,
	}
	go func() {
		// service connections
		log.Infof("Start the peer server at %s", haddr)
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return nil
}

// Stop stops the server
// TODO: The channel manager failed to stop
func (s *Server) Stop() {
	s.rpcServer.Stop()
	s.cm.stop()
	log.Info("Succeed to stop the peer service")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout after 1 second.")
	}
	log.Println("Server exiting")
}
