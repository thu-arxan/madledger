package server

import (
	"github.com/gin-gonic/gin"
	pb "madledger/protos"
	"net/http"
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

	status, err := hs.ChannelManager.GetTxStatus(chID, txID, true)
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
	addr := []byte(j.Addr)

	history := hs.ChannelManager.ListTxHistory(addr)
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
