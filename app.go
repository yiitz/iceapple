package main

import (
	"github.com/yiitz/iceapple/app"
	"github.com/yiitz/iceapple/log"
	"github.com/yiitz/iceapple/media"
)

func main() {
	media.Init()
	app.Run()
	media.Destroy()
	log.Flush()
}
