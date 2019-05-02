package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/thanbaiks/vinylstack/core"
	"github.com/thanbaiks/vinylstack/exporter"

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

	begin := time.Now()
	playlists, err := dld.Download()
	if err != nil {
		panic(err)
	}
	songs := core.Simplify(playlists)
	totalTime := time.Since(begin)
	color.Green("Fetched %d playlists with %d songs (in %s)", len(playlists), len(songs), totalTime.String())
	color.Blue("Start downloading")
	exporter := exporter.Exporter{Target: "_dist_"}
	err = exporter.Prepare()
	if err != nil {
		panic(err)
	}
	begin = time.Now()
	err = exporter.DownloadAndExport(songs)
	totalTime = time.Since(begin)
	if err != nil {
		panic(err)
	}
	color.Green("Finished downloading (in %s)", totalTime.String())
	color.White("Have fun!")
}
