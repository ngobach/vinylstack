package core

// Playlist describe a playlist
type Playlist struct {
	Title string
	List  []Song
}

// Simplify merge playlists into songs with appended genres
func Simplify(pll []Playlist) []Song {
	return []Song{}
}
