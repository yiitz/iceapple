package main

import (
	"github.com/yiitz/iceapple/app"
	"github.com/yiitz/iceapple/log"
	"github.com/yiitz/iceapple/media"
	"github.com/yiitz/iceapple/api"
)

func main() {
	media.Init()
	app.Run()
	media.Destroy()
	log.Flush()
	api.SaveCookie()
}
