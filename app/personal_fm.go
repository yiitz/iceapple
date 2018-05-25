package app

import (
	"github.com/yiitz/iceapple/api"
	"github.com/yiitz/iceapple/ui"
	"github.com/yiitz/iceapple/entity"
	"strings"
)

var songs []ui.PlayListItem

func enterPersonalFM()  {

	pl.Selectable = false

	pl.SetItems(songs)

	pb.OnSongFinished = func() {
		songs = songs[1:]
		playNext()
	}

	playNextFunc = playNext

	playNext()
}

func playNext() {
	for len(songs) <= 0 {
		queryNext()
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
		s.Uri = api.SongGetUrl(int(v["id"].(float64)))
		songs = append(songs, s)
	}
	return
}