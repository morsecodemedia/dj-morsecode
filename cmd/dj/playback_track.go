package main

import (
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

type playbackTrack struct {
	Track metadata.NowPlaying

	Album    string
	Duration time.Duration
	TrackID  string

	Observed bool
}

func (m model) resolvePlaybackTrack(
	snapshot playbackSnapshot,
) (model, playbackTrack) {

	result := playbackTrack{
		Track: metadata.Resolve(
			snapshot.RawTitle,
		),
		Album:    snapshot.Album,
		Duration: snapshot.Duration,
		TrackID:  snapshot.TrackID,
	}

	if snapshot.IsRadio {

		observedItem := metadata.Normalize(
			snapshot.Metadata,
		)

		if observedItem.Observed() {

			result.Observed = true
			m.PlaybackItem = observedItem

		}

		if m.PlaybackItem.IsTrack() {

			result.Track.RawTitle =
				m.PlaybackItem.RawTitle

			result.Track.Artist =
				m.PlaybackItem.Artist

			result.Track.Title =
				m.PlaybackItem.Title

			result.Track.Valid = true

			result.Album =
				m.PlaybackItem.Album

			if m.PlaybackItem.Duration > 0 {

				result.Duration =
					m.PlaybackItem.Duration

			}

			result.TrackID =
				snapshot.Path +
					"\x00" +
					m.PlaybackItem.Artist +
					"\x00" +
					m.PlaybackItem.Title

		} else {

			result.Track.Valid = false

		}

		return m, result
	}

	m.PlaybackItem =
		metadata.PlaybackItem{}

	if snapshot.Artist != "" {

		result.Track.Artist =
			snapshot.Artist

	}

	if snapshot.Title != "" {

		result.Track.Title =
			snapshot.Title

	}

	return m, result

}
