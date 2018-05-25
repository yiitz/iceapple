package main

import (
	"github.com/yiitz/iceapple/config"
	"github.com/yiitz/iceapple/api"
	"github.com/yiitz/iceapple/app"
	"github.com/yiitz/iceapple/log"
	"github.com/yiitz/iceapple/media"
	"flag"
	"os"
	"fmt"
	"bufio"
	"syscall"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	flag.StringVar(&config.Proxy, "x", "", "set http proxy")
	flag.StringVar(&config.LogLevel, "l", "debug", "set log level,debug,info,warn...")
}

func main() {
	flag.Parse()
	if len(config.Proxy) <=0 {
		p, ok := os.LookupEnv("http_proxy")
		if ok {
			config.Proxy = p
		}
	}
	log.SetLevel(config.LogLevel)
	api.InitClient()
	media.Init()

	fmt.Println("check login state...")
	for _,ok := api.UserInfo();!ok; {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("need login, choose method to login.\n1: cellphone\n2: email\nyour choice:")
		method, _ := reader.ReadString('\n')
		fmt.Print("account:")
		u, _ := reader.ReadString('\n')
		fmt.Print("password:")
		p, _ := terminal.ReadPassword(int(syscall.Stdin))
		if "1" == method {
			api.UserLoginPhone(u, string(p))
		} else {
			api.UserLogin(u, string(p))
		}
	}

	app.Run()
	media.Destroy()
	api.SaveCookie()
	log.Flush()
}
