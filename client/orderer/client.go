package orderer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/modood/table"
	"google.golang.org/grpc"

	pb "madledger/protos"
)

// Client is the Client to communicate with orderer
type Client struct {
	ordererClient pb.OrdererClient
}

// NewClient is the constructor of pb.OrdereClient
func NewClient() (*Client, error) {
	var conn *grpc.ClientConn
	var err error
	conn, err = grpc.Dial("localhost:12345", grpc.WithInsecure(), grpc.WithTimeout(2000*time.Millisecond))
	if err != nil {
		return nil, err
	}
	ordererClient := pb.NewOrdererClient(conn)

	return &Client{
		ordererClient: ordererClient,
	}, nil
}

type channelInfo struct {
	Name      string
	System    bool
	BlockSize uint64
}

// ListChannel list the info of channel
func (c *Client) ListChannel() error {
	infos, err := c.ordererClient.ListChannels(context.Background(), &pb.ListChannelsRequest{
		System: true,
	})
	if err != nil {
		fmt.Println(err)
	}
	var channelInfos []channelInfo
	for i, channel := range infos.Channels {
		channelInfos = append(channelInfos, channelInfo{
			Name:      channel.ChannelID,
			System:    false,
			BlockSize: channel.BlockSize,
		})
		if strings.HasPrefix(channel.ChannelID, "_") {
			channelInfos[i].System = true
		}
	}
	table.Output(channelInfos)
	return nil
}

// CreateChannel create a channel
func (c *Client) CreateChannel(channelID string) error {
	_, err := c.ordererClient.AddChannel(context.Background(), &pb.AddChannelRequest{
		ChannelID: channelID,
	})
	if err != nil {
		fmt.Printf("Failed to create channel %s because %s\n", channelID, err)
	} else {
		fmt.Println("Succeed!")
	}
	return nil
}
