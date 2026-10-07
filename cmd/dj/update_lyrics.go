package main

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/morsecodemedia/dj-morsecode/internal/library"
	"github.com/morsecodemedia/dj-morsecode/internal/player"
)

func (m model) updateLyrics(
	msg lrclibSongMsg,
) (tea.Model, tea.Cmd) {

	if msg.TrackID != m.LastTrack {
		return m, nil
	}

	if msg.Duration !=
		m.LyricsLookupDuration {
		return m, nil
	}

	if msg.Err != nil {

		m.LyricsState =
			lyricsUnavailable

		return m, nil
	}

	_, err := library.Store(
		m.Song.Artist,
		m.Song.Title,
		msg.Content,
	)
	if err != nil {

		m.LyricsState =
			lyricsUnavailable

		return m, nil
	}

	m.Song.Lyrics = msg.Song.Lyrics
	m.Song.Timeline = msg.Song.Timeline
	m.LyricsState = lyricsRemote

	m.CurrentCue = player.CurrentCue(
		m.Song.Timeline,
		m.playbackPosition(),
	)

	return m, nil

}
