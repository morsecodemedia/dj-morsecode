package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

func (m model) updateContext(
	msg contextMsg,
) (tea.Model, tea.Cmd) {
	if msg.TrackID != m.LastTrack {
		return m, nil
	}

	if msg.Err != nil {
		return m, nil
	}

	if !msg.OK {
		return m, nil
	}

	m.TrackContext =
		metadata.MergeTrackContext(
			m.TrackContext,
			msg.Context,
		)

	return m, nil
}
