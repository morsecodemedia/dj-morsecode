package lyrics

import (
	"testing"
	"time"
)

func TestFromLines(t *testing.T) {

	lines := []string{
		"[ti:Interstate Love Song]",
		"[ar:Stone Temple Pilots]",
		"[al:Purple]",
		"[length:03:16]",
		"[00:13.00]Waiting on a Sunday afternoon",
		"[00:18.00]For what I read between the lines",
	}

	song := FromLines(lines)

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

	if song.Duration != 3*time.Minute+16*time.Second {
		t.Errorf(
			"expected duration %s, got %s",
			3*time.Minute+16*time.Second,
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
func TestFromString(t *testing.T) {

	content := `[ti:Interstate Love Song]
[ar:Stone Temple Pilots]
[al:Purple]
[length:03:16]
[00:13.00]Waiting on a Sunday afternoon
[00:18.00]For what I read between the lines`

	song := FromString(content)

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

	if len(song.Timeline) != 2 {
		t.Fatalf(
			"expected 2 timeline cues, got %d",
			len(song.Timeline),
		)
	}

	if song.Timeline[1].Time != 18*time.Second {
		t.Errorf(
			"expected second cue at %s, got %s",
			18*time.Second,
			song.Timeline[1].Time,
		)
	}

	if song.Timeline[1].Text != "For what I read between the lines" {
		t.Errorf(
			"unexpected second cue text %q",
			song.Timeline[1].Text,
		)
	}

}
