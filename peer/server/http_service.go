package server

import (
	"encoding/hex"
	pb "madledger/protos"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTxStatusReq ...
type GetTxStatusReq struct {
	ChannelID string `json:"channelID"`
	TxID      string `json:"txID"`
}

// GetTxStatusByHTTP gets tx status by http
func (hs *Server) GetTxStatusByHTTP(c *gin.Context) {
	log.Info("start get tx status")
	var j GetTxStatusReq
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	chID := j.ChannelID
	txID := j.TxID

	log.Infof("before get tx status %s", txID)
	status, err := hs.cm.GetTxStatus(chID, txID, true)
	log.Infof("after get tx status %s, %v", txID, err)

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
	log.Info("finish get tx status")
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
	log.Info("addr is ", hex.EncodeToString(addr))
	history := hs.cm.GetTxHistory(addr)
	log.Info("get history is ", history)
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
