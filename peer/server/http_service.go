package server

import (
	"encoding/binary"
	"encoding/hex"
	"madledger/common"
	"madledger/common/util"
	pb "madledger/protos"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetTxStatusReq ...
type GetTxStatusReq struct {
	ChannelID string `json:"channelID"`
	TxID      string `json:"txID"`
}

// GetTxStatusByHTTP gets tx status by http
func (hs *Server) GetTxStatusByHTTP(c *gin.Context) {
	var j GetTxStatusReq
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	chID := j.ChannelID
	txID := j.TxID

	status, err := hs.cm.GetTxStatus(chID, txID, true)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := &pb.TxStatus{
		Err:             status.Err,
		BlockNumber:     status.BlockNumber,
		BlockIndex:      int32(status.BlockIndex),
		Output:          status.Output,
		ContractAddress: status.ContractAddress,
	}
	c.JSON(http.StatusOK, gin.H{"status": result})
	return
}

// ListTxHistoryReq ...
type ListTxHistoryReq struct {
	Addr string `json:"address"`
}

// ListTxHistoryByHTTP lists Tx history by http
func (hs *Server) ListTxHistoryByHTTP(c *gin.Context) {
	var j ListTxHistoryReq
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addr, _ := hex.DecodeString(j.Addr)
	history := hs.cm.GetTxHistory(addr)
	var pbHistory = make(map[string]*pb.StringList)
	for channelID, ids := range history {
		value := new(pb.StringList)
		for _, id := range ids {
			value.Value = append(value.Value, id)
		}
		pbHistory[channelID] = value
	}
	res := &pb.TxHistory{
		Txs: pbHistory,
	}
	c.JSON(http.StatusOK, gin.H{"txhistory": res})
	return
}

// GetTokenInfoReq ...
type GetTokenInfoReq struct {
	Addr      string `json:"address"`
	ChannelID string `json:"channelid"`
}

// GetTokenInfoByHTTP Get Token Info By HTTP
func (hs *Server) GetTokenInfoByHTTP(c *gin.Context) {
	var j GetTokenInfoReq
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	channelID := j.ChannelID
	addr, err := hex.DecodeString(j.Addr)
	key := util.BytesCombine(common.AddressFromChannelID(channelID).Bytes(), []byte("token"), addr)
	tokenBytes, err := hs.cm.db.Get(key, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var token uint64
	if tokenBytes != nil {
		token = uint64(binary.BigEndian.Uint64(tokenBytes))
	}
	info := &pb.TokenInfo{
		Balance: token,
	}
	c.JSON(http.StatusOK, gin.H{"tokeninfo": info})
	return
}

//GetBlockReq ...
type GetBlockReq struct {
	ChannelID string `json:"channelid"`
	Num       string `json:"num"`
}

//GetBlockByHTTP Get Block By HTTP
func (hs *Server) GetBlockByHTTP(c *gin.Context) {
	var j GetBlockReq
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	channelID := j.ChannelID
	num, err := strconv.ParseUint(j.Num, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	block, err := hs.cm.db.GetBlock(channelID, num)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"block": block})
	return
}
