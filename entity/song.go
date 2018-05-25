package entity

import "github.com/yiitz/iceapple/api"

type Song struct {
	Id     int
	Name   string
	Artist string
	Uri    string
	Album  string
}

func (s *Song) GetName() string {
	return s.Name
}

func (s *Song) GetUri() string {
	return api.SongGetUrl(s.Id)
}

func (s *Song) GetArtist() string {
	return s.Artist
}

func (s *Song) GetAlbum() string {
	return s.Album
}
