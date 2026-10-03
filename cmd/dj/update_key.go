package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/morsecodemedia/dj-morsecode/internal/observation"
	"github.com/morsecodemedia/dj-morsecode/internal/radio"
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

			m.observeStation(
				observation.StationTuneRequested,
				station.ID,
			)

			m.ActiveIntent = radio.Intent{}
			m.PendingStationID = ""
			m.PendingSince = time.Time{}
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

	if m.CommandMode !=
		commandModeNone {

		switch key {

		case "esc":

			if m.CommandMode ==
				commandModeVolume {

				m = m.enterCommandMode(
					commandModeControls,
				)

				return m, nil
			}

			m = m.leaveCommandMode()

			return m, nil

		}

		switch m.CommandMode {

		case commandModeTune:

			switch key {

			case "s":
				m = m.leaveCommandMode()
				m.StationPickerOpen = true
				m.StationIndex = 0

				return m, nil

			case "g":
				m = m.leaveCommandMode()
				m.GenrePickerOpen = true
				m.GenreIndex = 0

				return m, nil

			case "m":
				m = m.leaveCommandMode()
				m.MoodPickerOpen = true
				m.MoodIndex = 0

				return m, nil

			case "v":
				m = m.leaveCommandMode()
				m.VibePickerOpen = true
				m.VibeIndex = 0

				return m, nil

			}

		case commandModeInfo:

			switch key {

			case "o":
				m = m.leaveCommandMode()
				m.ObservationHistoryOpen = true

				return m, nil

			case "h":
				m = m.leaveCommandMode()
				m.StationHistoryOpen = true

				return m, nil

			}

		case commandModeEnhancements:

			switch key {

			case "l":
				m.LyricsVisible =
					!m.LyricsVisible

				m = m.leaveCommandMode()

				return m, nil

			}

		case commandModeControls:

			switch key {

			case "b":

				if m.StationHistory.CanBack() {

					tune, ok :=
						m.StationHistory.Back()

					if !ok {
						return m, nil
					}

					m = m.tuneHistory(
						tune,
					)

					m = m.leaveCommandMode()

					return m, nil
				}

				m = m.previousStation()
				m = m.leaveCommandMode()

				return m, nil

			case "n":

				if m.StationHistory.CanForward() {

					tune, ok :=
						m.StationHistory.Forward()

					if !ok {
						return m, nil
					}

					m = m.tuneHistory(
						tune,
					)

					m = m.leaveCommandMode()

					return m, nil
				}

				m = m.nextStation()
				m = m.leaveCommandMode()

				return m, nil

			case "p":

				paused := m.Player.Paused()

				if err := m.Player.SetPaused(
					!paused,
				); err != nil {

					return m, nil
				}

				return m, nil

			case "m":

				muted := m.Player.Muted()

				if err := m.Player.SetMuted(
					!muted,
				); err != nil {

					return m, nil
				}

				return m, nil

			case "v":

				m = m.enterCommandMode(
					commandModeVolume,
				)

				return m, nil

			}
		case commandModeVolume:

			switch key {

			case "up":

				volume := m.Player.Volume() + 5

				if volume > 100 {
					volume = 100
				}

				if err := m.Player.SetVolume(
					volume,
				); err != nil {

					return m, nil
				}

				return m, nil

			case "down":

				volume := m.Player.Volume() - 5

				if volume < 0 {
					volume = 0
				}

				if err := m.Player.SetVolume(
					volume,
				); err != nil {

					return m, nil
				}

				return m, nil

			}
		}

		return m, nil
	}

	switch key {

	case "ctrl+c", "q":
		return m, tea.Quit

	case "c":
		m = m.enterCommandMode(
			commandModeControls,
		)
		return m, nil

	case "t":
		m = m.enterCommandMode(
			commandModeTune,
		)
		return m, nil

	case "e":
		m = m.enterCommandMode(
			commandModeEnhancements,
		)
		return m, nil

	case "i":
		m = m.enterCommandMode(
			commandModeInfo,
		)
		return m, nil

	}

	return m, nil

}
