package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/morsecodemedia/dj-morsecode/internal/player"
	"github.com/morsecodemedia/dj-morsecode/internal/radio"
)

func (m model) updateTick(
	msg tickMsg,
) (tea.Model, tea.Cmd) {

	_ = msg

	snapshot := snapshotPlayback(
		m.Player,
	)

	position := m.playbackPosition()
	path := snapshot.Path
	if snapshot.Duration > 0 {
		m.Song.Duration =
			snapshot.Duration
	}

	var station *radio.Station
	var stationFound bool

	m, station, stationFound =
		m.reconcileStationIdentity(
			path,
		)

	var stop bool
	var resolved playbackTrack

	m, resolved =
		m.resolvePlaybackTrack(
			snapshot,
		)

	if snapshot.IsNetwork &&
		resolved.Observed {

		m.observePlayback(
			m.PlaybackItem,
			m.CurrentStationID,
			resolved.TrackID,
		)

	}

	m, stop =
		m.reconcileStationLifecycle(
			station,
			stationFound,
			snapshot.IsIdle,
			time.Now(),
		)

	if stop {
		return m, tick()
	}

	track := resolved.Track

	if snapshot.IsNetwork &&
		resolved.Observed &&
		!m.PlaybackItem.IsTrack() {

		m = m.clearTrackState()
		return m, tick()
	}

	if track.RawTitle == "" {
		return m, tick()
	}

	if track.Valid {

		if snapshot.IsNetwork {

			m.NowPlaying = m.PlaybackItem.DisplayTitle()

		} else {

			m.NowPlaying = track.RawTitle

		}

		if resolved.TrackID != m.LastTrack {
			return m.startTrack(
				resolved,
				snapshot.IsNetwork,
			)
		}

	}

	m.CurrentCue = player.CurrentCue(
		m.Song.Timeline,
		position,
	)

	m.OnAir = !m.OnAir
	return m, tick()
}
