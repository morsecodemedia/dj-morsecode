package controls

type Action int

const (
	ActionNone Action = iota

	ActionQuit
	ActionCancel

	ActionEnterControls
	ActionEnterTune
	ActionEnterEnhancements
	ActionEnterInfo
	ActionEnterVolume

	ActionBack
	ActionNext
	ActionTogglePause
	ActionToggleMute
	ActionVolumeUp
	ActionVolumeDown

	ActionOpenStations
	ActionOpenGenres
	ActionOpenMoods
	ActionOpenVibes

	ActionToggleLyrics

	ActionOpenObservations
	ActionOpenHistory
)
