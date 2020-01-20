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

// FetchBlockByHTTP gets block by http
func (hs *Server) FetchBlockByHTTP(c *gin.Context) {
	channelID := c.Query("channelID")
	number, _ := strconv.ParseUint(c.Query("number"), 0, 64)
	//TODO: behavior is not a bool, defined in pb
	behavior, _ := strconv.ParseBool(c.Query("behavior"))
	block, err := hs.cc.FetchBlock(channelID, uint64(number), bool(behavior))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}
	c.JSON(http.StatusOK, gin.H{"block": block.Bytes})
	return
}

// ListChannelReq Binding from JSON
type ListChannelReq struct {
	System string `form:"system" json:"system" xml:"system"  binding:"required"`
	PK     string `form:"pk" json:"pk" xml:"pk" binding:"required"`
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
	log.Info(json.PK)
	req := &pb.ListChannelsRequest{
		System: system,
		PK:     pk,
	}
	info, err := hs.cc.ListChannels(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"channelinfo": info})
	return
}

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
		pk, err := crypto.NewPublicKey(coreTx.Data.Sig.PK)
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
