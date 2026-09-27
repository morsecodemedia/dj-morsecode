package observation

import (
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

func TestMemorySinkRecordsStationObservations(
	t *testing.T,
) {

	sink := NewMemorySink()

	observation := StationObservation{
		Kind:      StationTuneConfirmed,
		StationID: "z100",
	}

	if err := sink.RecordStation(
		observation,
	); err != nil {

		t.Fatalf(
			"RecordStation returned error: %v",
			err,
		)

	}

	recorded := sink.Stations()

	if len(recorded) != 1 {
		t.Fatalf(
			"expected one station observation, got %d",
			len(recorded),
		)
	}

	if recorded[0].StationID != "z100" {
		t.Errorf(
			"expected station %q, got %q",
			"z100",
			recorded[0].StationID,
		)
	}

}
func TestMemorySinkRecordsPlaybackObservations(
	t *testing.T,
) {

	sink := NewMemorySink()

	observation := PlaybackObservation{
		Item: metadata.PlaybackItem{
			Type:     metadata.PlaybackTrack,
			Artist:   "Hozier",
			Title:    "Too Sweet",
			RawTitle: "Hozier - Too Sweet",
		},
		StationID: "station-1",
		TrackID:   "track-1",
	}

	if err := sink.RecordPlayback(
		observation,
	); err != nil {

		t.Fatalf(
			"RecordPlayback returned error: %v",
			err,
		)

	}

	recorded := sink.Playback()

	if len(recorded) != 1 {
		t.Fatalf(
			"expected one playback observation, got %d",
			len(recorded),
		)
	}

	if recorded[0].TrackID != "track-1" {
		t.Errorf(
			"expected track ID %q, got %q",
			"track-1",
			recorded[0].TrackID,
		)
	}

}
func TestMemorySinkReturnsPlaybackCopy(
	t *testing.T,
) {

	sink := NewMemorySink()

	_ = sink.RecordPlayback(
		PlaybackObservation{
			TrackID: "original",
		},
	)

	recorded := sink.Playback()

	recorded[0].TrackID = "mutated"

	again := sink.Playback()

	if again[0].TrackID != "original" {
		t.Errorf(
			"expected stored observation to remain unchanged, got %q",
			again[0].TrackID,
		)
	}

}
