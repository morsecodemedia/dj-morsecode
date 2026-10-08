package observation

import (
	"errors"
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
	"github.com/morsecodemedia/dj-morsecode/internal/source"
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

func (s failingSink) RecordMedia(
	MediaObservation,
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

func TestServiceObserveMedia(
	t *testing.T,
) {

	sink := NewMemorySink()

	service := NewService(
		NewRecorder(),
		sink,
	)

	item := source.MediaItem{
		Kind: source.MediaVideo,

		Ref: source.ItemRef{
			Source: source.Source{
				Kind: source.KindYouTube,
			},

			ID:   "UnqR5XUcLew",
			URI:  "https://www.youtube.com/watch?v=UnqR5XUcLew",
			Name: "Beastie Boys - Intergalactic",
		},

		Artist: "Beastie Boys",
		Title:  "Intergalactic",
	}

	observed, recorded, err :=
		service.ObserveMedia(
			item,
		)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if !recorded {
		t.Fatal(
			"expected media observation",
		)
	}

	if observed.Item.Ref.ID !=
		item.Ref.ID {

		t.Errorf(
			"unexpected media ID %q",
			observed.Item.Ref.ID,
		)

	}

	if len(sink.Media()) != 1 {

		t.Fatalf(
			"expected one stored media observation, got %d",
			len(sink.Media()),
		)

	}

}

func TestServiceObserveMediaDeduplicates(
	t *testing.T,
) {

	sink := NewMemorySink()

	service := NewService(
		NewRecorder(),
		sink,
	)

	item := source.MediaItem{
		Kind: source.MediaVideo,

		Ref: source.ItemRef{
			Source: source.Source{
				Kind: source.KindYouTube,
			},

			ID:  "UnqR5XUcLew",
			URI: "https://www.youtube.com/watch?v=UnqR5XUcLew",
		},
	}

	if _, recorded, err :=
		service.ObserveMedia(
			item,
		); err != nil ||
		!recorded {

		t.Fatalf(
			"expected first observation, recorded=%t err=%v",
			recorded,
			err,
		)

	}

	if _, recorded, err :=
		service.ObserveMedia(
			item,
		); err != nil {

		t.Fatalf(
			"unexpected duplicate error: %v",
			err,
		)

	} else if recorded {

		t.Fatal(
			"expected duplicate observation to be ignored",
		)

	}

	if len(sink.Media()) != 1 {

		t.Fatalf(
			"expected one stored observation, got %d",
			len(sink.Media()),
		)

	}

}

func TestServiceObserveMediaPropagatesSinkFailure(
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

	item := source.MediaItem{
		Kind: source.MediaVideo,

		Ref: source.ItemRef{
			Source: source.Source{
				Kind: source.KindYouTube,
			},

			ID: "UnqR5XUcLew",
		},
	}

	_, recorded, err :=
		service.ObserveMedia(
			item,
		)

	if !recorded {
		t.Fatal(
			"expected observation to be recorded",
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

}
