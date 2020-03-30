package core

type Genre struct {
	ID    ID      `json:"id"`
	Name  string  `json:"name"`
	Cover *string `json:"cover"`
}
