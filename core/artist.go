package core

type Artist struct {
	ID    ID      `json:"id"`
	Name  string  `json:"name"`
	Cover *string `json:"cover"`
}
