package observation

import (
	"testing"
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

func TestStationObservation(t *testing.T) {

	observedAt := time.Date(
		2026,
		time.September,
		27,
		10,
		0,
		0,
		0,
		time.UTC,
	)

	observation := StationObservation{
		Kind:       StationTuneRequested,
		StationID:  "z100",
		ObservedAt: observedAt,
	}

	if observation.Kind !=
		StationTuneRequested {

		t.Errorf(
			"expected tune requested, got %q",
			observation.Kind,
		)

	}

	if observation.StationID != "z100" {
		t.Errorf(
			"expected station %q, got %q",
			"z100",
			observation.StationID,
		)
	}

	if !observation.ObservedAt.Equal(
		observedAt,
	) {

		t.Errorf(
			"unexpected observation time %v",
			observation.ObservedAt,
		)

	}

}

func TestPlaybackObservationPreservesItem(
	t *testing.T,
) {

	item := metadata.PlaybackItem{
		Type:     metadata.PlaybackTrack,
		Artist:   "Hozier",
		Title:    "Too Sweet",
		RawTitle: "Hozier - Too Sweet",
	}

	observation := PlaybackObservation{
		Item:      item,
		StationID: "station-1",
		TrackID:   "track-1",
	}

	if observation.Item.Type !=
		metadata.PlaybackTrack {

		t.Errorf(
			"expected playback track, got %q",
			observation.Item.Type,
		)

	}

	if observation.Item.Artist != "Hozier" {
		t.Errorf(
			"expected artist %q, got %q",
			"Hozier",
			observation.Item.Artist,
		)
	}

	if observation.TrackID != "track-1" {
		t.Errorf(
			"expected track ID %q, got %q",
			"track-1",
			observation.TrackID,
		)
	}

}
