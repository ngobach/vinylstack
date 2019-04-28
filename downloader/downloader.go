package downloader

import "github.com/thanbaiks/vinylstack/core"

// Downloader provide the base interface for other downloaders
type Downloader interface {
	Info() string
	Download() ([]core.Playlist, error)
}
