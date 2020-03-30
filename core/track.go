package core

type Track struct {
	ID        ID      `json:"id"`
	URL       string  `json:"url"`
	Title     string  `json:"title"`
	ArtistIds []ID    `json:"artists"`
	Cover     *string `json:"cover"`
	GenreIds  []ID    `json:"genres"`
}
