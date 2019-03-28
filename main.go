package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	// Check if argument is empty
	if len(args) != 1 {
		fmt.Println("You must provide one CSN User ID")
		os.Exit(1)
	}
	fmt.Println("[*] User to be queried: " + args[0])
	playlistUrls, err := getPlaylistUrls(args[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("[*] Found playlists")
	songsByID := map[string]song{}
	for _, url := range playlistUrls {
		pl, err := getPlaylist(url)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		for _, ss := range pl.list {
			s, exists := songsByID[ss.id()]
			if exists {
				s.Genres = append(s.Genres, pl.title)
			} else {
				ss.Genres = append(ss.Genres, pl.title)
				songsByID[ss.id()] = ss
			}
		}
	}

	songs := []song{}
	for _, s := range songsByID {
		songs = append(songs, s)
	}

	e := exporter{"dist"}
	e.prepare()
	err = e.export(songs)
	if err != nil {
		fmt.Println(err)
	}
}
