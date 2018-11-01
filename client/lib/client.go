package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"madledger/common/crypto"
	"madledger/core/types"
	"strings"
	"time"

	"google.golang.org/grpc"

	cc "madledger/blockchain/config"
	"madledger/client/config"
	pb "madledger/protos"
)

// Client is the Client to communicate with orderer
type Client struct {
	ordererClient pb.OrdererClient
	peerClients   []pb.PeerClient
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
	ordererClient, err := getOrdererClient(cfg)
	if err != nil {
		return nil, err
	}
	peerClients, err := getPeerClients(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		ordererClient: ordererClient,
		peerClients:   peerClients,
		privKey:       keyStore.Keys[0],
	}, nil
}

func getOrdererClient(cfg *config.Config) (pb.OrdererClient, error) {
	conn, err := grpc.Dial(cfg.Orderer.Address[0], grpc.WithInsecure(), grpc.WithTimeout(2000*time.Millisecond))
	if err != nil {
		return nil, err
	}
	ordererClient := pb.NewOrdererClient(conn)
	return ordererClient, nil
}

func getPeerClients(cfg *config.Config) ([]pb.PeerClient, error) {
	var clients []pb.PeerClient
	for _, address := range cfg.Peer.Address {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(2000*time.Millisecond))
		if err != nil {
			return nil, err
		}
		peerClient := pb.NewPeerClient(conn)
		clients = append(clients, peerClient)
	}

	return clients, nil
}

// GetPrivKey return the private key
func (c *Client) GetPrivKey() crypto.PrivateKey {
	return c.privKey
}

// ListChannel list the info of channel
func (c *Client) ListChannel(system bool) ([]ChannelInfo, error) {
	var channelInfos []ChannelInfo
	pk, err := c.GetPrivKey().PubKey().Bytes()
	if err != nil {
		return channelInfos, err
	}
	infos, err := c.ordererClient.ListChannels(context.Background(), &pb.ListChannelsRequest{
		System: system,
		PK:     pk,
	})
	if err != nil {
		return channelInfos, err
	}
	for i, channel := range infos.Channels {
		channelInfos = append(channelInfos, ChannelInfo{
			Name:      channel.ChannelID,
			System:    false,
			BlockSize: channel.BlockSize,
			Identity:  channel.Identity.String(),
		})
		if strings.HasPrefix(channel.ChannelID, "_") {
			channelInfos[i].System = true
		}
	}
	return channelInfos, nil
}

// CreateChannel create a channel
func (c *Client) CreateChannel(channelID string, public bool, admins, members []*types.Member) error {
	self, err := types.NewMember(c.GetPrivKey().PubKey(), "admin")
	if err != nil {
		return err
	}
	admins = unionMembers(admins, []*types.Member{self})
	// if this is a public channel, there is no need to contain members
	if public {
		members = make([]*types.Member, 0)
	} else {
		members = unionMembers(admins, members)
	}
	payload, _ := json.Marshal(cc.Payload{
		ChannelID: channelID,
		Profile: &cc.Profile{
			Public:  public,
			Admins:  admins,
			Members: members,
		},
		Version: 1,
	})
	typesTx, _ := types.NewTx(types.CONFIGCHANNELID, types.CreateChannelContractAddress, payload, c.GetPrivKey())
	pbTx, _ := pb.NewTx(typesTx)

	_, err = c.ordererClient.CreateChannel(context.Background(), &pb.CreateChannelRequest{
		Tx: pbTx,
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
	collector := NewCollector(len(c.peerClients))
	for i := range c.peerClients {
		go func(i int) {
			status, err := c.peerClients[i].GetTxStatus(context.Background(), &pb.GetTxStatusRequest{
				ChannelID: tx.Data.ChannelID,
				TxID:      tx.ID,
				Behavior:  pb.Behavior_RETURN_UNTIL_READY,
			})
			if err != nil {
				collector.Add(nil, err)
			} else {
				collector.Add(status, nil)
			}
		}(i)
	}

	result, err := collector.Wait()
	if err != nil {
		return nil, err
	}
	return result.(*pb.TxStatus), nil
}

// GetHistory return the history of address
func (c *Client) GetHistory(address []byte) (*pb.TxHistory, error) {
	collector := NewCollector(len(c.peerClients))
	for i := range c.peerClients {
		go func(i int) {
			history, err := c.peerClients[i].ListTxHistory(context.Background(), &pb.ListTxHistoryRequest{
				Address: address,
			})
			if err != nil {
				collector.Add(nil, err)
			} else {
				collector.Add(history, nil)
			}
		}(i)
	}

	result, err := collector.Wait()
	if err != nil {
		return nil, err
	}

	return result.(*pb.TxHistory), err
}

func unionMembers(first, second []*types.Member) []*types.Member {
	union := make([]*types.Member, 0)
	members := append(first, second...)
	for _, member := range members {
		if !membersContain(union, member) {
			union = append(union, member)
		}
	}
	return union
}

func membersContain(members []*types.Member, member *types.Member) bool {
	for i := range members {
		if member.Equal(members[i]) {
			return true
		}
	}
	return false
}
