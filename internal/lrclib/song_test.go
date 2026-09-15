package lrclib

import (
	"testing"
	"time"
)

func TestSong(t *testing.T) {

	result := Result{
		ID:         123,
		TrackName:  "Interstate Love Song",
		ArtistName: "Stone Temple Pilots",
		AlbumName:  "Purple",
		Duration:   196,
		SyncedLyrics: `[00:13.00]Waiting on a Sunday afternoon
[00:18.00]For what I read between the lines`,
	}

	song := Song(result)

	if song.Title != "Interstate Love Song" {
		t.Errorf(
			"expected title %q, got %q",
			"Interstate Love Song",
			song.Title,
		)
	}

	if song.Artist != "Stone Temple Pilots" {
		t.Errorf(
			"expected artist %q, got %q",
			"Stone Temple Pilots",
			song.Artist,
		)
	}

	if song.Album != "Purple" {
		t.Errorf(
			"expected album %q, got %q",
			"Purple",
			song.Album,
		)
	}

	if song.Duration != 196*time.Second {
		t.Errorf(
			"expected duration %s, got %s",
			196*time.Second,
			song.Duration,
		)
	}

	if len(song.Timeline) != 2 {
		t.Fatalf(
			"expected 2 timeline cues, got %d",
			len(song.Timeline),
		)
	}

	if song.Timeline[0].Time != 13*time.Second {
		t.Errorf(
			"expected first cue at %s, got %s",
			13*time.Second,
			song.Timeline[0].Time,
		)
	}

	if song.Timeline[0].Text != "Waiting on a Sunday afternoon" {
		t.Errorf(
			"unexpected first cue text %q",
			song.Timeline[0].Text,
		)
	}

}
