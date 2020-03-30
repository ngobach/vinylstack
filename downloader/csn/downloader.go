package csn

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/thanbaiks/vinylstack/core"
	"strconv"
	"strings"
	"sync"
)

const concurrentDownload = 4

type ChiaSeNhac struct {
	UserID string
}

func toID(s string) core.ID {
	return core.ID("csn-" + s)
}

func (csn *ChiaSeNhac) Download(store *core.Store) error {
	pllURLs, err := getUserPlaylists(csn.UserID)
	if err != nil {
		return nil
	}

	for _, pllURL := range pllURLs {
		err := getPlaylist(pllURL, store)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractPlaylistID(s string) core.ID {
	return toID(strings.Split(s, "/")[2])
}

func extractPlaylistCover(url string) *string {
	if len(url) == 0 || url == "https://data.chiasenhac.com/imgs/no_cover.jpg" {
		return nil
	}
	return &url
}

func extractArtistCover(url string) *string {
	if len(url) == 0 || strings.Contains(url, "no_cover") {
		return nil
	}
	return &url
}

func extractTrackCover(url string) *string {
	if len(url) == 0 || strings.Contains(url, "no_cover") {
		return nil
	}
	return &url
}

func getUserPlaylists(uid string) ([]string, error) {
	var result []string
	doc, err := goquery.NewDocument("https://chiasenhac.vn/user/" + uid)
	if err != nil {
		return nil, err
	}
	doc.Find("#playlist .card").Each(func(_ int, el *goquery.Selection) {
		href, _ := el.Find(".card-title>a").Attr("href")
		result = append(result, href)
	})
	return result, nil
}

func getPlaylist(url string, store *core.Store) error {
	fmt.Println("[csn] Downloading", url)
	doc, err := goquery.NewDocument("https://chiasenhac.vn" + url)
	if err != nil {
		return err
	}
	genre := core.Genre{
		ID:    extractPlaylistID(url),
		Name:  doc.Find("meta[name=title]").AttrOr("content", "untitled"),
		Cover: extractPlaylistCover(doc.Find("link[rel=image_src]").AttrOr("href", "")),
	}
	store.Genres[genre.ID] = genre
	plSize := doc.Find(".card-footer").Length()
	ids := make(chan int, plSize)
	fetched := make(chan core.Track)
	for i := 1; i <= plSize; i++ {
		ids <- i
	}
	close(ids)
	storeMutex := sync.Mutex{}
	for i := 0; i < concurrentDownload; i++ {
		go func() {
			for idx := range ids {
				track := core.Track{}
				songURL := "https://chiasenhac.vn" + url + "?playlist=" + strconv.Itoa(idx)
				doc, err := goquery.NewDocument(songURL)
				if err != nil {
					panic(err)
				}
				track.Title = doc.Find(".card-details .card-title").Text()
				track.ID = toID(core.GenString(track.Title))
				storeMutex.Lock()
				localTrack, ok := store.Tracks[track.ID]
				if ok {
					localTrack.GenreIds = append(localTrack.GenreIds, genre.ID)
					store.Tracks[track.ID] = localTrack
					storeMutex.Unlock()
					fetched <- localTrack
					continue
				}
				storeMutex.Unlock()
				track.GenreIds = []core.ID{genre.ID}
				if len(track.Title) == 0 {
					panic(fmt.Errorf("found empty title while fetching track %s", track.ID))
				}
				storeMutex.Lock()
				doc.Find(".card-details .list-unstyled li").First().Find("a").Each(func(_ int, selection *goquery.Selection) {
					artistUrl, _ := selection.Attr("href")
					artistID := toID(core.GenString(artistUrl))
					track.ArtistIds = append(track.ArtistIds, artistID)
					_, ok := store.Artists[artistID]
					if !ok {
						artist, err := getArtist(artistUrl, artistID)
						if err != nil {
							panic(fmt.Errorf("cant get artist information %s %w", artistID, err))
						}
						store.Artists[artistID] = *artist
					}
				})
				storeMutex.Unlock()
				trackCover, _ := doc.Find(".card-details img").Attr("src")
				track.Cover = extractTrackCover(trackCover)
				doc.Find(".download_item").Each(func(_ int, el *goquery.Selection) {
					if strings.Contains(el.Text(), "128") {
						track.URL, _ = el.Attr("href")
					}
				})
				storeMutex.Lock()
				store.Tracks[track.ID] = track
				storeMutex.Unlock()
				fetched <- track
			}
		}()
	}

	for i := 0; i < plSize; i++ {
		<-fetched
	}
	return nil
}

func getArtist(url string, id core.ID) (*core.Artist, error) {
	if strings.Contains(url, "tim-kiem") {
		return &core.Artist{
			ID:    id,
			Name:  "Various artist",
			Cover: nil,
		}, nil
	}
	doc, err := goquery.NewDocument("https://chiasenhac.vn" + url)
	if err != nil {
		return nil, err
	}
	artist := core.Artist{
		ID: id,
	}
	artist.Name = doc.Find(".artist_name_box").Text()
	cover, _ := doc.Find(".box_profile img").Attr("src")
	artist.Cover = extractArtistCover(cover)
	return &artist, nil
}
