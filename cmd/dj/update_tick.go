package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/morsecodemedia/dj-morsecode/internal/library"
	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
	"github.com/morsecodemedia/dj-morsecode/internal/music"
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
	album := resolved.Album
	duration := resolved.Duration
	trackID := resolved.TrackID

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

		if trackID != m.LastTrack {

			m.Song = music.Song{
				Title:    track.Title,
				Artist:   track.Artist,
				Album:    album,
				Duration: duration,
			}

			m.CurrentCue = 0
			m.EnrichmentMatch =
				metadata.EnrichmentMatch{}
			m.TrackContext =
				metadata.TrackContext{}

			enrichmentDuration := time.Duration(0)

			if snapshot.IsNetwork &&
				m.PlaybackItem.Duration > 0 {

				enrichmentDuration =
					m.PlaybackItem.Duration

			} else if !snapshot.IsNetwork {

				enrichmentDuration = duration

			}

			enrichmentItem := metadata.PlaybackItem{
				Type:     metadata.PlaybackTrack,
				Artist:   track.Artist,
				Title:    track.Title,
				Duration: enrichmentDuration,
			}

			enrichmentCmd := enrichTrack(
				m.EnrichmentService,
				trackID,
				enrichmentItem,
			)

			song, ok := library.Load(
				track.Artist,
				track.Title,
			)

			if ok {
				m.Song.Lyrics = song.Lyrics
				m.Song.Timeline = song.Timeline
				m.LyricsState = lyricsLocal
				m.LastTrack = trackID

				return m, tea.Batch(
					tick(),
					enrichmentCmd,
				)

			}

			m.LastTrack = trackID
			m.LyricsState = lyricsSearching
			m.LyricsLookupDuration = duration

			return m, tea.Batch(
				tick(),
				loadLRCLIBSong(
					trackID,
					track.Artist,
					track.Title,
					duration,
				),
				enrichmentCmd,
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
