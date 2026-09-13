package player

import "strings"

func Progress(
	current float64,
	total float64,
) float64 {

	if total <= 0 {
		return 0
	}

	progress := current / total

	if progress < 0 {
		progress = 0
	}

	if progress > 1 {
		progress = 1
	}

	return progress

}

func ProgressBar(
	width int,
	progress float64,
) string {

	filled := int(progress * float64(width))

	if filled > width {
		filled = width
	}

	if filled < 0 {
		filled = 0
	}

	bar :=
		strings.Repeat("█", filled) +
			strings.Repeat("░", width-filled)

	return bar

}
