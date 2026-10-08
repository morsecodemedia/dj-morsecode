package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/morsecodemedia/dj-morsecode/internal/controls"
	"github.com/morsecodemedia/dj-morsecode/internal/observation"
	"github.com/morsecodemedia/dj-morsecode/internal/radio"
	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

func (m model) updateKey(
	msg tea.KeyMsg,
) (tea.Model, tea.Cmd) {

	key := msg.String()

	if m.MoodPickerOpen {

		switch key {

		case "esc":
			m.MoodPickerOpen = false
			return m, nil

		case "up", "k":

			if m.MoodIndex > 0 {
				m.MoodIndex--
			}

			return m, nil

		case "down", "j":

			if m.MoodIndex < len(radio.MoodFamilies)-1 {
				m.MoodIndex++
			}

			return m, nil

		case "enter":

			if len(radio.MoodFamilies) == 0 {
				return m, nil
			}

			family := radio.MoodFamilies[m.MoodIndex]
			intent := radio.MoodIntent(
				family,
			)
			m.FailedStationIDs = nil
			station, ok := radio.Choose(
				intent.Criteria,
				m.StationHistory,
				radio.ChooseOptions{
					RecentLimit: 3,
					Chooser:     radio.RandomCandidate,
					ExcludeIDs:  m.FailedStationIDs,
				},
			)
			if !ok {
				return m, nil
			}

			err := m.Player.Load(
				station.StreamURL,
			)
			if err != nil {
				return m, nil
			}

			m.SourceMedia =
				source.MediaItem{}

			m.observeStation(
				observation.StationTuneRequested,
				station.ID,
			)

			m.PendingStationID = station.ID
			m.PendingSince = time.Now()

			m.ActiveIntent = intent
			m.MoodPickerOpen = false

			return m, nil

		}

		return m, nil

	}

	if m.GenrePickerOpen {

		genres := radio.Genres()

		switch key {

		case "esc":
			m.GenrePickerOpen = false
			return m, nil

		case "up", "k":

			if m.GenreIndex > 0 {
				m.GenreIndex--
			}

			return m, nil

		case "down", "j":

			if m.GenreIndex < len(genres)-1 {
				m.GenreIndex++
			}

			return m, nil

		case "enter":

			if len(genres) == 0 {
				return m, nil
			}

			genre := genres[m.GenreIndex]
			intent := radio.GenreIntent(
				genre,
			)
			m.FailedStationIDs = nil
			station, ok := radio.Choose(
				intent.Criteria,
				m.StationHistory,
				radio.ChooseOptions{
					RecentLimit: 3,
					Chooser:     radio.RandomCandidate,
					ExcludeIDs:  m.FailedStationIDs,
				},
			)
			if !ok {
				return m, nil
			}

			err := m.Player.Load(
				station.StreamURL,
			)
			if err != nil {
				return m, nil
			}

			m.SourceMedia =
				source.MediaItem{}

			m.observeStation(
				observation.StationTuneRequested,
				station.ID,
			)

			m.PendingStationID = station.ID
			m.PendingSince = time.Now()

			m.ActiveIntent = intent
			m.GenrePickerOpen = false

			return m, nil

		}

		return m, nil

	}

	if m.VibePickerOpen {

		switch key {

		case "esc":
			m.VibePickerOpen = false
			return m, nil

		case "up", "k":

			if m.VibeIndex > 0 {
				m.VibeIndex--
			}

			return m, nil

		case "down", "j":

			if m.VibeIndex < len(radio.Presets)-1 {
				m.VibeIndex++
			}

			return m, nil

		case "enter":

			if len(radio.Presets) == 0 {
				return m, nil
			}

			preset := radio.Presets[m.VibeIndex]
			intent := radio.VibeIntent(
				preset,
			)
			m.FailedStationIDs = nil
			station, ok := radio.Choose(
				intent.Criteria,
				m.StationHistory,
				radio.ChooseOptions{
					RecentLimit: 3,
					Chooser:     radio.RandomCandidate,
					ExcludeIDs:  m.FailedStationIDs,
				},
			)
			if !ok {
				return m, nil
			}

			err := m.Player.Load(
				station.StreamURL,
			)
			if err != nil {
				return m, nil
			}

			m.SourceMedia =
				source.MediaItem{}

			m.observeStation(
				observation.StationTuneRequested,
				station.ID,
			)

			m.PendingStationID = station.ID
			m.PendingSince = time.Now()

			m.ActiveIntent = intent
			m.VibePickerOpen = false

			return m, nil

		}

		return m, nil

	}

	if m.ObservationHistoryOpen {

		switch key {

		case "esc", "o":
			m.ObservationHistoryOpen = false
			return m, nil

		}

		return m, nil

	}

	if m.StationHistoryOpen {

		switch key {

		case "esc", "h":
			m.StationHistoryOpen = false
			return m, nil

		}

		return m, nil

	}

	if m.StationPickerOpen {

		switch key {

		case "enter":

			if len(radio.Stations) == 0 {
				return m, nil
			}

			station := radio.Stations[m.StationIndex]

			err := m.Player.Load(
				station.StreamURL,
			)
			if err != nil {
				return m, nil
			}

			m.SourceMedia =
				source.MediaItem{}

			m.observeStation(
				observation.StationTuneRequested,
				station.ID,
			)

			m.ActiveIntent = radio.Intent{}

			m.PendingStationID =
				station.ID

			m.PendingSince =
				time.Now()

			m.StationPickerOpen = false
			m.FailedStationIDs = nil

			return m, nil

		case "esc":
			m.StationPickerOpen = false
			return m, nil

		case "up", "k":

			if m.StationIndex > 0 {
				m.StationIndex--
			}

			return m, nil

		case "down", "j":

			if m.StationIndex < len(radio.Stations)-1 {
				m.StationIndex++
			}

			return m, nil

		}

		return m, nil

	}

	if m.YouTubePickerOpen {

		items := youtubeItems()

		switch key {

		case "esc":

			m.YouTubePickerOpen = false
			return m, nil

		case "up", "k":

			if m.YouTubeIndex > 0 {
				m.YouTubeIndex--
			}

			return m, nil

		case "down", "j":

			if m.YouTubeIndex <
				len(items)-1 {

				m.YouTubeIndex++
			}

			return m, nil

		case "enter":

			if len(items) == 0 {
				return m, nil
			}

			item :=
				items[m.YouTubeIndex]

			// Playlists are catalog references,
			// not directly supported here yet.
			if item.Kind !=
				source.MediaVideo {

				return m, nil
			}

			if err := m.Player.Load(
				item.Ref.URI,
			); err != nil {

				return m, nil
			}

			m = m.clearTrackState()

			m.SourceMedia = item

			m.ActiveIntent =
				radio.Intent{}

			m.PendingStationID = ""
			m.PendingSince =
				time.Time{}

			m.FailedStationIDs =
				nil

			m.YouTubePickerOpen =
				false

			return m, nil
		}

		return m, nil
	}

	action := controls.Resolve(
		m.CommandMode,
		key,
	)

	return m.applyControlAction(
		action,
	)

}
