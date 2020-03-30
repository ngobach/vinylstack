package downloader

import "github.com/thanbaiks/vinylstack/core"

type Downloader interface {
	Download(store *core.Store) error
}
