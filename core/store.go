package core

import "fmt"

type Store struct {
	Genres  map[ID]Genre  `json:"genres"`
	Artists map[ID]Artist `json:"artists"`
	Tracks  map[ID]Track  `json:"tracks"`
}

var DefaultStore = Store{
	Genres:  map[ID]Genre{},
	Artists: map[ID]Artist{},
	Tracks:  map[ID]Track{},
}

func (store *Store) Dump() {
	for _, g := range store.Genres {
		fmt.Printf("Genre: %s - %s\n", g.ID, g.Name)
	}
	for _, a := range store.Artists {
		fmt.Printf("Artist: %s - %s\n", a.ID, a.Name)
	}
	for _, t := range store.Tracks {
		sub := ""
		for _, g := range t.GenreIds {
			sub += "G:" + string(g) + ";"
		}
		for _, a := range t.ArtistIds {
			sub += "A:" + string(a) + ";"
		}
		fmt.Printf("Track: %s - %s (%s)\n", t.ID, t.Title, sub)
	}
}
