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
	"sync"

	"github.com/thanbaiks/vinylstack/core"
)

const concurrentDownload = 4

func filenameFromURL(u string) string {
	realURL, _ := url.PathUnescape(u)
	tmp := strings.Split(realURL, "/")
	return tmp[len(tmp)-1]
}

type Exporter struct {
	target string
}

func NewExporter(target string) Exporter {
	return Exporter{target}
}

func (e *Exporter) Prepare() error {
	stat, err := os.Stat(e.target)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if stat != nil && !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", e.target)
	}
	if stat == nil {
		err = os.Mkdir(e.target, os.ModePerm)
	}
	if err != nil {
		return err
	}
	return nil
}

func (e *Exporter) DownloadAndExport(store *core.Store) error {
	count := len(store.Tracks)
	inputs := make(chan core.ID, count)
	done := make(chan core.Track)
	for _, track := range store.Tracks {
		inputs <- track.ID
	}
	close(inputs)
	mtx := sync.Mutex{}
	for i := 0; i < concurrentDownload; i++ {
		go func() {
			for id := range inputs {
				mtx.Lock()
				track := store.Tracks[id]
				mtx.Unlock()
				filename := filenameFromURL(track.URL)
				_, err := os.Stat(path.Join(e.target, filename))
				if os.IsNotExist(err) {
					resp, err := http.Get(track.URL)
					if err != nil {
						panic(err)
					}
					file, err := os.Create(path.Join(e.target, filename))
					if err != nil {
						panic(err)
					}
					_, err = io.Copy(file, resp.Body)
					if err != nil {
						panic(err)
					}
					err = file.Close()
					if err != nil {
						panic(err)
					}
					fmt.Println("++", "Downloaded", filename)
				} else {
					fmt.Println("==", "Skipped", filename)
				}
				track.URL = filename
				mtx.Lock()
				store.Tracks[id] = track
				mtx.Unlock()
				done <- track
			}
		}()
	}

	for i := 0; i < count; i++ {
		<-done
	}

	encoded, err := json.Marshal(store)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(e.target, "index.json"), encoded, os.ModePerm)
}
