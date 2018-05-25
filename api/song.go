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