package core

// Playlist describe a playlist
type Genres struct {
	ID     string
	Name   string
	Cover  *string
	Tracks []*Track
}

// Simplify merge playlists into songs with appended genres
func Simplify(playlists []Playlist) []Track {
	m := map[string]Track{}
	for _, playlist := range playlists {
		for _, song := range playlist.List {
			if _, found := m[song.ID()]; found {
				// Use the existing one instead
				song = m[song.ID()]
			}
			song.Genres = append(song.Genres, playlist.Title)
			m[song.ID()] = song
		}
	}
	songs := make([]Track, 0, len(m))
	for _, song := range m {
		songs = append(songs, song)
	}
	return songs
}
