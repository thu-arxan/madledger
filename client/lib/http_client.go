package lib

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	cc "madledger/blockchain/config"
	"madledger/client/config"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/core"
	pb "madledger/protos"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// HTTPClient is the Client to communicate with orderer and peers
type HTTPClient struct {
	ordererHTTPClients []string
	peerHTTPClients    []string
	privKey            crypto.PrivateKey
}

// NewHTTPClient is the constructor of HTTPClient
func NewHTTPClient(cfgFile string) (*HTTPClient, error) {
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return nil, err
	}

	return NewHTTPClientFromConfig(cfg)
}

// NewHTTPClientFromConfig will construct http client from cfg
func NewHTTPClientFromConfig(cfg *config.Config) (*HTTPClient, error) {
	// get clients
	ordererClients, err := getOrdererHTTPClients(cfg)
	if err != nil {
		return nil, err
	}
	peerClients, err := getPeerHTTPClients(cfg)
	if err != nil {
		return nil, err
	}

	return &HTTPClient{
		ordererHTTPClients: ordererClients,
		peerHTTPClients:    peerClients,
		privKey:            cfg.KeyStore.Privs[0],
	}, nil
}

func getOrdererHTTPClients(cfg *config.Config) ([]string, error) {
	var clients []string
	for _, address := range cfg.Orderer.HTTPAddress {
		clients = append(clients, address)
	}
	return clients, nil
}

func getPeerHTTPClients(cfg *config.Config) ([]string, error) {
	var clients []string
	for _, address := range cfg.Peer.HTTPAddress {
		clients = append(clients, address)
	}
	return clients, nil
}

//ListChannelResp ...
type ListChannelResp struct {
	ChannelInfos *pb.ChannelInfos `json:"channelinfo"`
}

// ListChannelByHTTP list the info of channel
func (c *HTTPClient) ListChannelByHTTP(system bool) ([]ChannelInfo, error) {
	var channelInfos []ChannelInfo
	pubKey := c.GetPrivKey().PubKey()
	pk, err := pubKey.Bytes()
	if err != nil {
		return channelInfos, err
	}
	var infos ListChannelResp
	// var result map[string]interface{}
	for i, ordererHTTPClient := range c.ordererHTTPClients {
		requestBody, _ := json.Marshal(map[string]string{
			"system": strconv.FormatBool(system),
			"pk":     hex.EncodeToString(pk),
			"algo":   strconv.FormatInt(int64(pubKey.Algo()), 16),
		})
		log.Infof("client send pk is %v", hex.EncodeToString(pk))

		resp, err := http.Post("http://"+ordererHTTPClient+"/v1/listchannels", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			return channelInfos, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &infos)

		times := i + 1
		if err != nil {
			if times == len(c.ordererHTTPClients) {
				return channelInfos, err
			}
		} else {
			break
		}
	}
	for i, channel := range infos.ChannelInfos.Channels {
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

//CreateChannelResp ...
type CreateChannelResp struct {
	Error string `json:"error"`
}

// CreateChannelByHTTP create a channel
func (c *HTTPClient) CreateChannelByHTTP(channelID string, public bool, admins, members []*core.Member,
	gasPrice uint64, ratio uint64, maxGas uint64) error {
	// log.Infof("Create channel %s", channelID)
	self, err := core.NewMember(c.privKey.PubKey(), "admin")
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
			Public:          public,
			Admins:          admins,
			Members:         members,
			GasPrice:        gasPrice,
			AssetTokenRatio: ratio,
			MaxGas:          maxGas,
		},
		Version: 1,
	})
	coreTx, _ := core.NewTx(core.CONFIGCHANNELID, core.CreateChannelContractAddress, payload, 0, "", c.privKey)

	var info CreateChannelResp
	var times int
	for i, ordererHTTPClient := range c.ordererHTTPClients {
		coreTxBytes, _ := coreTx.Bytes()
		requestBody, _ := json.Marshal(map[string]string{
			"tx": string(coreTxBytes),
		})
		resp, err := http.Post("http://"+ordererHTTPClient+"/v1/createchannel", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &info)
		if info.Error != "" {
			return errors.New(info.Error)
		}
		times = i + 1
		if err != nil {
			// try to use other ordererClients until the last one still returns an error
			if times == len(c.ordererHTTPClients) {
				return err
			}
		} else {
			break
		}
	}

	return nil
}

// AddTxResp ...
type AddTxResp struct {
	Error string `json:"error"`
}

// GetTxStatusResp ...
type GetTxStatusResp struct {
	Error  string       `json:"error"`
	Status *pb.TxStatus `json:"status"`
}

// AddTxByHTTP try to add a tx
// TODO: Support bft
func (c *HTTPClient) AddTxByHTTP(tx *core.Tx) (*pb.TxStatus, error) {
	var info AddTxResp
	for i, ordererHTTPClient := range c.ordererHTTPClients {
		coreTxBytes, _ := tx.Bytes()
		requestBody, _ := json.Marshal(map[string]string{
			"tx": string(coreTxBytes),
		})
		resp, err := http.Post("http://"+ordererHTTPClient+"/v1/addtx", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		log.Infof("create resp is: %s", string(body))
		err = json.Unmarshal(body, &info)
		if info.Error != "" {
			return nil, errors.New(info.Error)
		}
		times := i + 1
		if err != nil {
			// if the client is not system admin, just exit the loop
			if strings.Contains(err.Error(), "the client is not system admin and can not update validator") {
				return nil, err
			}
			// try to use other ordererClients until the last one still returns an error
			if times == len(c.ordererHTTPClients) {
				return nil, err
			}
		} else {
			// add tx successfully and exit the loop
			break
		}
	}

	collector := NewCollector(len(c.peerHTTPClients), 1)
	for i := range c.peerHTTPClients {
		go func(i int) {
			var info GetTxStatusResp
			requestBody, _ := json.Marshal(map[string]string{
				"channelID": tx.Data.ChannelID,
				"txID":      tx.ID,
			})
			log.Infof("ask %s for txstatus", c.peerHTTPClients[i])
			resp, err := http.Post("http://"+c.peerHTTPClients[i]+"/v1/gettxstatus", "application/json", bytes.NewBuffer(requestBody))
			log.Infof("get response from %s", c.peerHTTPClients[i])
			if err != nil {
				collector.AddError(err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			log.Info(string(body))
			json.Unmarshal(body, &info)

			if err != nil || info.Error != "" {
				collector.AddError(err)
			} else {
				collector.Add(info.Status)
			}
		}(i)
	}

	result, err := collector.Wait()
	if err != nil {
		return nil, err
	}
	return result.(*pb.TxStatus), nil
}

// GetTxHistoryResp ...
type GetTxHistoryResp struct {
	Error     string        `json:"error"`
	TxHistory *pb.TxHistory `json:"txhistory"`
}

// GetHistoryByHTTP return the history of address
// TODO: Support bft
func (c *HTTPClient) GetHistoryByHTTP(address []byte) (*pb.TxHistory, error) {
	collector := NewCollector(len(c.peerHTTPClients), 1)
	for i := range c.peerHTTPClients {
		go func(i int) {
			var info GetTxHistoryResp
			requestBody, _ := json.Marshal(map[string]string{
				"address": hex.EncodeToString(address),
			})
			resp, err := http.Post("http://"+c.peerHTTPClients[i]+"/v1/listtxhistory", "application/json", bytes.NewBuffer(requestBody))
			if err != nil {
				collector.AddError(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			json.Unmarshal(body, &info)

			if err != nil || info.Error != "" {
				collector.AddError(errors.New(info.Error))
			} else {
				collector.Add(info.TxHistory)
			}
		}(i)
	}

	result, err := collector.Wait()
	if err != nil {
		return nil, err
	}

	return result.(*pb.TxHistory), err
}

// GetPrivKey return the private key
func (c *HTTPClient) GetPrivKey() crypto.PrivateKey {
	return c.privKey
}

//GetAccountBalanceResp ...
type GetAccountBalanceResp struct {
	Error   string         `json:"error"`
	Account pb.AccountInfo `json:"accountinfo"`
}

//GetAccountBalanceByHTTP Get Account Balance By HTTP
func (c *HTTPClient) GetAccountBalanceByHTTP(address common.Address) (uint64, error) {
	var times int
	var info GetAccountBalanceResp
	for i := range c.ordererHTTPClients {
		requestBody, _ := json.Marshal(map[string]string{
			"address": hex.EncodeToString(address.Bytes()),
		})
		resp, err := http.Post("http://"+c.ordererHTTPClients[i]+"/v1/getaccountinfo", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			return 0, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		err = json.Unmarshal(body, &info)
		if err != nil {
			return 0, err
		}
		times = i + 1
		if err != nil {
			// try to use other ordererClients until the last one still returns an error
			if times == len(c.ordererHTTPClients) {
				return 0, err
			}
		} else {
			break
		}

	}
	return info.Account.GetBalance(), nil
}

// GetTokenInfoResp ...
type GetTokenInfoResp struct {
	Error string        `json:"error"`
	Token *pb.TokenInfo `json:"tokeninfo"`
}

// GetTokenInfoByHTTP Get Token Info By HTTP
func (c *HTTPClient) GetTokenInfoByHTTP(address common.Address, channelID []byte) (uint64, error) {
	collector := NewCollector(len(c.peerHTTPClients), 1)
	var info GetTokenInfoResp
	for i := range c.peerHTTPClients {
		go func(i int) {
			requestBody, _ := json.Marshal(map[string]string{
				"address":   hex.EncodeToString(address.Bytes()),
				"channelid": string(channelID),
			})
			resp, err := http.Post("http://"+c.peerHTTPClients[i]+"/v1/gettokeninfo", "application/json", bytes.NewBuffer(requestBody))
			if err != nil {
				collector.AddError(err)
			}
			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			json.Unmarshal(body, &info)

			if err != nil || info.Error != "" {
				collector.AddError(errors.New(info.Error))
			} else {
				collector.Add(info.Token)
			}
		}(i)
	}
	result, err := collector.Wait()
	if err != nil {
		return 0, err
	}
	return result.(*pb.TokenInfo).GetBalance(), err

}
