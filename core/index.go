package core

type Index struct {
	Genres       []struct{} `json:"genres"`
	Artists      []struct{} `json:"artists"`
	Tracks       []Track    `json:"tracks"`
	DefautlCover string     `json:"default_cover"`
}
