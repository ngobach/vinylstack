package patcher

import (
	"encoding/json"
	"github.com/thanbaiks/vinylstack/core"
	"io/ioutil"
	"os"
)

type PatchList struct {
	Genres       map[core.ID]core.Genre
	Artists      map[core.ID]core.Artist
	Tracks       map[core.ID]core.Track
	DefaultCover *string `json:"default_cover,omitempty"`
}

func LoadPatchList(file string) *PatchList {
	_, err := os.Stat(file)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if os.IsNotExist(err) {
		return nil
	}
	doc, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	pl := PatchList{
		Genres:  map[core.ID]core.Genre{},
		Artists: map[core.ID]core.Artist{},
		Tracks:  map[core.ID]core.Track{},
	}
	err = json.Unmarshal(doc, &pl)
	if err != nil {
		panic(err)
	}
	return &pl
}

func (pl *PatchList) Patch(store *core.Store) {
	for id, track := range pl.Tracks {
		if target, ok := store.Tracks[id]; ok {
			if len(track.Title) > 0 {
				target.Title = track.Title
			}
			if len(*track.Cover) > 0 {
				target.Cover = track.Cover
			}
			if len(track.URL) > 0 {
				target.URL = track.URL
			}
			store.Tracks[id] = target
		}
	}
	for id, genre := range pl.Genres {
		if target, ok := store.Genres[id]; ok {
			if len(genre.Name) > 0 {
				target.Name = genre.Name
			}
			if len(*genre.Cover) > 0 {
				target.Cover = genre.Cover
			}
			store.Genres[id] = target
		}
	}
	for id, artist := range pl.Artists {
		if target, ok := store.Genres[id]; ok {
			if len(artist.Name) > 0 {
				target.Name = artist.Name
			}
			if len(*artist.Cover) > 0 {
				target.Cover = artist.Cover
			}
			store.Genres[id] = target
		}
	}
	if len(*pl.DefaultCover) > 0 {
		store.DefaultCover = pl.DefaultCover
	}
}
