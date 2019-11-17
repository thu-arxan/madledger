package lib

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials"
	"madledger/common/crypto"
	"madledger/core/types"
	"sort"
	"strings"
	"time"

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

// 获取ordererClient数组
func getOrdererClients(cfg *config.Config) ([]pb.OrdererClient, error) {
	var clients []pb.OrdererClient
	for _, address := range cfg.Orderer.Address {
		var opts []grpc.DialOption
		var conn *grpc.ClientConn
		var err error
		if cfg.TLS.Enable {
			creds := credentials.NewTLS(&tls.Config{
				ServerName:   "orderer.madledger.com",
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
				ServerName:   "peer.madledger.com",
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

	// 有多个ordererClient，遍历ordererClient直到成功获取the info of channels
	for i, ordererClient := range c.ordererClients {
		infos, err = ordererClient.ListChannels(context.Background(), &pb.ListChannelsRequest{
			System: system,
			PK:     pk,
		})
		times := i + 1
		if err != nil {
			// 打印出每一个出错信息
			// 如果最后一个ordererClient仍然失败，需要return
			if times == len(c.ordererClients) {
				fmt.Printf("lib/client/ListChannel: try %d times (the last time) but failed to get the info of channels because %s\n", times, err)
				return channelInfos, err
			} else {
				fmt.Printf("lib/client/ListChannel: try %d times but failed to get the info of channels because %s\n", times, err)
			}
		} else {
			fmt.Printf("lib/client/ListChannel: try %d times and success to get %d channels' info\n", times, len(infos.Channels))
			// 获取信息成功，break
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
func (c *Client) CreateChannel(channelID string, public bool, admins, members []*types.Member) error {
	// log.Infof("Create channel %s", channelID)
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
	typesTx, _ := types.NewTx(types.CONFIGCHANNELID, types.CreateChannelContractAddress, payload, c.GetPrivKey(), types.NORMAL)
	pbTx, _ := pb.NewTx(typesTx)

	var times int
	for i, ordererClient := range c.ordererClients {
		//fmt.Printf("Going to create channel %s by ordererClient[%d]\n", channelID, i)
		_, err = ordererClient.CreateChannel(context.Background(), &pb.CreateChannelRequest{
			Tx: pbTx,
		})
		times = i + 1
		if err != nil {
			// try to use other ordererClients until the last one still returns an error
			if times == len(c.ordererClients) {
				fmt.Printf("lib/client/CreateChannel: try %d times (the last time) but failed to create channel %s because %s\n", times, channelID, err)
				return err
			} else {
				fmt.Printf("lib/client/CreateChannel: try %d times but failed to create channel %s because %s\n", times, channelID, err)
			}
		} else {
			// create channel successfully and exit the loop
			fmt.Printf("lib/client/CreateChannel: try %d times and success to create channel %s\n", times, channelID)
			break
		}
	}

	return nil
}

// AddTx try to add a tx
func (c *Client) AddTx(tx *types.Tx) (*pb.TxStatus, error) {
	pbTx, err := pb.NewTx(tx)
	if err != nil {
		return nil, err
	}
	// give client.pubkey to orderer
	pk, err := c.GetPrivKey().PubKey().Bytes()
	for i, ordererClient := range c.ordererClients {
		_, err = ordererClient.AddTx(context.Background(), &pb.AddTxRequest{
			Tx: pbTx,
			PK: pk,
		})

		times := i + 1
		if err != nil {
			// if the client is not system admin, just exit the loop
			if strings.Contains(err.Error(), "the client is not system admin and can not update validator") {
				return nil, err
			}
			// try to use other ordererClients until the last one still returns an error
			if times == len(c.ordererClients) {
				fmt.Printf("lib/client/AddTx: try %d times(the last time) but fail to add tx %s because %s\n", times, tx.ID, err)
				return nil, err
			} else {
				fmt.Printf("lib/client/AddTx: try %d times but fail to add tx %s because %s\n", times, tx.ID, err)
			}
		} else {
			// add tx successfully and exit the loop
			fmt.Printf("lib/client/AddTx: try %d times and success to add tx %s\n", times, tx.ID)
			break
		}
	}

	collector := NewCollector(len(c.peerClients))
	for i := range c.peerClients {
		go func(i int) {
			// fmt.Printf("lib/client/AddTx: going to get tx status from peerClient[%d]\n", i)
			status, err := c.peerClients[i].GetTxStatus(context.Background(), &pb.GetTxStatusRequest{
				ChannelID: tx.Data.ChannelID,
				TxID:      tx.ID,
				Behavior:  pb.Behavior_RETURN_UNTIL_READY,
			})
			if err != nil {
				collector.Add(nil, err)
				fmt.Printf("lib/client/AddTx: peerClient[%d] failed to get tx status because %s\n", i, err)
			} else {
				collector.Add(status, nil)
				// fmt.Printf("lib/client/AddTx: peerClient[%d] success to get tx status\n", i)
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
	} /*else {
		txHistorys := result.(*pb.TxHistory)
		for channel, txs := range txHistorys.Txs {
			for i, id := range txs.Value {
				fmt.Printf("No.%d: Channel %s, tx.ID %s\n", i, channel, id)
			}
		}
	}*/

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

// GetPrivKey return the private key
func (c *Client) GetPrivKey() crypto.PrivateKey {
	return c.privKey
}
