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
	"fmt"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc/grpclog"
	"madledger/common/util"
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
	glog = logrus.WithFields(logrus.Fields{"app": "peer", "package": "server/grpc"})
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
	ActionGetBlock      = "getblock"
)

// Server provide the serve of peer
type Server struct {
	config       *config.Config
	rpcServer    *grpc.Server
	rpcWebServer *grpcweb.WrappedGrpcServer
	cm           *ChannelManager
	srv          *http.Server
	engine       *gin.Engine
}

// NewServer is the constructor of server
func NewServer(cfg *config.Config) (*Server, error) {
	server := new(Server)
	server.config = cfg
	var err error
	// set channel manager
	server.cm, err = NewChannelManager(cfg)
	if err != nil {
		return nil, err
	}

	server.engine = gin.New()
	server.engine.Use(gin.Recovery())
	server.initServer(server.engine)
	server.srv = &http.Server{
		Handler:      server.engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
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
		v1.POST(ActionGetBlock, s.GetBlockByHTTP)

	}
	return nil
}

// Start starts the server
func (s *Server) Start() error {
	log.Infof("Server start...")
	grpclog.SetLoggerV2(&util.GrpcLogger{Entry: glog}) // Export GRPC's log

	err := s.cm.start()
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	if s.config.TLS.Enable {
		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{*(s.config.TLS.Cert)},
			ClientCAs:    s.config.TLS.Pool,
		})
		opts = append(opts, grpc.Creds(creds))
	}
	s.rpcServer = grpc.NewServer(opts...)
	pb.RegisterPeerServer(s.rpcServer, s)

	go func() { // Start Native GRPC Server at Address:Port
		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Address, s.config.Port))
		if err != nil {
			log.Fatalf("Failed to start the peer server because %s", err.Error())
		}
		log.Infof("rpcServer serve at %s:%d", s.config.Address, s.config.Port)
		err = s.rpcServer.Serve(lis)
	}()

	s.rpcWebServer = grpcweb.WrapServer(s.rpcServer)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		//These header are intentionally add to solve browsers' CORS problem.
		resp.Header().Add("Access-Control-Allow-Origin","*")
		resp.Header().Add("Access-Control-Allow-Headers", "x-grpc-web, content-type")
		grpclog.Infof("Handle grpc request : %v", req)
		s.rpcWebServer.ServeHTTP(resp, req)
	}

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port + 11),
		Handler: http.HandlerFunc(handler),
	}
	s.srv = &httpServer

	go func() { // Start GRPC-WEB Server at Address:(Port+11)
		if s.config.TLS.Enable {
			httpServer.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{*(s.config.TLS.Cert)},
				ClientCAs:    s.config.TLS.Pool,
			}
			grpclog.Infof("Start tls rpc-web server at %d", s.config.Port + 11)
			grpclog.Infof("tls config : ca = %s, key = %s", s.config.TLS.RawCert, s.config.TLS.Key)
			if err := httpServer.ListenAndServeTLS(s.config.TLS.RawCert, s.config.TLS.Key); err != nil {
				if err.Error() == "http: Server closed" {
					grpclog.Infof("grpc-web server exit: %v")
				} else {
					grpclog.Fatalf("failed starting rpc-web server: %v", err.Error())
				}
			}
		} else {
			grpclog.Infof("Start insecure rpc-web server at %d", s.config.Port + 11)
			if err := httpServer.ListenAndServe(); err != nil {
				if err.Error() == "http: Server closed" {
					grpclog.Infof("grpc-web server exit: %v", err)
				} else {
					grpclog.Fatalf("failed starting rpc-web server: %v", err)
				}
			}
		}
	}()

	return nil
}

// Stop stops the server
// TODO: The channel manager failed to stop
func (s *Server) Stop() {
	log.Info("Stop Server")
	s.rpcServer.Stop()

	log.Info("Succeed to stop the peer service")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// catching ctx.Done(). timeout of 1 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout after 1 second.")
	}

	s.cm.stop()
	log.Println("Server exiting")
}
