package api

import "fmt"

func SongGetUrl(id int) string {
	r, err := client.Get(httpRoot + fmt.Sprintf("/song/media/outer/url?id=%d.mp3", id))
	if err != nil {
		return ""
	}
	logger.Debugf("get song url response : %+v", r.Header)
	return r.Header["Location"][0]
}

func SongLike(id int,like bool) bool {
	r := post(httpRoot+"/weapi/radio/like", map[string]interface{}{
		"trackId": id,
		"like":    like,
	})
	if r == nil || int(r["code"].(float64)) != 200 {
		return false
	}
	return true
}