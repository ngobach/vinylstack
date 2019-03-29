package main

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
)

type exporter struct {
	target string
}

func (e *exporter) prepare() error {
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

func getFileName(u string) string {
	realurl, _ := url.PathUnescape(u)
	tmp := strings.Split(realurl, "/")
	return tmp[len(tmp)-1]
}

func (e *exporter) export(songs []song) error {
	const workerCount = 4
	ch := make(chan song)
	done := make(chan song)

	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				song, more := <-ch
				if more {
					filename := getFileName(song.URL)
					song.URL = filename
					_, err := os.Stat(path.Join(e.target, filename))
					if os.IsNotExist(err) {
						fmt.Println("Downloading", filename)
						resp, err := http.Get(song.URL)
						if err != nil {
							panic(err)
						}
						file, err := os.Create(path.Join(e.target, filename))
						if err != nil {
							panic(err)
						}
						io.Copy(file, resp.Body)
						file.Close()
					} else {
						fmt.Println("Skip downloading", filename)
					}
					done <- song
				} else {
					break
				}
			}
		}()
	}

	go func() {
		for _, song := range songs {
			ch <- song
		}
	}()
	newList := make([]song, 0)
	for range songs {
		newList = append(newList, <-done)
	}
	encoded, err := json.Marshal(newList)
	if err != nil {
		return err
	}
	fmt.Println("Target file", path.Join(e.target, "index.json"))
	return ioutil.WriteFile(path.Join(e.target, "index.json"), encoded, os.ModePerm)
}
