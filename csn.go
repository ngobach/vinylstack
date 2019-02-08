package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func parse(albumId string) ([]song, error) {
	resp, err := http.Get("https://beta.chiasenhac.vn/playlist/" + albumId + ".html")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	return []song{}, nil
}
