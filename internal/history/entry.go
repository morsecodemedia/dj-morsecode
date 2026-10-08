package history

import (
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

type Kind string

const (
	KindStation Kind = "station"
	KindMedia   Kind = "media"
)

type Entry struct {
	Kind Kind

	StationID string
	Media     source.MediaItem

	PlayedAt time.Time
}
