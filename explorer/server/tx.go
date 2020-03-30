package server

import (
	"github.com/gin-gonic/gin"
	"madledger/common"
)


func TxHistory(c *gin.Context) {
	Client, err := GetClient(c)
	var address common.Address
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if address, err = Client.GetPrivKey().PubKey().Address(); err != nil {
		c.AbortWithError(500, err)
		return
	}

	history, err := Client.GetHistory(address.Bytes())
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, gin.H{
		"history": history.Txs,
	})
}