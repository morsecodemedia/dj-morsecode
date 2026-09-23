package metadata

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

type EnrichmentMatch struct {
	Track CanonicalTrack

	Provider   string
	Confidence float64
}
