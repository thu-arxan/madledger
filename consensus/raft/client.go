package raft

import (
	"context"
	"crypto/tls"
	"madledger/consensus"
	pb "madledger/consensus/raft/protos"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client is the clients keep connections to blockchain
type Client struct {
	sync.RWMutex
	addr string
	conn *grpc.ClientConn
	TLS  consensus.TLSConfig
}

// NewClient is the constructor of Client
func NewClient(addr string, tlsConfig consensus.TLSConfig) (*Client, error) {
	return &Client{
		addr: addr,
		TLS:  tlsConfig,
	}, nil
}

// newConn check whether conn is nil and init it if conn is nil
func (c *Client) newConn() error {
	c.RLock()

	if c.conn != nil {
		c.RUnlock()
		return nil
	}
	c.RUnlock()

	c.Lock()
	defer c.Unlock()

	var opts []grpc.DialOption
	var conn *grpc.ClientConn
	var err error
	if c.TLS.Enable {
		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{*(c.TLS.Cert)},
			RootCAs:      c.TLS.Pool,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithTimeout(2000*time.Millisecond))
	conn, err = grpc.Dial(c.addr, opts...)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) addTx(channelID string, tx []byte, caller uint64) error {
	// call newConn every time, avoid repeated lock
	if err := c.newConn(); err != nil {
		return err
	}

	client := pb.NewBlockChainClient(c.conn)
	_, err := client.AddTx(context.Background(), &pb.RaftTX{
		Tx:      tx,
		Caller:  caller,
		Channel: channelID,
	})
	log.Infof("after call %s, %v", c.addr, err)
	return err
}

// close closes grpc connections and set conn as nil
func (c *Client) close() {
	c.Lock()
	c.Unlock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}
