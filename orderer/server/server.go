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
	ActionFetchBlock    = "fetchblock"
	ActionListChannels  = "listchannels"
	ActionCreateChannel = "createchannel"
	ActionAddTx         = "addtx"
)

// Server provide the serve of orderer
type Server struct {
	sync.RWMutex
	config    *config.ServerConfig
	rpcServer *grpc.Server
	srv       *http.Server
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
func (s *Server) initServer(engine *gin.Engine) error {
	v1 := engine.Group(Version)
	{
		v1.POST(ActionFetchBlock, s.FetchBlockByHTTP)
		v1.POST(ActionListChannels, s.ListChannelsByHTTP)
		v1.POST(ActionCreateChannel, s.CreateChannelByHTTP)
		v1.POST(ActionAddTx, s.AddTxByHTTP)
	}
	return nil
}

// Start starts the server
func (s *Server) Start() error {
	s.Lock()
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

	s.Unlock()

	err = s.rpcServer.Serve(lis)
	
	// TODO: TLS support not implemented
	haddr := fmt.Sprintf("%s:%d", s.config.Address, s.config.Port-100)
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
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return nil
}

// Stop will stop the rpc service and the consensus service
func (s *Server) Stop() {
	s.Lock()
	defer s.Unlock()
	// if s.rpcServer != nil {
	s.rpcServer.Stop()
	// }
	s.cc.Stop()
	time.Sleep(500 * time.Millisecond)
	log.Info("Succeed to stop the orderer service")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
