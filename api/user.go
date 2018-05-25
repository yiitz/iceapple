package api

import (
	"fmt"
	"crypto/md5"
)

func UserLoginPhone(u string, p string) bool {
	r := post(httpRoot+"/weapi/login/cellphone", map[string]interface{}{
		"phone":         u,
		"password":      fmt.Sprintf("%x", md5.Sum([]byte(p))),
		"rememberLogin": "true",
	})
	if r == nil || int(r["code"].(float64)) != 200{
		return false
	}

	return true
}

func UserLogin(u string, p string) bool {
	r := post(httpRoot+"/weapi/login", map[string]interface{}{
		"username":         u,
		"password":      fmt.Sprintf("%x", md5.Sum([]byte(p))),
		"rememberLogin": "true",
	})
	if r == nil || int(r["code"].(float64)) != 200{
		return false
	}

	return true
}

func UserInfo() (r map[string]interface{},login bool) {
	login = false
	r = post(httpRoot+"/weapi/subcount",nil)
	if r == nil || int(r["code"].(float64)) == 301{
		return
	}
	login = true
	return
}