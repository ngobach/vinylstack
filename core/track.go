package core

// Track describe a Media track
type Track struct {
	ID        *string   `json:"id"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	ArtistIds []*string `json:"artists"`
	Cover     string    `json:"cover"`
	GenreIds  []*string `json:"genres"`
}
