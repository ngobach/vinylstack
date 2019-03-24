package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	targetDir string = "dist"
)

type exporter struct{}

/**
 * Remove
 */
func (exporter) prepare() error {
	fi, err := os.Stat(targetDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		// No error, should remove all
		err := os.RemoveAll(targetDir)
		if err != nil {
			return err
		}
	}
	err = os.Mkdir(targetDir, os.ModeDir)
	if err != nil {
		return err
	}
}

func export(file string, items []song) error {
	encoded, err := json.Marshal(items)
	if err != nil {
		return err
	}
	ioutil.WriteFile(file, encoded, 0644)
	return nil
}
