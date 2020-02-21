package lib

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"madledger/common/crypto"
	"madledger/core"
	"sort"
	"strings"
	"time"
	"madledger/common"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"

	cc "madledger/blockchain/config"
	"madledger/client/config"
	pb "madledger/protos"
)

var (
	log = logrus.WithFields(logrus.Fields{"app": "client", "package": "lib"})
)

// Client is the Client to communicate with orderer
type Client struct {
	ordererClients []pb.OrdererClient
	peerClients    []pb.PeerClient
	privKey        crypto.PrivateKey
}

// NewClient is the constructor of pb.OrdereClient
func NewClient(cfgFile string) (*Client, error) {
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return nil, err
	}

	return NewClientFromConfig(cfg)
}

// NewClientFromConfig will construct client from cfg
func NewClientFromConfig(cfg *config.Config) (*Client, error) {
	keyStore, err := cfg.GetKeyStoreConfig()
	if err != nil {
		return nil, err
	}
	// get clients
	ordererClients, err := getOrdererClients(cfg)
	if err != nil {
		return nil, err
	}
	peerClients, err := getPeerClients(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		ordererClients: ordererClients,
		peerClients:    peerClients,
		privKey:        keyStore.Keys[0],
	}, nil
}

func getOrdererClients(cfg *config.Config) ([]pb.OrdererClient, error) {
	var clients []pb.OrdererClient
	for _, address := range cfg.Orderer.Address {
		var opts []grpc.DialOption
		var conn *grpc.ClientConn
		var err error
		if cfg.TLS.Enable {
			creds := credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{*(cfg.TLS.Cert)},
				RootCAs:      cfg.TLS.Pool,
			})
			opts = append(opts, grpc.WithTransportCredentials(creds))
		} else {
			opts = append(opts, grpc.WithInsecure())
		}
		opts = append(opts, grpc.WithTimeout(2000*time.Millisecond))
		conn, err = grpc.Dial(address, opts...)
		if err != nil {
			return nil, err
		}
		ordererClient := pb.NewOrdererClient(conn)
		clients = append(clients, ordererClient)
	}

	return clients, nil
}

func getPeerClients(cfg *config.Config) ([]pb.PeerClient, error) {
	var clients []pb.PeerClient
	for _, address := range cfg.Peer.Address {
		var opts []grpc.DialOption
		var conn *grpc.ClientConn
		var err error
		if cfg.TLS.Enable {
			creds := credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{*(cfg.TLS.Cert)},
				RootCAs:      cfg.TLS.Pool,
			})
			opts = append(opts, grpc.WithTransportCredentials(creds))
		} else {
			opts = append(opts, grpc.WithInsecure())
		}
		opts = append(opts, grpc.WithTimeout(2000*time.Millisecond))

		conn, err = grpc.Dial(address, opts...)
		if err != nil {
			return nil, err
		}

		peerClient := pb.NewPeerClient(conn)
		clients = append(clients, peerClient)
	}

	return clients, nil
}

// ListChannel list the info of channel
func (c *Client) ListChannel(system bool) ([]ChannelInfo, error) {
	var channelInfos []ChannelInfo
	pk, err := c.GetPrivKey().PubKey().Bytes()
	if err != nil {
		return channelInfos, err
	}
	var infos *pb.ChannelInfos

	for i, ordererClient := range c.ordererClients {
		infos, err = ordererClient.ListChannels(context.Background(), &pb.ListChannelsRequest{
			System: system,
			PK:     pk,
		})
		times := i + 1
		if err != nil {
			if times == len(c.ordererClients) {
				return channelInfos, err
			}
		} else {
			break
		}
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

	// sort slice
	sort.Slice(channelInfos, func(i, j int) bool {
		if channelInfos[i].System != channelInfos[j].System {
			if channelInfos[i].System {
				return true
			}
			return false
		}
		return channelInfos[i].Name < channelInfos[j].Name
	})

	return channelInfos, nil
}

// CreateChannel create a channel
func (c *Client) CreateChannel(channelID string, public bool, admins, members []*core.Member) error {
	// log.Infof("Create channel %s", channelID)
	self, err := core.NewMember(c.GetPrivKey().PubKey(), "admin")
	if err != nil {
		return err
	}
	admins = unionMembers(admins, []*core.Member{self})
	// if this is a public channel, there is no need to contain members
	if public {
		members = make([]*core.Member, 0)
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
	coreTx, _ := core.NewTx(core.CONFIGCHANNELID, core.CreateChannelContractAddress, payload, 0, "", c.GetPrivKey())
	pbTx, _ := pb.NewTx(coreTx)

	var times int
	for i, ordererClient := range c.ordererClients {
		_, err = ordererClient.CreateChannel(context.Background(), &pb.CreateChannelRequest{
			Tx: pbTx,
		})
		times = i + 1
		if err != nil {
			// try to use other ordererClients until the last one still returns an error
			if times == len(c.ordererClients) {
				return err
			}
		} else {
			break
		}
	}

	return nil
}

// AddTx try to add a tx
// TODO: Support bft
func (c *Client) AddTx(tx *core.Tx) (*pb.TxStatus, error) {
	pbTx, err := pb.NewTx(tx)
	if err != nil {
		return nil, err
	}

	for i, ordererClient := range c.ordererClients {
		_, err = ordererClient.AddTx(context.Background(), &pb.AddTxRequest{
			Tx: pbTx,
		})

		times := i + 1
		if err != nil {
			// if the client is not system admin, just exit the loop
			if strings.Contains(err.Error(), "the client is not system admin and can not update validator") {
				return nil, err
			}
			// try to use other ordererClients until the last one still returns an error
			if times == len(c.ordererClients) {
				return nil, err
			}
		} else {
			// add tx successfully and exit the loop
			break
		}
	}

	collector := NewCollector(len(c.peerClients), 1)
	for i := range c.peerClients {
		go func(i int) {
			status, err := c.peerClients[i].GetTxStatus(context.Background(), &pb.GetTxStatusRequest{
				ChannelID: tx.Data.ChannelID,
				TxID:      tx.ID,
				Behavior:  pb.Behavior_RETURN_UNTIL_READY,
			})
			if err != nil {
				collector.AddError(err)
			} else {
				collector.Add(status)
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
// TODO: Support bft
func (c *Client) GetHistory(address []byte) (*pb.TxHistory, error) {
	collector := NewCollector(len(c.peerClients), 1)
	for i := range c.peerClients {
		go func(i int) {
			history, err := c.peerClients[i].ListTxHistory(context.Background(), &pb.ListTxHistoryRequest{
				Address: address,
			})
			if err != nil {
				collector.AddError(err)
			} else {
				collector.Add(history)
			}
		}(i)
	}

	result, err := collector.Wait()
	if err != nil {
		return nil, err
	}

	return result.(*pb.TxHistory), err
}

func unionMembers(first, second []*core.Member) []*core.Member {
	union := make([]*core.Member, 0)
	members := append(first, second...)
	for _, member := range members {
		if !membersContain(union, member) {
			union = append(union, member)
		}
	}
	return union
}

func membersContain(members []*core.Member, member *core.Member) bool {
	for i := range members {
		if member.Equal(members[i]) {
			return true
		}
	}
	return false
}

// GetPrivKey return the private key
func (c *Client) GetPrivKey() crypto.PrivateKey {
	return c.privKey
}

func (c *Client) GetAccountBalance(address common.Address) uint64 {
	var times int
	var acc *pb.AccountInfo
	var err error
	for i, ordererClient := range c.ordererClients {
		acc, err = ordererClient.GetAccountInfo(context.Background(), &pb.GetAccountInfoRequest{
			Address: address.Bytes(),
		})
		times = i + 1
		if err != nil {
			// try to use other ordererClients until the last one still returns an error
			if times == len(c.ordererClients) {
				panic(err.Error())
				return 0
			}
		} else {
			break
		}
	}
	return acc.GetBalance()
}