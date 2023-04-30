package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func main() {
	linkFileContent, err := ioutil.ReadFile("./links.json")
	if err != nil {
		panic(err)
	}

	links := []string{}
	err = json.Unmarshal(linkFileContent, &links)

	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("./downloads", 0777)
	if err != nil {
		panic(err)
	}

	for _, v := range links {
		parts := strings.Split(v, "/")
		fileName, err := url.PathUnescape(strings.TrimSpace(parts[len(parts)-1]))
		if err != nil {
			panic(err)
		}
		downloadFile(v, path.Join("./downloads", fileName))
		log.Println("Downloaded", fileName)
	}
}

func downloadFile(url, saveAs string) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	file, err := os.Create(saveAs)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	io.Copy(file, resp.Body)
}
