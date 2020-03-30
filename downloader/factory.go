package downloader

type Factory interface {
	Name() string
	CommandHelp() string
	CommandPrefix() string
	Create(string) Downloader
}

var factories []Factory = nil

func RegisterFactory(factory Factory) {
	factories = append(factories, factory)
}

func Factories() []Factory {
	return factories
}
