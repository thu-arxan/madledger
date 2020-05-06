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
	"madledger/common/util"
	"madledger/orderer/channel"
	"madledger/orderer/config"
	pb "madledger/protos"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"github.com/gin-gonic/gin"

	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "orderer", "package": "server"})
)

// Here defines some consts
const (
	Version = "v1"
)

// These define api for gin
const (
	// ActionFetchBlock     = "fetchblock"
	ActionListChannels   = "listchannels"
	ActionCreateChannel  = "createchannel"
	ActionAddTx          = "addtx"
	ActionGetAccountInfo = "getaccountinfo"
)

// Server provide the serve of orderer
type Server struct {
	sync.RWMutex
	cfg          *config.Config
	cc           *channel.Coordinator
	rpcServer    *grpc.Server
	srv          *http.Server
	rpcWebServer *grpcweb.WrappedGrpcServer
	engine       *gin.Engine
}

// NewServer is the constructor of server
func NewServer(cfg *config.Config) (*Server, error) {
	server := new(Server)
	server.cfg = cfg
	// load db config
	dbCfg, err := cfg.GetDBConfig()
	if err != nil {
		return nil, err
	}
	// load consensus config
	consensusCfg, err := cfg.GetConsensusConfig()
	if err != nil {
		return nil, err
	}
	// get channel coordinator
	cc, err := channel.NewCoordinator(dbCfg.LevelDB.Path, cfg, consensusCfg)
	if err != nil {
		return nil, err
	}
	server.cc = cc

	/*
		server.engine = gin.New()
		server.engine.Use(gin.Recovery())
		server.initServer(server.engine)
		server.srv = &http.Server{
			Handler:      server.engine,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		}*/

	return server, nil
}
func (s *Server) initServer(engine *gin.Engine) error {
	v1 := engine.Group(Version)
	{
		v1.POST(ActionListChannels, s.ListChannelsByHTTP)
		v1.POST(ActionCreateChannel, s.CreateChannelByHTTP)
		v1.POST(ActionAddTx, s.AddTxByHTTP)
		v1.POST(ActionGetAccountInfo, s.GetAccountInfoByHTTP)
	}
	return nil
}

// Start starts the server
func (s *Server) Start() error {
	s.Lock()
	log.Infof("Server start...")
	util.MountLogger()

	err := s.cc.Start()
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
	pb.RegisterOrdererServer(s.rpcServer, s)
	go func() { // Start Native GRPC Server at Address:Port
		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Address, s.cfg.Port))
		if err != nil {
			log.Fatalf("Failed to start the orderer server because %s", err.Error())
		}
		fmt.Printf("Start the orderer at %s:%d\n", s.cfg.Address, s.cfg.Port)

		log.Infof("rpcServer serve at %s:%d", s.cfg.Address, s.cfg.Port)
		err = s.rpcServer.Serve(lis)
	}()

	s.rpcWebServer = grpcweb.WrapServer(s.rpcServer)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		//These header are intentionally add to solve browsers' CORS problem.
		resp.Header().Add("Access-Control-Allow-Origin", "*")
		resp.Header().Add("Access-Control-Allow-Headers", "x-grpc-web, content-type")
		grpclog.Infof("Handle grpc request : %v", req)
		s.rpcWebServer.ServeHTTP(resp, req)
	}

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.Port+11),
		Handler: http.HandlerFunc(handler),
	}
	s.srv = &httpServer

	go func() { // Start GRPC-WEB Server at Address:(Port+11)
		if s.cfg.TLS.Enable {
			httpServer.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{*(s.cfg.TLS.Cert)},
				ClientCAs:    s.cfg.TLS.Pool,
			}
			grpclog.Infof("Start tls rpc-web server at %d", s.cfg.Port+11)
			grpclog.Infof("tls config : ca = %s, key = %s", s.cfg.TLS.RawCert, s.cfg.TLS.Key)
			if err := httpServer.ListenAndServeTLS(s.cfg.TLS.RawCert, s.cfg.TLS.Key); err != nil {
				if err.Error() == "http: Server closed" {
					grpclog.Infof("grpc-web server exit: %v")
				} else {
					grpclog.Fatalf("failed starting rpc-web server: %v", err.Error())
				}
				grpclog.SetLoggerV2(nil)
				log.Info("grpcLogger shutdown...")
			}
		} else {
			grpclog.Infof("Start insecure rpc-web server at %d", s.cfg.Port+11)
			if err := httpServer.ListenAndServe(); err != nil {
				if err.Error() == "http: Server closed" {
					grpclog.Infof("grpc-web server exit: %v", err)
				} else {
					grpclog.Fatalf("failed starting rpc-web server: %v", err)
				}
			}
		}
	}()
	// MHY ADD END

	s.Unlock()

	time.Sleep(100 * time.Millisecond)
	return nil
}

// Stop will stop the rpc service and the consensus service
func (s *Server) Stop() {
	log.Info("Stop Server")
	s.Lock()
	defer s.Unlock()

	s.rpcServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 1 seconds.")
	}

	s.cc.Stop()
	time.Sleep(500 * time.Millisecond)
	log.Info("Succeed to stop the orderer service")
}
