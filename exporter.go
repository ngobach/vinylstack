package main

import (
	"encoding/json"
	"io/ioutil"
)

func export(file string, items []song) error {
	encoded, err := json.Marshal(items)
	if err != nil {
		return err
	}
	ioutil.WriteFile(file, encoded, 0644)
	return nil
}
