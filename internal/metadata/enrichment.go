package metadata

import "time"

type Identifier struct {
	Scheme string
	Value  string
}

const (
	IdentifierISRC        = "isrc"
	IdentifierMusicBrainz = "musicbrainz-recording"
)

type CanonicalTrack struct {
	Artist string
	Title  string

	Identifiers []Identifier
}

type EnrichmentCandidate struct {
	Artist string
	Title  string

	Duration time.Duration

	Variant string

	Identifiers []Identifier

	Provider      string
	ProviderScore float64
}

type EnrichmentMatch struct {
	Track CanonicalTrack

	Provider   string
	Confidence float64
}
