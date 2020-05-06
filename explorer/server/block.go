package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetBlock get block
func GetBlock(c *gin.Context) {
	Client, err := GetClient(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	channelID := c.DefaultQuery("channelid", "")
	if channelID == "" {
		c.AbortWithError(500, err)
		return
	}
	blockIndexstr := c.DefaultQuery("blockindex", "0")
	blockIndex, err := strconv.ParseUint(blockIndexstr, 10, 64)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	block, err := Client.GetBlock(blockIndex, channelID)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, gin.H{
		"block": block,
	})
}
