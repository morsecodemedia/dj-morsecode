package observation

import (
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

type Recorder struct {
	now func() time.Time

	lastStation  *StationObservation
	lastPlayback *PlaybackObservation
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
