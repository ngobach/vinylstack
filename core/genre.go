package core

type Genre struct {
	ID    string
	Name  string
	Cover *string
}

// TODO: Remove this
func Simplify(playlists []Genre) []Track {
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
