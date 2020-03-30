package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	clib "madledger/client/lib"
	"math/rand"
	"net/http"
	"time"
)

type Session struct {
	client *clib.Client
	lastAccessTime int64
}

func NewSession() (*Session, error){
	var result = new(Session)
	client, err := clib.NewClient(defaultConfig)
	if err != nil {
		return nil, err
	}
	result.client = client
	result.lastAccessTime = time.Now().Unix()
	return result, nil
}

func  GetRandomString(l int) string {
	str := "0123456789abcdef"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetClient(c *gin.Context) (*clib.Client, error) {
	if rand.Uint32() % 1 == 0 {
		go func() {
			tooOld := time.Now().Unix() - 10
			for k, v := range SessionPool {
				fmt.Println(v.lastAccessTime, tooOld)
				if v.lastAccessTime < tooOld {
					delete(SessionPool, k)
				}
			}
		}()
	}
	cookie, err := c.Cookie("client_id")
	if err == nil {
		client := SessionPool[cookie]
		if client != nil {
			fmt.Println("Use exist client ", cookie)
			SessionPool[cookie].lastAccessTime = time.Now().Unix()
			return SessionPool[cookie].client, nil
		}
	}
	cookie = GetRandomString(32)
	c.SetCookie("client_id", cookie, 3600, "/", "", http.SameSiteLaxMode, false, true)

	session, err := NewSession()
	if err != nil {
		return nil, err
	}
	SessionPool[cookie] = session
	fmt.Println("New client ", cookie)

	return session.client, nil
}