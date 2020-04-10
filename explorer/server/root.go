package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var SessionPool = make(map[string] *Session)
var defaultConfig = ""

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}


func RunServer(host string, port int, config string) error {
	defaultConfig = config
	fmt.Printf("Server run on %s:%d\n", host, port)
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.Static("/static", "./static")
	r.GET("/api/client/account/list", AccountList)
	r.GET("/api/client/channel/list", ChannelList)
	r.GET("/api/client/tx/history", TxHistory)
	r.GET("/api/client/network/info", NetworkInfo)
	r.Run(fmt.Sprintf("%s:%d", host, port))
	return nil
}
