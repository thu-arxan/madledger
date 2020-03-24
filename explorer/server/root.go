package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	clib "madledger/client/lib"
)

var Client *clib.Client = nil

func RunServer(host string, port int, config string) error {
	var err error
	fmt.Printf("Server run on %s:%d\n", host, port)
	Client, err = clib.NewClient(config)
	if err != nil {
		return err
	}
	r := gin.Default()
	r.GET("/api/client/account/list", AccountList)
	r.Run(fmt.Sprintf("%s:%d", host, port))
	return nil
}