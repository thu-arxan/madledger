package server

import (
	"github.com/gin-gonic/gin"
	pb "madledger/protos"
	"net/http"
)

// GetTxStatus gets tx status by http
func (hs *HTTPServer) GetTxStatus(c *gin.Context) {
	chID := c.Query("channelID")
	txID := c.Query("txID")
	status, err := hs.ChannelManager.GetTxStatus(chID, txID, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}
	result := &pb.TxStatus{
		Err:             status.Err,
		BlockNumber:     status.BlockNumber,
		BlockIndex:      int32(status.BlockIndex),
		Output:          status.Output,
		ContractAddress: status.ContractAddress,
	}
	c.JSON(http.StatusOK, gin.H{"status": result.String})
	return
}

// ListTxHistory lists Tx history by http
func (hs *HTTPServer) ListTxHistory(c *gin.Context) {
	addr := []byte(c.Query("address"))
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
	c.JSON(http.StatusOK, gin.H{"txhistory": res.String})
	return
}
