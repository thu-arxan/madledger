package orderer

import (
	"context"
	"time"

	"madledger/core/types"
	pb "madledger/protos"

	"google.golang.org/grpc"
)

// Client is the client of orderer
type Client struct {
	ordererClient pb.OrdererClient
}

// NewClient is the constructor of Client
func NewClient(addr string) (*Client, error) {
	var conn *grpc.ClientConn
	var err error
	conn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(2000*time.Millisecond))
	if err != nil {
		return nil, err
	}
	ordererClient := pb.NewOrdererClient(conn)

	return &Client{
		ordererClient: ordererClient,
	}, nil
}

// FetchBlock return block if exist, else return error
func (c *Client) FetchBlock(channelID string, num uint64, async bool) (*types.Block, error) {
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
	typesBlock, err := ConvertBlockFromPbToTypes(pbBlock)
	return typesBlock, err
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
