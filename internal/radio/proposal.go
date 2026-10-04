package radio

import "github.com/morsecodemedia/dj-morsecode/internal/source"

type StationProposal struct {
	Name      string
	StreamURL string
	Homepage  string

	Genre    string
	Tags     []string
	Moods    []string
	Contexts []string
	Energy   int

	Source source.ItemRef
}

func (p StationProposal) Valid() bool {

	return p.StreamURL != "" &&
		p.Source.Valid()

}
