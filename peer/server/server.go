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
)

// Server provide the serve of peer
type Server struct {
	config         *config.ServerConfig
	rpcServer      *grpc.Server
	srv            *http.Server
	ChannelManager *ChannelManager
	ordererClients []*orderer.Client
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
	// load orderer config
	ordererClients, err := getOrdererClients(cfg)
	if err != nil {
		return nil, err
	}
	// load chain config
	chainCfg, err := cfg.GetBlockChainConfig()
	if err != nil {
		return nil, err
	}
	// load identity
	identity, err := cfg.GetIdentity()
	if err != nil {
		return nil, err
	}
	channelManager, err := NewChannelManager(dbCfg.LevelDB.Dir, identity, chainCfg, ordererClients)
	if err != nil {
		return nil, err
	}
	server.ChannelManager = channelManager
	server.ordererClients = ordererClients

	return server, nil
}

// 获取ordererClient数组
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
func (s *Server) initServer(engine *gin.Engine) error {
	v1 := engine.Group(Version)
	{
		v1.POST(ActionGetTxStatus, s.GetTxStatusByHTTP)
		v1.POST(ActionListTxHistory, s.ListTxHistoryByHTTP)
	}
	return nil
}

// Start starts the server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Address, s.config.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.New("Failed to start the peer server")
	}

	err = s.ChannelManager.start()
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

	go func() {
		log.Infof("Start the peer server at %s", addr)
		err = s.rpcServer.Serve(lis)
		if err != nil {
			log.Error("gRPC Serve error: ", err)
			return
		}
	}()

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
	s.ChannelManager.stop()
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
