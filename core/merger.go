package core

type Merger struct {
	// Tracks by its ID
	tracks map[string]Track
}

func NewMerger() Merger {
	return Merger{
		tracks: nil,
	}
}

func (m *Merger) Add(tracks []Track) {
	// TODO:
}

func (m *Merger) Genres() []Genre {
	// TODO:
	return nil
}

func (m *Merger) Artists() []Artist {
	// TODO:
	return nil
}

func (m *Merger) Tracks() []Track {
	// TODO:
	return nil
}
