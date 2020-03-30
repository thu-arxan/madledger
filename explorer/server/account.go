package server

import (
	"github.com/gin-gonic/gin"
)

func AccountList(c *gin.Context) {
	Client, err := GetClient(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	address, err := Client.GetPrivKey().PubKey().Address()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	balance, err := Client.GetAccountBalance(address)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, gin.H{
		"address": address.String(),
		"balance": balance,
	})
}
