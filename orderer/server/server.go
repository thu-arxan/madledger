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
	"fmt"
	"madledger/orderer/channel"
	"madledger/orderer/config"
	pb "madledger/protos"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc/credentials"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "orderer", "package": "server"})
)

// Server provide the serve of orderer
type Server struct {
	sync.RWMutex
	config    *config.ServerConfig
	rpcServer *grpc.Server
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


type ZapLogger struct {
}

// NewZapLogger 创建封装了zap的对象，该对象是对LoggerV2接口的实现
func NewZapLogger() *ZapLogger {
	return &ZapLogger{
	}
}

// Info returns
func (zl *ZapLogger) Info(args ...interface{}) {
	fmt.Println(args...)
}

// Infoln returns
func (zl *ZapLogger) Infoln(args ...interface{}) {
	fmt.Println(args...)
}

// Infof returns
func (zl *ZapLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Warning returns
func (zl *ZapLogger) Warning(args ...interface{}) {
	fmt.Println(args...)
}

// Warningln returns
func (zl *ZapLogger) Warningln(args ...interface{}) {
	fmt.Println(args...)
}

// Warningf returns
func (zl *ZapLogger) Warningf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Error returns
func (zl *ZapLogger) Error(args ...interface{}) {
	fmt.Println(args...)
}

// Errorln returns
func (zl *ZapLogger) Errorln(args ...interface{}) {
	fmt.Println(args...)
}

// Errorf returns
func (zl *ZapLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Fatal returns
func (zl *ZapLogger) Fatal(args ...interface{}) {
	fmt.Println(args...)
}

// Fatalln returns
func (zl *ZapLogger) Fatalln(args ...interface{}) {
	fmt.Println(args...)
}

// Fatalf logs to fatal level
func (zl *ZapLogger) Fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// V reports whether verbosity level l is at least the requested verbose level.
func (zl *ZapLogger) V(v int) bool {
	return false
}


// Start starts the server
func (s *Server) Start() error {
	var logger = NewZapLogger()
	grpclog.SetLoggerV2(logger);
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
	if err != nil {
		return err
	}

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
}
