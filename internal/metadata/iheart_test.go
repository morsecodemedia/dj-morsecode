package metadata

import (
	"testing"
	"time"
)

func TestNormalizeIHeartTrack(t *testing.T) {

	fields := loadMetadataFixture(
		t,
		"z100-track.json",
	)

	item := Normalize(
		fields,
	)

	if item.Type != PlaybackTrack {
		t.Fatalf(
			"expected track, got %q",
			item.Type,
		)
	}

	if item.Artist != "Dominic Fike" {
		t.Errorf(
			"expected artist %q, got %q",
			"Dominic Fike",
			item.Artist,
		)
	}

	if item.Title != "Babydoll" {
		t.Errorf(
			"expected title %q, got %q",
			"Babydoll",
			item.Title,
		)
	}

	if item.Duration != 1*time.Minute+37*time.Second {
		t.Errorf(
			"expected duration %s, got %s",
			1*time.Minute+37*time.Second,
			item.Duration,
		)
	}

}

func TestNormalizeIHeartStationID(t *testing.T) {

	fields := loadMetadataFixture(
		t,
		"z100-station-id.json",
	)

	item := Normalize(
		fields,
	)

	if item.Type != PlaybackStationID {
		t.Fatalf(
			"expected station ID, got %q",
			item.Type,
		)
	}

	if item.Title != "Z100" {
		t.Errorf(
			"expected title %q, got %q",
			"Z100",
			item.Title,
		)
	}

	if item.IsTrack() {
		t.Fatal(
			"expected station ID not to be a track",
		)
	}

}

func TestParseIHeartFields(t *testing.T) {

	fields := parseIHeartFields(
		`Dominic Fike - text="Babydoll" song_spot="M" TPID="87757447"`,
	)

	if fields["text"] != "Babydoll" {
		t.Errorf(
			"expected text %q, got %q",
			"Babydoll",
			fields["text"],
		)
	}

	if fields["song_spot"] != "M" {
		t.Errorf(
			"expected song spot %q, got %q",
			"M",
			fields["song_spot"],
		)
	}

	if fields["TPID"] != "87757447" {
		t.Errorf(
			"expected TPID %q, got %q",
			"87757447",
			fields["TPID"],
		)
	}

}
