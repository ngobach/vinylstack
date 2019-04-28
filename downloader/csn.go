package downloader

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/thanbaiks/vinylstack/core"
)

const concurrentDownload = 16

// ChiaSeNhac downloader
type ChiaSeNhac struct {
	UserID string
}

// Info return info of the downloader
func (csn ChiaSeNhac) Info() string {
	return fmt.Sprintf("ChiaSeNhac downloader (UID: %s)", csn.UserID)
}

// Download start download process
func (csn ChiaSeNhac) Download() ([]core.Playlist, error) {
	pllURLs, err := getPlaylistUrls(csn.UserID)
	if err != nil {
		return nil, err
	}

	playlists := []core.Playlist{}
	for _, pllURL := range pllURLs {
		pl, err := getPlaylist(pllURL)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, pl)
	}
	return playlists, nil
}

func getPlaylistUrls(uid string) ([]string, error) {
	result := []string{}
	doc, err := goquery.NewDocument("https://chiasenhac.vn/user/" + uid)
	if err != nil {
		return nil, err
	}
	doc.Find("#playlist .card").Each(func(idx int, el *goquery.Selection) {
		href, _ := el.Find(".card-title>a").Attr("href")
		result = append(result, href)
	})
	return result, nil
}

func getPlaylist(url string) (core.Playlist, error) {
	pl := core.Playlist{}
	doc, err := goquery.NewDocument("https://chiasenhac.vn" + url)
	if err != nil {
		return pl, err
	}
	pl.Title = doc.Find("meta[name=title]").AttrOr("content", "Untitled")
	fmt.Println("Scanning playlist", pl.Title)
	plSize := doc.Find(".card-footer").Length()
	ids := make(chan int, plSize)
	fetched := make(chan core.Song)
	for i := 1; i <= plSize; i++ {
		ids <- i
	}
	close(ids)
	for i := 0; i < concurrentDownload; i++ {
		go func() {
			for {
				idx, found := <-ids
				if !found {
					return
				}
				s := core.Song{}
				doc, err := goquery.NewDocument("https://chiasenhac.vn" + url + "?playlist=" + strconv.Itoa(idx))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				s.Title = doc.Find(".card-details .card-title").Text()
				s.Artist = doc.Find(".card-details .list-unstyled li").First().Text()
				s.Artist = strings.TrimPrefix(s.Artist, "Ca sĩ: ")
				s.Cover, _ = doc.Find(".card-details img").Attr("src")
				if s.Cover == "/imgs/no_cover.jpg" {
					s.Cover = ""
				}
				// Print song info
				doc.Find(".download_item").Each(func(idx int, el *goquery.Selection) {
					if strings.Contains(el.Text(), "128") {
						s.URL, _ = el.Attr("href")
					}
				})
				fetched <- s
			}
		}()
	}

	for i := 0; i < plSize; i++ {
		s := <-fetched
		fmt.Printf("-- %s (%s)\n", s.Title, s.Artist)
		pl.List = append(pl.List, s)
	}

	return pl, nil
}
