package main

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

func (m model) updateEnrichment(
	msg enrichmentMsg,
) (tea.Model, tea.Cmd) {

	if msg.Err != nil {
		return m, nil
	}

	if msg.TrackID != m.LastTrack {
		return m, nil
	}

	if msg.Status != metadata.MatchAccepted {
		return m, nil
	}

	m.EnrichmentMatch = msg.Match

	contextCmd := loadTrackContexts(
		m.ContextService,
		msg.TrackID,
		msg.Match.Track,
	)

	if msg.Match.Duration <= 0 ||
		msg.Match.Duration == m.LyricsLookupDuration {

		return m, contextCmd
	}

	if m.LyricsState == lyricsLocal ||
		m.LyricsState == lyricsRemote {

		return m, contextCmd
	}

	m.LyricsLookupDuration =
		msg.Match.Duration

	m.LyricsState =
		lyricsSearching

	lyricsCmd := loadLRCLIBSong(
		msg.TrackID,
		msg.Match.Track.Artist,
		msg.Match.Track.Title,
		msg.Match.Duration,
	)

	return m, tea.Batch(
		contextCmd,
		lyricsCmd,
	)

}
