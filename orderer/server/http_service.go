package server

import (
	"encoding/hex"
	"encoding/json"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/core"
	pb "madledger/protos"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListChannelReq Binding from JSON
type ListChannelReq struct {
	System string `form:"system" json:"system" xml:"system"  binding:"required"`
	PK     string `form:"pk" json:"pk" xml:"pk" binding:"required"`
	Algo   string `form:"algo" json:"algo" xml:"algo" binding:"required"`
}

// ListChannelsByHTTP list channels by http
func (hs *Server) ListChannelsByHTTP(c *gin.Context) {
	var json ListChannelReq
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	system, _ := strconv.ParseBool(json.System)
	pk, _ := hex.DecodeString(json.PK)
	algo, _ := strconv.ParseInt(json.Algo, 16, 64)
	req := &pb.ListChannelsRequest{
		System: system,
		PK:     pk,
		Algo:   int32(algo),
	}
	log.Infof("orderer receive pk is %v", json.PK)
	info, err := hs.cc.ListChannels(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"channelinfo": info})
	return
}

// CreateChannelReq ...
type CreateChannelReq struct {
	Tx string `json:"tx"`
}

// CreateChannelByHTTP create channel by http
func (hs *Server) CreateChannelByHTTP(c *gin.Context) {
	var j CreateChannelReq
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var coreTx core.Tx
	json.Unmarshal([]byte(j.Tx), &coreTx)
	if !coreTx.Verify() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The tx is not a valid tx"})
		return
	}
	if coreTx.GetReceiver().String() != core.CreateChannelContractAddress.String() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The receiver of the tx is not the valid contract address"})
		return
	}
	_, err := hs.cc.CreateChannel(&coreTx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
	return
}

// AddTxReq ...
type AddTxReq struct {
	Tx string `json:"tx"`
}

// AddTxByHTTP add tx by http
func (hs *Server) AddTxByHTTP(c *gin.Context) {
	var j CreateChannelReq
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var status pb.TxStatus

	var coreTx core.Tx
	json.Unmarshal([]byte(j.Tx), &coreTx)

	txType, err := core.GetTxType(common.BytesToAddress(coreTx.Data.Recipient).String())
	if err == nil && (txType == core.VALIDATOR || txType == core.NODE) {
		pk, err := crypto.NewPublicKey(coreTx.Data.Sig.PK, coreTx.Data.Sig.Algo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// create member to check if the client is system admin
		member, err := core.NewMember(pk, "")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !hs.cc.CM.IsSystemAdmin(member) { // not system admin, return error
			c.JSON(http.StatusBadRequest, gin.H{"error": "The client is not system admin and can't config the cluster"})
			return
		}
	}
	err = hs.cc.AddTx(&coreTx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"txstatus": status})
	return
}

// AccountInfoReq ...
type AccountInfoReq struct {
	Addr string `json:"address"`
}

// GetAccountInfoByHTTP get account info by http
func (hs *Server) GetAccountInfoByHTTP(c *gin.Context) {
	var j AccountInfoReq
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var accountInfo pb.AccountInfo
	str, err := hex.DecodeString(j.Addr)
	addr := common.BytesToAddress(str)
	account, err := hs.cc.AM.GetAccount(addr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accountInfo.Balance = account.GetBalance()
	c.JSON(http.StatusOK, gin.H{"accountinfo": accountInfo})
	return
}
