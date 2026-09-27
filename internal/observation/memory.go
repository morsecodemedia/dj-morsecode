package observation

import "time"

type MemorySink struct {
	stations []StationObservation
	playback []PlaybackObservation
}

func (s *MemorySink) PlaybackElapsed(
	trackID string,
	now time.Time,
) (time.Duration, bool) {

	if s == nil ||
		trackID == "" {

		return 0, false
	}

	for i := len(s.playback) - 1; i >= 0; i-- {

		observed := s.playback[i]

		if observed.TrackID != trackID {
			continue
		}

		if !observed.Item.IsTrack() {
			continue
		}

		elapsed := now.Sub(
			observed.ObservedAt,
		)

		if elapsed < 0 {
			elapsed = 0
		}

		return elapsed, true

	}

	return 0, false

}

func NewMemorySink() *MemorySink {

	return &MemorySink{}

}

func (s *MemorySink) RecordStation(
	observation StationObservation,
) error {

	s.stations = append(
		s.stations,
		observation,
	)

	return nil

}

func (s *MemorySink) RecordPlayback(
	observation PlaybackObservation,
) error {

	s.playback = append(
		s.playback,
		observation,
	)

	return nil

}

func (s *MemorySink) Stations() []StationObservation {

	return append(
		[]StationObservation(nil),
		s.stations...,
	)

}

func (s *MemorySink) Playback() []PlaybackObservation {

	return append(
		[]PlaybackObservation(nil),
		s.playback...,
	)

}
