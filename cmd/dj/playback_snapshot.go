package main

import (
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/player"
)

type playbackSnapshot struct {
	Duration time.Duration

	RawTitle string
	Artist   string
	Title    string
	Album    string
	TrackID  string
	Path     string

	IsNetwork bool
	IsIdle    bool

	Metadata map[string]string
}

func snapshotPlayback(
	playback *player.Player,
) playbackSnapshot {

	if playback == nil {
		return playbackSnapshot{}
	}

	snapshot := playbackSnapshot{
		Duration: playback.Duration(),
		RawTitle: playback.Title(),
		Artist:   playback.Artist(),
		Title:    playback.TrackTitle(),
		Album:    playback.Album(),
		TrackID:  playback.TrackID(),
		Path:     playback.Path(),

		IsNetwork: playback.IsNetwork(),
		IsIdle:    playback.IsIdle(),
	}

	if snapshot.IsNetwork {
		snapshot.Metadata =
			playback.Metadata()
	}

	return snapshot

}
