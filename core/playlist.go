package core

// Playlist describe a playlist
type Playlist struct {
	Title string
	List  []Song
}

// Simplify merge playlists into songs with appended genres
func Simplify(plls []Playlist) []Song {
	m := map[string]Song{}
	for _, playlist := range plls {
		for _, song := range playlist.List {
			if _, found := m[song.ID()]; found {
				// Use the existing one instead
				song = m[song.ID()]
			}
			song.Genres = append(song.Genres, playlist.Title)
			m[song.ID()] = song
		}
	}
	songs := make([]Song, 0, len(m))
	for _, song := range m {
		songs = append(songs, song)
	}
	return songs
}
