package player

import (
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/music"
)

func CurrentCue(
	timeline []music.Cue,
	elapsed time.Duration,
) int {

	if len(timeline) == 0 {
		return 0
	}

	for i := len(timeline) - 1; i >= 0; i-- {

		if elapsed >= timeline[i].Time {
			return i
		}

	}

	return 0

}
