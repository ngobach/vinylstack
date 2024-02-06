package core

type Index struct {
	Genres       []struct{} `json:"genres"`
	Artists      []struct{} `json:"artists"`
	Tracks       []Track    `json:"tracks"`
	DefautlCover string     `json:"default_cover"`
}

type Track struct {
	URL    string   `json:"url"`
	Title  string   `json:"title"`
	Artist string   `json:"artist"`
	Cover  string   `json:"cover"`
	Genres []string `json:"genres"`
}
