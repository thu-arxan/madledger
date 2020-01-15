package server

import (
	"encoding/json"
	"madledger/common"
	"madledger/common/crypto"
	"madledger/core"
	pb "madledger/protos"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FetchBlock gets block by http
func (hs *HTTPServer) FetchBlock(c *gin.Context) {
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

// ListChannels list channels by http
func (hs *HTTPServer) ListChannels(c *gin.Context) {
	system, _ := strconv.ParseBool(c.Query("system"))
	pk := []byte(c.Query("pk"))
	req := &pb.ListChannelsRequest{
		System: system,
		PK:     pk,
	}
	info, err := hs.cc.ListChannels(req)
	c.JSON(http.StatusOK, gin.H{"channelinfo": info.String, "error": err.Error})
	return
}

// CreateChannel create channel by http
func (hs *HTTPServer) CreateChannel(c *gin.Context) {
	tx := []byte(c.Query("tx"))
	var coreTx core.Tx
	json.Unmarshal(tx, &coreTx)
	if !coreTx.Verify() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The tx is not a valid tx"})
		return
	}
	if coreTx.GetReceiver().String() != core.CreateChannelContractAddress.String() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The receiver of the tx is not the valid contract address"})
		return
	}
	info, err := hs.cc.CreateChannel(&coreTx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
	}
	c.JSON(http.StatusOK, gin.H{"channelinfo": info.String, "error": err.Error})
	return
}

// AddTx add tx by http
func (hs *HTTPServer) AddTx(c *gin.Context) {
	var status pb.TxStatus
	tx := []byte(c.Query("tx"))
	var coreTx core.Tx
	json.Unmarshal(tx, &coreTx)

	txType, err := core.GetTxType(common.BytesToAddress(coreTx.Data.Recipient).String())
	if err == nil && (txType == core.VALIDATOR || txType == core.NODE) {
		pk, err := crypto.NewPublicKey(coreTx.Data.Sig.PK)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}
		// create member to check if the client is system admin
		member, err := core.NewMember(pk, "")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}
		if !hs.cc.CM.IsSystemAdmin(member) { // not system admin, return error
			c.JSON(http.StatusBadRequest, gin.H{"error": "The client is not system admin and can't config the cluster"})
			return
		}
	}
	err = hs.cc.AddTx(&coreTx)
	c.JSON(http.StatusOK, gin.H{"txstatus": status.String, "error": err.Error})
	return
}
