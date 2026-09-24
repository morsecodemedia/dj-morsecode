package musicbrainz

import "time"

type Recording struct {
	ID string

	Artist string
	Title  string

	Score    int
	Duration time.Duration

	ISRCs []string

	Disambiguation string
	Releases       []Release
}

type Release struct {
	ID      string
	Title   string
	Date    string
	Country string
}
