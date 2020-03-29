package core

type Artist struct {
	Id     string
	Name   string
	Cover  *string
	Tracks []*Track
}
