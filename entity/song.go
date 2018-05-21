package entity

type Song struct {
	Name   string
	Artist string
	Uri    string
}

func (s *Song) GetName() string {
	return s.Name
}

func (s *Song) GetUri() string {
	return s.Uri
}
