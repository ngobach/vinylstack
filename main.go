package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
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
	for _, url := range playlistUrls {
		pl, err := getPlaylist(url)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		pl.list = []song{}
	}
}
