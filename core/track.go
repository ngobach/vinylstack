package core

type Track struct {
	URL    string   `json:"url"`
	Title  string   `json:"title"`
	Artist string   `json:"artist"`
	Cover  string   `json:"cover"`
	Genres []string `json:"genres"`
}
