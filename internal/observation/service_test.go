package observation

import (
	"errors"
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

func TestServiceRecordsPlaybackTransition(
	t *testing.T,
) {

	sink := NewMemorySink()

	service := NewService(
		NewRecorder(),
		sink,
	)

	item := metadata.PlaybackItem{
		Type:     metadata.PlaybackTrack,
		Artist:   "Hozier",
		Title:    "Too Sweet",
		RawTitle: "Hozier - Too Sweet",
	}

	observation, recorded, err :=
		service.ObservePlayback(
			item,
			"station-1",
			"track-1",
		)

	if err != nil {
		t.Fatalf(
			"ObservePlayback returned error: %v",
			err,
		)
	}

	if !recorded {
		t.Fatal(
			"expected playback transition",
		)
	}

	if observation.TrackID != "track-1" {
		t.Errorf(
			"expected track ID %q, got %q",
			"track-1",
			observation.TrackID,
		)
	}

	if len(sink.Playback()) != 1 {
		t.Fatalf(
			"expected one stored observation, got %d",
			len(sink.Playback()),
		)
	}

}

func TestServiceSuppressesRepeatedPlayback(
	t *testing.T,
) {

	sink := NewMemorySink()

	service := NewService(
		NewRecorder(),
		sink,
	)

	item := metadata.PlaybackItem{
		Type:     metadata.PlaybackTrack,
		Artist:   "Hozier",
		Title:    "Too Sweet",
		RawTitle: "Hozier - Too Sweet",
	}

	_, recorded, err :=
		service.ObservePlayback(
			item,
			"station-1",
			"track-1",
		)

	if err != nil || !recorded {
		t.Fatal(
			"expected first observation",
		)
	}

	_, recorded, err =
		service.ObservePlayback(
			item,
			"station-1",
			"track-1",
		)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if recorded {
		t.Fatal(
			"expected repeated observation to be suppressed",
		)
	}

	if len(sink.Playback()) != 1 {
		t.Fatalf(
			"expected one stored observation, got %d",
			len(sink.Playback()),
		)
	}

}

type failingSink struct {
	err error
}

func (s failingSink) RecordStation(
	StationObservation,
) error {

	return s.err

}

func (s failingSink) RecordPlayback(
	PlaybackObservation,
) error {

	return s.err

}

func TestServiceDoesNotReplayTransitionAfterSinkFailure(
	t *testing.T,
) {

	expectedErr := errors.New(
		"sink failed",
	)

	service := NewService(
		NewRecorder(),
		failingSink{
			err: expectedErr,
		},
	)

	item := metadata.PlaybackItem{
		Type:     metadata.PlaybackTrack,
		Artist:   "Hozier",
		Title:    "Too Sweet",
		RawTitle: "Hozier - Too Sweet",
	}

	_, recorded, err :=
		service.ObservePlayback(
			item,
			"station-1",
			"track-1",
		)

	if !recorded {
		t.Fatal(
			"expected transition to be observed",
		)
	}

	if !errors.Is(
		err,
		expectedErr,
	) {

		t.Fatalf(
			"expected sink error, got %v",
			err,
		)

	}

	_, recorded, err =
		service.ObservePlayback(
			item,
			"station-1",
			"track-1",
		)

	if err != nil {
		t.Fatalf(
			"unexpected second error: %v",
			err,
		)
	}

	if recorded {
		t.Fatal(
			"expected transition not to be replayed",
		)
	}

}
