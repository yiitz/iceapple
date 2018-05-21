package storage

import (
	"os/user"
	"os"
)

var appDir string
func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	appDir = usr.HomeDir + "/.iceapple"
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		if err = os.Mkdir(appDir, 0700); err != nil {
			panic(err)
		}
	}
}

func AppDir() string {
	return appDir
}