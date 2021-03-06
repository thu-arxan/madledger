// Copyright (c) 2020 THU-Arxan
// Madledger is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package orderer

import (
	"context"
	"crypto/tls"
	"madledger/peer/config"
	"time"

	"google.golang.org/grpc/credentials"

	"madledger/core"
	pb "madledger/protos"

	"google.golang.org/grpc"
)

// Client is the client of orderer
type Client struct {
	ordererClient pb.OrdererClient
}

// NewClient is the constructor of Client
func NewClient(addr string, cfg *config.Config) (*Client, error) {
	var opts []grpc.DialOption
	var conn *grpc.ClientConn
	var err error
	if cfg.TLS.Enable {
		creds := credentials.NewTLS(&tls.Config{
			//ServerName:   "orderer.madledger.com",
			Certificates: []tls.Certificate{*(cfg.TLS.Cert)},
			RootCAs:      cfg.TLS.Pool,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithTimeout(2000*time.Millisecond))

	conn, err = grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	ordererClient := pb.NewOrdererClient(conn)
	return &Client{
		ordererClient: ordererClient,
	}, nil
}

// FetchBlock return block if exist, else return error
func (c *Client) FetchBlock(channelID string, num uint64, async bool) (*core.Block, error) {
	var behavior = pb.Behavior_FAIL_IF_NOT_READY
	if async {
		behavior = pb.Behavior_RETURN_UNTIL_READY
	}
	pbBlock, err := c.ordererClient.FetchBlock(context.Background(), &pb.FetchBlockRequest{
		ChannelID: channelID,
		Number:    num,
		Behavior:  behavior,
	})
	if err != nil {
		return nil, err
	}
	return pbBlock.ToCore()
}

// ListChannels return all channels
func (c *Client) ListChannels() ([]string, error) {
	channelInfos, err := c.ordererClient.ListChannels(context.Background(), &pb.ListChannelsRequest{
		System: false,
	})
	var channels []string
	if err != nil {
		return nil, err
	}
	for _, channelInfo := range channelInfos.Channels {
		channels = append(channels, channelInfo.ChannelID)
	}
	return channels, nil
}
