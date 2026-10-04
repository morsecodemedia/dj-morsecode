package controls

func Resolve(
	mode Mode,
	key string,
) Action {

	if key == "esc" {

		if mode != ModeNone {
			return ActionCancel
		}

		return ActionNone
	}

	switch mode {

	case ModeNone:

		switch key {

		case "ctrl+c", "q":
			return ActionQuit

		case "c":
			return ActionEnterControls

		case "t":
			return ActionEnterTune

		case "e":
			return ActionEnterEnhancements

		case "i":
			return ActionEnterInfo

		}

	case ModeControls:

		switch key {

		case "b":
			return ActionBack

		case "n":
			return ActionNext

		case "p":
			return ActionTogglePause

		case "m":
			return ActionToggleMute

		case "v":
			return ActionEnterVolume

		}

	case ModeVolume:

		switch key {

		case "up":
			return ActionVolumeUp

		case "down":
			return ActionVolumeDown

		}

	case ModeTune:

		switch key {

		case "s":
			return ActionOpenStations

		case "g":
			return ActionOpenGenres

		case "m":
			return ActionOpenMoods

		case "v":
			return ActionOpenVibes

		}

	case ModeEnhancements:

		switch key {

		case "l":
			return ActionToggleLyrics

		}

	case ModeInfo:

		switch key {

		case "o":
			return ActionOpenObservations

		case "h":
			return ActionOpenHistory

		}

	}

	return ActionNone
}
