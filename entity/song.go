package entity

import "github.com/yiitz/iceapple/api"

type Song struct {
	Id      int
	Name    string
	Artist  string
	Album   string
	Starred bool
}

func (s *Song) GetUri() string {
	return api.SongGetUrl(s.Id)
}