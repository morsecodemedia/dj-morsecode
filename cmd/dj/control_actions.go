package main

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/morsecodemedia/dj-morsecode/internal/controls"
)

func (m model) applyControlAction(
	action controls.Action,
) (tea.Model, tea.Cmd) {

	switch action {

	case controls.ActionQuit:

		return m, tea.Quit

	case controls.ActionCancel:

		if m.CommandMode ==
			controls.ModeVolume {

			m = m.enterCommandMode(
				controls.ModeControls,
			)

			return m, nil
		}

		m = m.leaveCommandMode()

		return m, nil

	case controls.ActionEnterControls:

		m = m.enterCommandMode(
			controls.ModeControls,
		)

		return m, nil

	case controls.ActionEnterTune:

		m = m.enterCommandMode(
			controls.ModeTune,
		)

		return m, nil

	case controls.ActionEnterEnhancements:

		m = m.enterCommandMode(
			controls.ModeEnhancements,
		)

		return m, nil

	case controls.ActionEnterInfo:

		m = m.enterCommandMode(
			controls.ModeInfo,
		)

		return m, nil

	case controls.ActionEnterVolume:

		m = m.enterCommandMode(
			controls.ModeVolume,
		)

		return m, nil

	case controls.ActionOpenStations:

		m = m.leaveCommandMode()
		m.StationPickerOpen = true
		m.StationIndex = 0

		return m, nil

	case controls.ActionOpenGenres:

		m = m.leaveCommandMode()
		m.GenrePickerOpen = true
		m.GenreIndex = 0

		return m, nil

	case controls.ActionOpenMoods:

		m = m.leaveCommandMode()
		m.MoodPickerOpen = true
		m.MoodIndex = 0

		return m, nil

	case controls.ActionOpenVibes:

		m = m.leaveCommandMode()
		m.VibePickerOpen = true
		m.VibeIndex = 0

		return m, nil

	case controls.ActionOpenObservations:

		m = m.leaveCommandMode()
		m.ObservationHistoryOpen = true

		return m, nil

	case controls.ActionOpenHistory:

		m = m.leaveCommandMode()
		m.StationHistoryOpen = true

		return m, nil

	case controls.ActionToggleLyrics:

		m.LyricsVisible =
			!m.LyricsVisible

		m = m.leaveCommandMode()

		return m, nil
	case controls.ActionBack:

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

	case controls.ActionNext:

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

	case controls.ActionTogglePause:

		paused := m.Player.Paused()

		if err := m.Player.SetPaused(
			!paused,
		); err != nil {

			return m, nil
		}

		return m, nil

	case controls.ActionToggleMute:

		muted := m.Player.Muted()

		if err := m.Player.SetMuted(
			!muted,
		); err != nil {

			return m, nil
		}

		return m, nil

	case controls.ActionVolumeUp:

		volume :=
			m.Player.Volume() + 5

		if volume > 100 {
			volume = 100
		}

		if err := m.Player.SetVolume(
			volume,
		); err != nil {

			return m, nil
		}

		return m, nil

	case controls.ActionVolumeDown:

		volume :=
			m.Player.Volume() - 5

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

	return m, nil

}
