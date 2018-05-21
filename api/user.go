package api

import (
	"net/url"
	"encoding/json"
)

func UserLoginPhone(u string, p string) bool {
	resp, err := client.PostForm("http://music.163.com/weapi/login/cellphone", url.Values{
		"phone":         []string{u},
		"password":      []string{p},
		"rememberLogin": []string{"true"},
	})
	if err != nil {
		logger.Error(err)
	}
	r := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(r)
	if err == nil {
		logger.Error(err)
	}
	logger.Debugf("login response:%+v", r)
	return true
}

func UserInfo()  {

}