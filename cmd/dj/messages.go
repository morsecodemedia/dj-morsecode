package main

import (
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
	"github.com/morsecodemedia/dj-morsecode/internal/music"
)

type tickMsg time.Time

type lrclibSongMsg struct {
	TrackID  string
	Song     music.Song
	Content  string
	Err      error
	Duration time.Duration
}

type enrichmentMsg struct {
	TrackID string

	Match  metadata.EnrichmentMatch
	Status metadata.MatchStatus
	Err    error
}

type contextMsg struct {
	TrackID string

	Context metadata.TrackContext
	OK      bool
	Err     error
}
