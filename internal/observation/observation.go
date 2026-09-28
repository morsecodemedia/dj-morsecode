package observation

import (
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

type StationKind string

const (
	StationTuneRequested StationKind = "tune-requested"
	StationTuneConfirmed StationKind = "tune-confirmed"
	StationTuneFailed    StationKind = "tune-failed"
)

type StationObservation struct {
	Kind StationKind

	StationID  string
	ObservedAt time.Time
}

type PlaybackObservation struct {
	Item metadata.PlaybackItem

	StationID string
	TrackID   string

	ObservedAt time.Time
}
