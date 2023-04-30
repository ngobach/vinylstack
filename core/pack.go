package core

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/dhowden/tag"
)

func copyFile(src, dest string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}

	defer sf.Close()

	df, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer df.Close()

	_, err = io.Copy(df, sf)

	return err
}

func isFileExist(file string) bool {
	fi, _ := os.Stat(file)

	return fi != nil
}

type Pack struct {
	dest         string
	defaultCover string
	tracks       []Track
}

func (p *Pack) prepare() {
	os.MkdirAll(p.dest, 0777)
	p.tracks = []Track{}
}

func (p *Pack) SetDefaultCover(s string) {
	p.defaultCover = s
}

func (p *Pack) Clean() {
	os.RemoveAll(p.dest)
	p.prepare()
}

func (p *Pack) ImportMp3File(filePath string) {
	baseName := path.Base(filePath)
	destMp3File := path.Join(p.dest, baseName)
	destThumbnailFile := strings.Replace(destMp3File, ".mp3", ".jpg", 1)

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	metadata, err := tag.ReadFrom(file)
	if err != nil {
		panic(err)
	}

	if err := file.Close(); err != nil {
		panic(err)
	}

	if !isFileExist(destMp3File) {
		log.Println("Copied", destMp3File)
		if err = copyFile(filePath, destMp3File); err != nil {
			panic(err)
		}
	} else {
		log.Println("Skipped copying", destMp3File)
	}

	if !isFileExist(destThumbnailFile) {
		log.Println("Wrote", destThumbnailFile)
		if err = os.WriteFile(destThumbnailFile, metadata.Picture().Data, 0666); err != nil {
			panic(err)
		}
	} else {
		log.Println("Skipped writing", destThumbnailFile)
	}

	p.tracks = append(p.tracks, Track{
		Title:  metadata.Title(),
		Artist: metadata.Artist(),
		URL:    baseName,
		Cover:  strings.Replace(baseName, ".mp3", ".jpg", 1),
		Genres: []string{},
	})
}

func (p *Pack) ImportMp3Directory(dirPath string) {
	allFiles, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	mp3Files := []os.FileInfo{}
	for _, entry := range allFiles {
		if strings.HasSuffix(entry.Name(), ".mp3") {
			fi, err := entry.Info()
			if err != nil {
				panic(err)
			}

			mp3Files = append(mp3Files, fi)
		}
	}

	if len(mp3Files) == 0 {
		log.Fatalln("empty source directory:", dirPath)
	}

	log.Printf("Preparing copy %d file(s)", len(mp3Files))

	for _, fi := range mp3Files {
		p.ImportMp3File(path.Join(dirPath, fi.Name()))
	}
}

func (p *Pack) WriteIndex() {
	indexFilePath := path.Join(p.dest, "index.json")

	index := &Index{
		Tracks:       p.tracks,
		Genres:       []struct{}{},
		Artists:      []struct{}{},
		DefautlCover: p.defaultCover,
	}

	encoded, err := json.Marshal(index)
	if err != nil {
		panic(err)
	}

	os.WriteFile(indexFilePath, encoded, 0666)
}

func NewPack(dest string) *Pack {
	p := &Pack{
		dest: dest,
	}
	p.prepare()

	return p
}
