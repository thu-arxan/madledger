package server

import (
	"github.com/gin-gonic/gin"
	"madledger/common/crypto"
	"madledger/core"
	"strconv"
)

func ChannelList(c *gin.Context) {
	Client, err := GetClient(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	sysstr := c.DefaultQuery("system", "true")
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

func CreateChannel(c *gin.Context) {
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

	publicstr := c.DefaultQuery("public", "true")
	public, err := strconv.ParseBool(publicstr)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	var admins []*core.Member
	adminsarr := c.QueryArray("admins")
	for _, adminpk := range adminsarr {
		pk, err := crypto.NewPublicKey([]byte(adminpk), crypto.KeyAlgoSM2)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		admin, err := core.NewMember(pk, "admin")
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		admins = append(admins, admin)
	}

	var members []*core.Member
	membersarr := c.QueryArray("members")
	for _, memberpk := range membersarr {
		pk, err := crypto.NewPublicKey([]byte(memberpk), crypto.KeyAlgoSM2)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		member, err := core.NewMember(pk, "member")
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		members = append(members, member)
	}

	gasPricestr := c.DefaultQuery("gasprice", "0")
	gasPrice, err := strconv.ParseUint(gasPricestr, 10, 64)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	ratiostr := c.DefaultQuery("ratio", "1")
	ratio, err := strconv.ParseUint(ratiostr, 10, 64)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	maxGasstr := c.DefaultQuery("maxGas", "10000000")
	maxGas, err := strconv.ParseUint(maxGasstr, 10, 64)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	err = Client.CreateChannel(channelID, public, admins, members, gasPrice, ratio, maxGas)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, nil)
}
