package server

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func ChannelList(c *gin.Context) {
	Client, err := GetClient(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	sysstr := c.DefaultQuery("system","true")
	sysflag, err := strconv.ParseBool(sysstr)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	channels, err := Client.ListChannel(sysflag)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, gin.H{
		"channels": channels,
	})
}
