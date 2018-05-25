package app

import (
	"github.com/yiitz/iceapple/api"
	"github.com/yiitz/iceapple/entity"
	"github.com/yiitz/iceapple/ui"
	"strings"
	"time"
)

var songs []ui.PlayListItem

func enterPersonalFM() {

	pl.Selectable = false

	pb.OnSongFinished = func() {
		playNext()
	}

	playNextFunc = playNext

	playCurrent()
}

func playNext() {
	songs = songs[1:]
	playCurrent()
}

func playCurrent() {

	for len(songs) <= 0 {
		queryNext()
		if len(songs) <= 0 {
			time.Sleep(time.Second * 3)
		}
	}

	s := songs[0]
	pl.SetItems(songs)
	player.Play(s.GetUri())
}

func queryNext() {
	l := api.PersonalFM()
	if l == nil {
		return
	}
	for _, v := range l {
		v := v.(map[string]interface{})
		s := &entity.Song{Name: v["name"].(string)}
		var as []string
		for _, a := range v["artists"].([]interface{}) {
			as = append(as, a.(map[string]interface{})["name"].(string))
		}
		s.Artist = strings.Join(as, ",")
		s.Album = (v["album"].(map[string]interface{}))["name"].(string)
		s.Id = int(v["id"].(float64))
		songs = append(songs, s)
	}
	return
}
