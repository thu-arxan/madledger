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
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"madledger/orderer/channel"
	"madledger/orderer/config"
	pb "madledger/protos"
	"net"
	"net/http"
	"sync"
	"time"

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
	config       *config.ServerConfig
	rpcServer    *grpc.Server
	rpcWebServer *grpcweb.WrappedGrpcServer
	srv          *http.Server
	cc           *channel.Coordinator
	ln           net.Listener
	engine       *gin.Engine
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

// Wrapper to provide log for GRPC
type GrpcLogger struct {
	*logrus.Entry
}

func NewGrpcLogger() *GrpcLogger {
	return &GrpcLogger{
		logrus.WithFields(logrus.Fields{"app": "orderer", "package": "server/grpc"}),
	}
}

func (_ *GrpcLogger) V(_ int) bool {
	return false
}

// Start starts the server
func (s *Server) Start() error {
	log.Infof("Server start...")
	grpclog.SetLoggerV2(NewGrpcLogger()) // Export GRPC's log

	s.Lock()
	err := s.cc.Start()
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
	pb.RegisterOrdererServer(s.rpcServer, s)

	go func() { // Start Native GRPC Server at Address:Port
		lis, err := net.Listen("tcp",  fmt.Sprintf("%s:%d", s.config.Address, s.config.Port))
		if err != nil {
			log.Fatalf("Failed to start the orderer server because %s", err.Error())
		}
		fmt.Printf("Start the orderer at %s:%d\n", s.config.Address, s.config.Port)

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
				//RootCAs:    s.config.TLS.Pool,
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
					grpclog.Infof("grpc-web server exit: %v")
				} else {
					grpclog.Fatalf("failed starting rpc-web server: %v", err)
				}
			}
		}
	}()
	// MHY ADD END


	s.Unlock()

	/*
	var ln net.Listener
	go func() {
		if s.config.TLS.Enable && s.config.TLS.Cert != nil {
			tlsConfig := &tls.Config{
				Certificates: []tls.Certificate{*s.config.TLS.Cert},
			}

			ln, err = tls.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Address, s.config.Port-100), tlsConfig)
			if err != nil {
				log.Errorf("HTTPS listen failed: %v", err)
			}
		} else {
			ln, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Address, s.config.Port-100))
			if err != nil {
				log.Errorf("HTTP listen failed: %v", err)
			}
		}
	}()
	s.ln = ln
	go func() {
		err := s.srv.Serve(s.ln)
		fmt.Println("orderer listen at ", s.ln.Addr().String())
		if err != nil && err != http.ErrServerClosed {
			log.Error("Http Serve failed: ", err)
		}
	}()
	*/


	// TODO: TLS support not implemented
	// haddr := fmt.Sprintf("%s:%d", s.config.Address, s.config.Port-100)
	// router := gin.Default()
	// err = s.initServer(router)
	// if err != nil {
	// 	log.Error("Init router failed: ", err)
	// 	return err
	// }
	// s.srv = &http.Server{
	// 	Addr:    haddr,
	// 	Handler: router,
	// }
	// go func() {
	// 	// service connections
	// 	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("listen: %s\n", err)
	// 	}
	// }()

	return nil
}

// Stop will stop the rpc service and the consensus service
func (s *Server) Stop() {
	s.Lock()
	defer s.Unlock()
	s.rpcServer.Stop()

	s.srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		log.Println("timeout of 1 seconds.")
	}

	s.cc.Stop()
	time.Sleep(500 * time.Millisecond)
	log.Info("Succeed to stop the orderer service")
}
