package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/morsecodemedia/dj-morsecode/internal/library"
	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
	"github.com/morsecodemedia/dj-morsecode/internal/music"
)

func (m model) startTrack(
	resolved playbackTrack,
	isRadio bool,
) (tea.Model, tea.Cmd) {

	track := resolved.Track

	m.Song = music.Song{
		Title:    track.Title,
		Artist:   track.Artist,
		Album:    resolved.Album,
		Duration: resolved.Duration,
	}

	m.CurrentCue = 0

	m.EnrichmentMatch =
		metadata.EnrichmentMatch{}

	m.TrackContext =
		metadata.TrackContext{}

	enrichmentDuration :=
		time.Duration(0)

	if isRadio &&
		m.PlaybackItem.Duration > 0 {

		enrichmentDuration =
			m.PlaybackItem.Duration

	} else if !isRadio {

		enrichmentDuration =
			resolved.Duration

	}

	enrichmentItem :=
		metadata.PlaybackItem{
			Type: metadata.PlaybackTrack,

			Artist: track.Artist,
			Title:  track.Title,

			Duration: enrichmentDuration,
		}

	enrichmentCmd := enrichTrack(
		m.EnrichmentService,
		resolved.TrackID,
		enrichmentItem,
	)

	song, ok := library.Load(
		track.Artist,
		track.Title,
	)

	if ok {

		m.Song.Lyrics =
			song.Lyrics

		m.Song.Timeline =
			song.Timeline

		m.LyricsState =
			lyricsLocal

		m.LastTrack =
			resolved.TrackID

		return m, tea.Batch(
			tick(),
			enrichmentCmd,
		)
	}

	m.LastTrack =
		resolved.TrackID

	m.LyricsLookupDuration =
		resolved.Duration

	if !isRadio &&
		resolved.Duration <= 0 {

		m.LyricsState =
			lyricsSearching

		return m, tea.Batch(
			tick(),
			enrichmentCmd,
		)

	}

	m.LyricsState =
		lyricsSearching

	return m, tea.Batch(
		tick(),
		loadLRCLIBSong(
			resolved.TrackID,
			track.Artist,
			track.Title,
			resolved.Duration,
		),
		enrichmentCmd,
	)

}
