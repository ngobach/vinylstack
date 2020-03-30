package csn

import "github.com/thanbaiks/vinylstack/downloader"

type ChiaSeNhacFactory struct{}

func (c ChiaSeNhacFactory) CommandHelp() string {
	return "User ID"
}

func (c ChiaSeNhacFactory) Name() string {
	return "Chia Se Nhac"
}

func (c ChiaSeNhacFactory) CommandPrefix() string {
	return "csn"
}

func (c ChiaSeNhacFactory) Create(option string) downloader.Downloader {
	return &ChiaSeNhac{
		UserID: option,
	}
}
