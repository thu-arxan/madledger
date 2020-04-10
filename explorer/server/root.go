package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var SessionPool = make(map[string]*Session)
var defaultConfig = ""

func RunServer(host string, port int, config string) error {
	defaultConfig = config
	fmt.Printf("Server run on %s:%d\n", host, port)
	r := gin.Default()
	r.GET("/api/client/account/list", AccountList)
	r.GET("/api/client/channel/list", ChannelList)
	r.GET("/api/client/tx/history", TxHistory)
	r.GET("/api/client/network/info", NetworkInfo)
	r.GET("/api/client/channel/create", CreateChannel)
	r.Run(fmt.Sprintf("%s:%d", host, port))
	return nil
}
