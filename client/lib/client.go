package lib

import (
	"context"
	"fmt"
	"madledger/common/crypto"
	"madledger/core/types"
	"strings"
	"time"

	"google.golang.org/grpc"

	"madledger/client/config"
	pb "madledger/protos"
)

// Client is the Client to communicate with orderer
type Client struct {
	ordererClient pb.OrdererClient
	peerClient    pb.PeerClient
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
	// get clients
	ordererClient, err := getOrdererClient()
	if err != nil {
		return nil, err
	}
	peerClient, err := getPeerClient()
	if err != nil {
		return nil, err
	}

	return &Client{
		ordererClient: ordererClient,
		peerClient:    peerClient,
		privKey:       keyStore.Keys[0],
	}, nil
}

func getOrdererClient() (pb.OrdererClient, error) {
	conn, err := grpc.Dial("localhost:12345", grpc.WithInsecure(), grpc.WithTimeout(2000*time.Millisecond))
	if err != nil {
		return nil, err
	}
	ordererClient := pb.NewOrdererClient(conn)
	return ordererClient, nil
}

func getPeerClient() (pb.PeerClient, error) {
	conn, err := grpc.Dial("localhost:23456", grpc.WithInsecure(), grpc.WithTimeout(2000*time.Millisecond))
	if err != nil {
		return nil, err
	}
	peerClient := pb.NewPeerClient(conn)
	return peerClient, nil
}

// GetPrivKey return the private key
func (c *Client) GetPrivKey() crypto.PrivateKey {
	return c.privKey
}

// ListChannel list the info of channel
func (c *Client) ListChannel(system bool) ([]*ChannelInfo, error) {
	infos, err := c.ordererClient.ListChannels(context.Background(), &pb.ListChannelsRequest{
		System: system,
	})
	if err != nil {
		return nil, err
	}
	var channelInfos []*ChannelInfo
	for i, channel := range infos.Channels {
		channelInfos = append(channelInfos, &ChannelInfo{
			Name:      channel.ChannelID,
			System:    false,
			BlockSize: channel.BlockSize,
		})
		if strings.HasPrefix(channel.ChannelID, "_") {
			channelInfos[i].System = true
		}
	}

	return channelInfos, nil
}

// CreateChannel create a channel
func (c *Client) CreateChannel(channelID string) error {
	_, err := c.ordererClient.AddChannel(context.Background(), &pb.AddChannelRequest{
		ChannelID: channelID,
	})
	if err != nil {
		fmt.Printf("Failed to create channel %s because %s\n", channelID, err)
		return err
	}

	fmt.Println("Succeed!")
	return nil
}

// AddTx try to add a tx
func (c *Client) AddTx(tx *types.Tx) (*pb.TxStatus, error) {
	pbTx, err := pb.NewTx(tx)
	if err != nil {
		return nil, err
	}
	_, err = c.ordererClient.AddTx(context.Background(), &pb.AddTxRequest{
		Tx: pbTx,
	})
	if err != nil {
		return nil, err
	}
	time.Sleep(1 * time.Second)
	status, err := c.peerClient.GetTxStatus(context.Background(), &pb.GetTxStatusRequest{
		ChannelID: tx.Data.ChannelID,
		TxID:      tx.ID,
		Behavior:  pb.Behavior_RETURN_UNTIL_READY,
	})
	if err != nil {
		return nil, err
	}
	return status, nil
}
