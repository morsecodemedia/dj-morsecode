package musicbrainz

import (
	"testing"
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

func TestCandidate(t *testing.T) {

	recording := Recording{
		ID:       "recording-id",
		Artist:   "Hozier",
		Title:    "Too Sweet",
		Score:    98,
		Duration: 251 * time.Second,
		ISRCs: []string{
			"USSM12401865",
		},
	}

	candidate := Candidate(
		recording,
	)

	if candidate.Artist != "Hozier" {
		t.Errorf(
			"expected artist %q, got %q",
			"Hozier",
			candidate.Artist,
		)
	}

	if candidate.Title != "Too Sweet" {
		t.Errorf(
			"expected title %q, got %q",
			"Too Sweet",
			candidate.Title,
		)
	}

	if candidate.Provider != providerName {
		t.Errorf(
			"expected provider %q, got %q",
			providerName,
			candidate.Provider,
		)
	}

	if candidate.ProviderScore != 0.98 {
		t.Errorf(
			"expected provider score 0.98, got %f",
			candidate.ProviderScore,
		)
	}

	if candidate.Duration != 251*time.Second {
		t.Errorf(
			"expected duration %s, got %s",
			251*time.Second,
			candidate.Duration,
		)
	}

	if len(candidate.Identifiers) != 2 {
		t.Fatalf(
			"expected two identifiers, got %d",
			len(candidate.Identifiers),
		)
	}

	if candidate.Identifiers[0].Scheme !=
		metadata.IdentifierMusicBrainz {

		t.Errorf(
			"expected first identifier scheme %q, got %q",
			metadata.IdentifierMusicBrainz,
			candidate.Identifiers[0].Scheme,
		)

	}

	if candidate.Identifiers[0].Value !=
		"recording-id" {

		t.Errorf(
			"expected recording ID %q, got %q",
			"recording-id",
			candidate.Identifiers[0].Value,
		)

	}

	if candidate.Identifiers[1].Scheme !=
		metadata.IdentifierISRC {

		t.Errorf(
			"expected second identifier scheme %q, got %q",
			metadata.IdentifierISRC,
			candidate.Identifiers[1].Scheme,
		)

	}

	if candidate.Identifiers[1].Value !=
		"USSM12401865" {

		t.Errorf(
			"expected ISRC %q, got %q",
			"USSM12401865",
			candidate.Identifiers[1].Value,
		)

	}

}
func TestCandidateSkipsEmptyISRC(t *testing.T) {

	recording := Recording{
		ID:     "recording-id",
		Artist: "Hozier",
		Title:  "Too Sweet",
		Score:  100,
		ISRCs: []string{
			"",
			"USSM12401865",
			"",
		},
	}

	candidate := Candidate(
		recording,
	)

	if len(candidate.Identifiers) != 2 {
		t.Fatalf(
			"expected recording ID and one ISRC, got %d identifiers",
			len(candidate.Identifiers),
		)
	}

}
