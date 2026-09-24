package musicbrainz

import (
	"testing"
	"time"
)

func TestRecording(t *testing.T) {

	recording := Recording{
		ID:       "recording-id",
		Artist:   "Dua Lipa",
		Title:    "Blow Your Mind (Mwah)",
		Score:    100,
		Duration: 178 * time.Second,
		ISRCs: []string{
			"GBAHT1600302",
		},
	}

	if recording.ID != "recording-id" {
		t.Errorf(
			"expected ID %q, got %q",
			"recording-id",
			recording.ID,
		)
	}

	if recording.Duration != 178*time.Second {
		t.Errorf(
			"expected duration %s, got %s",
			178*time.Second,
			recording.Duration,
		)
	}

	if len(recording.ISRCs) != 1 {
		t.Fatalf(
			"expected one ISRC, got %d",
			len(recording.ISRCs),
		)
	}

}
