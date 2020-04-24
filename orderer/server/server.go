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
	"madledger/orderer/channel"
	"madledger/orderer/config"
	pb "madledger/protos"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc/credentials"

	"github.com/gin-gonic/gin"

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
	cfg       *config.Config
	rpcServer *grpc.Server
	srv       *http.Server
	cc        *channel.Coordinator
	ln        net.Listener
	engine    *gin.Engine
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
	addr := fmt.Sprintf("%s:%d", s.cfg.Address, s.cfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("Failed to start the orderer server because %s", err.Error())
	}
	fmt.Printf("Start the orderer at %s\n", addr)
	err = s.cc.Start()
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

	s.Unlock()

	var ln net.Listener
	if s.cfg.TLS.Enable && s.cfg.TLS.Cert != nil {
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{*s.cfg.TLS.Cert},
		}

		ln, err = tls.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Address, s.cfg.Port-100), tlsConfig)
		if err != nil {
			log.Errorf("HTTPS listen failed: %v", err)
			return err
		}
	} else {
		ln, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Address, s.cfg.Port-100))
		if err != nil {
			log.Errorf("HTTP listen failed: %v", err)
			return err
		}
	}
	s.ln = ln
	go func() {
		err := s.srv.Serve(s.ln)
		fmt.Println("orderer listen at ", s.ln.Addr().String())
		if err != nil && err != http.ErrServerClosed {
			log.Error("Http Serve failed: ", err)
		}
	}()

	err = s.rpcServer.Serve(lis)

	return nil
}

// Stop will stop the rpc service and the consensus service
func (s *Server) Stop() {
	s.Lock()
	defer s.Unlock()
	// if s.rpcServer != nil {
	s.rpcServer.Stop()
	// }

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 1 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 1 seconds.")
	}
	s.ln.Close()

	s.cc.Stop()
	time.Sleep(500 * time.Millisecond)
	log.Info("Succeed to stop the orderer service")
}
