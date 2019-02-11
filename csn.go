package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type song struct {
	url    string
	title  string
	artist string
	cover  string
}
type playlist struct {
	title string
	list  []song
}

func getPlaylistUrls(uid string) ([]string, error) {
	result := []string{}
	doc, err := goquery.NewDocument("https://beta.chiasenhac.vn/user/" + uid)
	if err != nil {
		return nil, err
	}
	doc.Find("#playlist .card").Each(func(idx int, el *goquery.Selection) {
		href, _ := el.Find(".card-title>a").Attr("href")
		result = append(result, href)
	})
	return result, nil
}

func getPlaylist(url string) (playlist, error) {
	pl := playlist{}
	doc, err := goquery.NewDocument("https://beta.chiasenhac.vn" + url)
	if err != nil {
		return pl, err
	}
	pl.title = doc.Find("meta[name=title]").AttrOr("content", "untitled")
	fmt.Println("[*] Found playlist", pl.title)
	plSize := doc.Find(".card-footer").Length()
	for i := 1; i <= plSize; i++ {
		s := song{}
		fmt.Println("----")
		doc, err := goquery.NewDocument("https://beta.chiasenhac.vn" + url + "?playlist=" + strconv.Itoa(i))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		s.title = doc.Find(".card-details .card-title").Text()
		s.artist = doc.Find(".card-details .list-unstyled li").First().Text()
		s.artist = strings.TrimPrefix(s.artist, "Ca sĩ: ")
		s.cover, _ = doc.Find(".card-details img").Attr("src")
		if s.cover == "/imgs/no_cover.jpg" {
			s.cover = ""
		}
		// Print song info
		doc.Find(".download_item").Each(func(idx int, el *goquery.Selection) {
			if strings.Contains(el.Text(), "320") {
				s.url, _ = el.Attr("href")
			}
		})
		fmt.Println("[+] Name:", s.title)
		fmt.Println("[+] Artist:", s.artist)
		fmt.Println("[+] Cover:", s.cover)
		fmt.Println("[+] 320bps:", s.url)
		pl.list = append(pl.list, s)
	}
	return pl, nil
}
