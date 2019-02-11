package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type song struct {
	Url    string   `json:"url"`
	Title  string   `json:"title"`
	Artist string   `json:"artist"`
	Cover  string   `json:"cover"`
	Genres []string `json:"genres"`
}

type playlist struct {
	title string
	list  []song
}

func (s *song) id() string {
	h := md5.New()
	h.Write([]byte(s.Title))
	h.Write([]byte(s.Artist))
	return hex.EncodeToString(h.Sum(nil))[:8]
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
		s.Title = doc.Find(".card-details .card-title").Text()
		s.Artist = doc.Find(".card-details .list-unstyled li").First().Text()
		s.Artist = strings.TrimPrefix(s.Artist, "Ca sĩ: ")
		s.Cover, _ = doc.Find(".card-details img").Attr("src")
		if s.Cover == "/imgs/no_cover.jpg" {
			s.Cover = ""
		}
		// Print song info
		doc.Find(".download_item").Each(func(idx int, el *goquery.Selection) {
			if strings.Contains(el.Text(), "320") {
				s.Url, _ = el.Attr("href")
			}
		})
		fmt.Println("[+] Name:", s.Title)
		fmt.Println("[+] Artist:", s.Artist)
		fmt.Println("[+] Cover:", s.Cover)
		fmt.Println("[+] 320bps:", s.Url)
		pl.list = append(pl.list, s)
	}
	return pl, nil
}
