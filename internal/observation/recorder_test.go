package observation

import (
	"testing"
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

func TestRecorderSuppressesRepeatedPlaybackObservation(
	t *testing.T,
) {

	now := time.Date(
		2026,
		time.September,
		27,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	recorder := newRecorderWithClock(
		func() time.Time {
			return now
		},
	)

	item := metadata.PlaybackItem{
		Type:     metadata.PlaybackUnknown,
		RawTitle: ` - text="Spot Block End"`,
	}

	first, ok := recorder.ObservePlayback(
		item,
		"z100",
		"",
	)

	if !ok {
		t.Fatal(
			"expected first observation",
		)
	}

	if first.Item.RawTitle != item.RawTitle {
		t.Errorf(
			"expected raw title %q, got %q",
			item.RawTitle,
			first.Item.RawTitle,
		)
	}

	now = now.Add(
		15 * time.Minute,
	)

	_, ok = recorder.ObservePlayback(
		item,
		"z100",
		"",
	)

	if ok {
		t.Fatal(
			"expected repeated observation to be suppressed",
		)
	}

}
func TestRecorderRecordsPlaybackTransition(
	t *testing.T,
) {

	recorder := NewRecorder()

	first := metadata.PlaybackItem{
		Type:     metadata.PlaybackTrack,
		Artist:   "Hozier",
		Title:    "Too Sweet",
		RawTitle: "Hozier - Too Sweet",
	}

	second := metadata.PlaybackItem{
		Type:     metadata.PlaybackTrack,
		Artist:   "Europe",
		Title:    "The Final Countdown",
		RawTitle: "Europe - The Final Countdown",
	}

	if _, ok := recorder.ObservePlayback(
		first,
		"station-1",
		"track-1",
	); !ok {

		t.Fatal(
			"expected first observation",
		)

	}

	observation, ok :=
		recorder.ObservePlayback(
			second,
			"station-1",
			"track-2",
		)

	if !ok {
		t.Fatal(
			"expected playback transition",
		)
	}

	if observation.Item.Title !=
		"The Final Countdown" {

		t.Errorf(
			"unexpected title %q",
			observation.Item.Title,
		)

	}

}
func TestRecorderRecordsStationTransitions(
	t *testing.T,
) {

	recorder := NewRecorder()

	if _, ok := recorder.ObserveStation(
		StationTuneRequested,
		"z100",
	); !ok {

		t.Fatal(
			"expected tune request",
		)

	}

	if _, ok := recorder.ObserveStation(
		StationTuneRequested,
		"z100",
	); ok {

		t.Fatal(
			"expected duplicate tune request to be suppressed",
		)

	}

	if _, ok := recorder.ObserveStation(
		StationTuneConfirmed,
		"z100",
	); !ok {

		t.Fatal(
			"expected tune confirmation",
		)

	}

}

func TestRecorderRejectsUnobservedPlayback(
	t *testing.T,
) {

	recorder := NewRecorder()

	_, ok := recorder.ObservePlayback(
		metadata.PlaybackItem{},
		"station-1",
		"",
	)

	if ok {
		t.Fatal(
			"expected unobserved playback item to be rejected",
		)
	}

}
