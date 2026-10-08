package observation

import (
	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

type Service struct {
	recorder *Recorder
	sink     Sink
}

func NewService(
	recorder *Recorder,
	sink Sink,
) *Service {

	return &Service{
		recorder: recorder,
		sink:     sink,
	}

}

func (s *Service) ObserveStation(
	kind StationKind,
	stationID string,
) (
	StationObservation,
	bool,
	error,
) {

	if s == nil ||
		s.recorder == nil {

		return StationObservation{},
			false,
			nil
	}

	observation, recorded :=
		s.recorder.ObserveStation(
			kind,
			stationID,
		)

	if !recorded {
		return StationObservation{},
			false,
			nil
	}

	if s.sink == nil {
		return observation,
			true,
			nil
	}

	err := s.sink.RecordStation(
		observation,
	)

	return observation,
		true,
		err

}

func (s *Service) ObservePlayback(
	item metadata.PlaybackItem,
	stationID string,
	trackID string,
) (
	PlaybackObservation,
	bool,
	error,
) {

	if s == nil ||
		s.recorder == nil {

		return PlaybackObservation{},
			false,
			nil
	}

	observation, recorded :=
		s.recorder.ObservePlayback(
			item,
			stationID,
			trackID,
		)

	if !recorded {
		return PlaybackObservation{},
			false,
			nil
	}

	if s.sink == nil {
		return observation,
			true,
			nil
	}

	err := s.sink.RecordPlayback(
		observation,
	)

	return observation,
		true,
		err

}

func (s *Service) ObserveMedia(
	item source.MediaItem,
) (
	MediaObservation,
	bool,
	error,
) {

	if s == nil ||
		s.recorder == nil {

		return MediaObservation{},
			false,
			nil
	}

	observation, recorded :=
		s.recorder.ObserveMedia(
			item,
		)

	if !recorded {
		return MediaObservation{},
			false,
			nil
	}

	if s.sink == nil {
		return observation,
			true,
			nil
	}

	err := s.sink.RecordMedia(
		observation,
	)

	return observation,
		true,
		err

}
