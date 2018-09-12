package orderer

import (
	"context"
	"fmt"
	"madledger/common/crypto"
	"madledger/core/types"
	"strings"
	"time"

	"github.com/modood/table"
	"google.golang.org/grpc"

	"madledger/client/config"
	pb "madledger/protos"
)

// Client is the Client to communicate with orderer
type Client struct {
	ordererClient pb.OrdererClient
	privKey       crypto.PrivateKey
}

// NewClient is the constructor of pb.OrdereClient
func NewClient(cfgFile string) (*Client, error) {
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return nil, err
	}
	keyStore, err := cfg.GetKeyStoreConfig()
	if err != nil {
		return nil, err
	}
	var conn *grpc.ClientConn

	conn, err = grpc.Dial("localhost:12345", grpc.WithInsecure(), grpc.WithTimeout(2000*time.Millisecond))
	if err != nil {
		return nil, err
	}
	ordererClient := pb.NewOrdererClient(conn)

	return &Client{
		ordererClient: ordererClient,
		privKey:       keyStore.Keys[0],
	}, nil
}

// GetPrivKey return the private key
func (c *Client) GetPrivKey() crypto.PrivateKey {
	return c.privKey
}

type channelInfo struct {
	Name      string
	System    bool
	BlockSize uint64
}

// ListChannel list the info of channel
func (c *Client) ListChannel(system bool) error {
	infos, err := c.ordererClient.ListChannels(context.Background(), &pb.ListChannelsRequest{
		System: system,
	})
	if err != nil {
		return err
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
	if len(channelInfos) == 0 {
		fmt.Println("No results!")
	} else {
		table.Output(channelInfos)
	}

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

// AddTx try to add a tx
func (c *Client) AddTx(tx *types.Tx) error {
	pbTx, err := pb.NewTx(tx)
	if err != nil {
		return err
	}
	_, err = c.ordererClient.AddTx(context.Background(), &pb.AddTxRequest{
		Tx: pbTx,
	})
	return err
}