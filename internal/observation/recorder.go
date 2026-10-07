package observation

import (
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

type Recorder struct {
	now func() time.Time

	lastStation  *StationObservation
	lastPlayback *PlaybackObservation
	lastMedia    *MediaObservation
}

func NewRecorder() *Recorder {

	return &Recorder{
		now: time.Now,
	}

}

func newRecorderWithClock(
	now func() time.Time,
) *Recorder {

	return &Recorder{
		now: now,
	}

}

func (r *Recorder) ObserveStation(
	kind StationKind,
	stationID string,
) (StationObservation, bool) {

	if stationID == "" {
		return StationObservation{}, false
	}

	observation := StationObservation{
		Kind:       kind,
		StationID:  stationID,
		ObservedAt: r.now(),
	}

	if r.lastStation != nil &&
		r.lastStation.Kind == observation.Kind &&
		r.lastStation.StationID == observation.StationID {

		return StationObservation{}, false
	}

	r.lastStation = &observation

	return observation, true

}

func (r *Recorder) ObservePlayback(
	item metadata.PlaybackItem,
	stationID string,
	trackID string,
) (PlaybackObservation, bool) {

	if !item.Observed() {
		return PlaybackObservation{}, false
	}

	observation := PlaybackObservation{
		Item:       item,
		StationID:  stationID,
		TrackID:    trackID,
		ObservedAt: r.now(),
	}

	if r.lastPlayback != nil &&
		samePlaybackObservation(
			*r.lastPlayback,
			observation,
		) {

		return PlaybackObservation{}, false
	}

	r.lastPlayback = &observation

	return observation, true

}

func samePlaybackObservation(
	left PlaybackObservation,
	right PlaybackObservation,
) bool {

	return left.Item.Type ==
		right.Item.Type &&
		left.Item.Artist ==
			right.Item.Artist &&
		left.Item.Title ==
			right.Item.Title &&
		left.Item.Album ==
			right.Item.Album &&
		left.Item.RawTitle ==
			right.Item.RawTitle &&
		left.StationID ==
			right.StationID &&
		left.TrackID ==
			right.TrackID

}

func (r *Recorder) ObserveMedia(
	item source.MediaItem,
) (MediaObservation, bool) {

	if !item.Valid() {
		return MediaObservation{}, false
	}

	observation := MediaObservation{
		Item:       item,
		ObservedAt: r.now(),
	}

	if r.lastMedia != nil &&
		sameMediaObservation(
			*r.lastMedia,
			observation,
		) {

		return MediaObservation{}, false
	}

	r.lastMedia = &observation

	return observation, true

}

func sameMediaObservation(
	left MediaObservation,
	right MediaObservation,
) bool {

	return left.Item.Kind ==
		right.Item.Kind &&
		left.Item.Ref.Source.Kind ==
			right.Item.Ref.Source.Kind &&
		left.Item.Ref.ID ==
			right.Item.Ref.ID &&
		left.Item.Ref.URI ==
			right.Item.Ref.URI

}
