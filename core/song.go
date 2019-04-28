package core

import (
	"crypto/md5"
	"encoding/hex"
)

// Song describe a song
type Song struct {
	URL    string   `json:"url"`
	Title  string   `json:"title"`
	Artist string   `json:"artist"`
	Cover  string   `json:"cover"`
	Genres []string `json:"genres"`
}

// ID return hashed ID of the song
func (s *Song) ID() string {
	h := md5.New()
	h.Write([]byte(s.Title))
	h.Write([]byte(s.Artist))
	return hex.EncodeToString(h.Sum(nil))[:8]
}
