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
	duration := snapshot.Duration

	if duration > 0 {
		m.Song.Duration = duration
	}

	rawTitle := snapshot.RawTitle
	artist := snapshot.Artist
	title := snapshot.Title
	album := snapshot.Album
	trackID := snapshot.TrackID
	path := snapshot.Path
	isNetwork := snapshot.IsNetwork

	var station *radio.Station
	var stationFound bool

	m, station, stationFound =
		m.reconcileStationIdentity(
			path,
		)

	track := metadata.Resolve(
		rawTitle,
	)

	observedPlaybackItem := false

	if isNetwork {

		observedItem := metadata.Normalize(
			snapshot.Metadata,
		)

		if observedItem.Observed() {

			observedPlaybackItem = true
			m.PlaybackItem = observedItem

		}

		if m.PlaybackItem.IsTrack() {

			track.RawTitle =
				m.PlaybackItem.RawTitle

			track.Artist =
				m.PlaybackItem.Artist

			track.Title =
				m.PlaybackItem.Title

			track.Valid = true

			album = m.PlaybackItem.Album

			if m.PlaybackItem.Duration > 0 {
				duration =
					m.PlaybackItem.Duration
			}

			trackID = path +
				"\x00" +
				m.PlaybackItem.Artist +
				"\x00" +
				m.PlaybackItem.Title

		} else {

			track.Valid = false

		}

	} else {

		m.PlaybackItem =
			metadata.PlaybackItem{}

	}

	if isNetwork &&
		observedPlaybackItem {

		m.observePlayback(
			m.PlaybackItem,
			m.CurrentStationID,
			trackID,
		)

	}

	var stop bool

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

	if !isNetwork {
		if artist != "" {
			track.Artist = artist
		}

		if title != "" {
			track.Title = title
		}
	}

	if isNetwork &&
		observedPlaybackItem &&
		!m.PlaybackItem.IsTrack() {

		m = m.clearTrackState()
		return m, tick()
	}

	if track.RawTitle == "" {
		return m, tick()
	}

	if track.Valid {

		if isNetwork {

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

			if isNetwork &&
				m.PlaybackItem.Duration > 0 {

				enrichmentDuration =
					m.PlaybackItem.Duration

			} else if !isNetwork {

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
