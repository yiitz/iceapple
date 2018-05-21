package api

import (
	"net/http"
	"github.com/yiitz/persistent-cookiejar"
	"github.com/yiitz/iceapple/storage"
	"github.com/yiitz/iceapple/log"
)

var client http.Client

var logger = log.NewLogger("api")

func init() {
	jar, err := cookiejar.New(&cookiejar.Options{Filename: storage.AppDir() + "/.cookie"})
	if err != nil {
		panic(err)
	}
	client.Jar = jar
}

func SaveCookie() {
	client.Jar.(*cookiejar.Jar).Save()
}
