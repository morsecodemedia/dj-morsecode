package observation

type MemorySink struct {
	stations []StationObservation
	playback []PlaybackObservation
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
