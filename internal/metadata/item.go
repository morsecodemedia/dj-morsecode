package metadata

import "time"

type PlaybackItemType string

const (
	PlaybackUnknown       PlaybackItemType = ""
	PlaybackTrack         PlaybackItemType = "track"
	PlaybackAdvertisement PlaybackItemType = "advertisement"
	PlaybackStationID     PlaybackItemType = "station-id"
	PlaybackProgram       PlaybackItemType = "program"
)

type PlaybackItem struct {
	Type PlaybackItemType

	Artist string
	Title  string
	Album  string

	RawTitle string
	Duration time.Duration

	SourceFields map[string]string
}

func (i PlaybackItem) IsTrack() bool {
	return i.Type == PlaybackTrack
}
