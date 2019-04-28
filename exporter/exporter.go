package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/thanbaiks/vinylstack/core"
)

const concurrentDownload = 16

func filenameFromURL(u string) string {
	realurl, _ := url.PathUnescape(u)
	tmp := strings.Split(realurl, "/")
	return tmp[len(tmp)-1]
}

// Exporter download related files to local disk
type Exporter struct {
	// Target directory. Must be related to working directory
	Target string
}

// Prepare target directory to be exportable
func (e *Exporter) Prepare() error {
	stat, err := os.Stat(e.Target)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if stat != nil && !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", e.Target)
	}
	if stat == nil {
		err = os.Mkdir(e.Target, os.ModePerm)
	}
	if err != nil {
		return err
	}
	return nil
}

// DownloadAndExport download all related media to local disk.
// It also export index.json file
func (e *Exporter) DownloadAndExport(songs []core.Song) error {
	count := len(songs)
	inputs := make(chan core.Song, count)
	done := make(chan core.Song)
	for _, song := range songs {
		inputs <- song
	}
	close(inputs)
	songs = make([]core.Song, 0, count)
	for i := 0; i < concurrentDownload; i++ {
		go func() {
			for {
				song, more := <-inputs
				if !more {
					return
				}
				// fmt.Println("::", "Working on", song.Title)
				filename := filenameFromURL(song.URL)
				_, err := os.Stat(path.Join(e.Target, filename))
				if os.IsNotExist(err) {
					resp, err := http.Get(song.URL)
					if err != nil {
						panic(err)
					}
					file, err := os.Create(path.Join(e.Target, filename))
					if err != nil {
						panic(err)
					}
					io.Copy(file, resp.Body)
					file.Close()
					fmt.Println("++", "Downloaded", filename)
				} else {
					fmt.Println("==", "Skipped", filename)
				}
				song.URL = filename
				done <- song
			}
		}()
	}

	for i := 0; i < count; i++ {
		songs = append(songs, <-done)
	}

	encoded, err := json.Marshal(songs)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(e.Target, "index.json"), encoded, os.ModePerm)
}
