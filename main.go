package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/thanbaiks/vinylstack/downloader"
)

func main() {
	csn := flag.String("csn", "", "ChiaSeNhac user ID")
	var dld downloader.Downloader
	flag.Parse()
	switch {
	case len(*csn) > 0:
		dld = downloader.ChiaSeNhac{UserID: *csn}
	}

	if dld == nil {
		flag.Usage()
		os.Exit(1)
	}
	color.Blue("Downloading playlists with downloader")
	fmt.Println(dld.Info())
	// Start download
	dld.Download()
}
