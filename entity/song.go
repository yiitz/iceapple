package entity

type Song struct {
	Name   string
	Artist string
	Uri    string
	Album  string
}

func (s *Song) GetName() string {
	return s.Name
}

func (s *Song) GetUri() string {
	return s.Uri
}

func (s *Song) GetArtist() string {
	return s.Artist
}

func (s *Song) GetAlbum() string {
	return s.Album
}
