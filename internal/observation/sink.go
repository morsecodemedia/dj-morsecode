package observation

type Sink interface {
	RecordStation(
		StationObservation,
	) error

	RecordPlayback(
		PlaybackObservation,
	) error
}
