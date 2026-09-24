package metadata

import (
	"testing"
	"time"
)

func TestCanonicalTrackIdentifier(t *testing.T) {

	track := CanonicalTrack{
		Artist: "Hozier",
		Title:  "Too Sweet",
		Identifiers: []Identifier{
			{
				Scheme: IdentifierMusicBrainz,
				Value:  "recording-id",
			},
		},
	}

	if len(track.Identifiers) != 1 {
		t.Fatalf(
			"expected one identifier, got %d",
			len(track.Identifiers),
		)
	}

	identifier := track.Identifiers[0]

	if identifier.Scheme != IdentifierMusicBrainz {
		t.Errorf(
			"expected scheme %q, got %q",
			IdentifierMusicBrainz,
			identifier.Scheme,
		)
	}

	if identifier.Value != "recording-id" {
		t.Errorf(
			"expected value %q, got %q",
			"recording-id",
			identifier.Value,
		)
	}

}

func TestEnrichmentCandidate(t *testing.T) {

	candidate := EnrichmentCandidate{
		Artist:   "Hozier",
		Title:    "Too Sweet",
		Duration: 251 * time.Second,
		Identifiers: []Identifier{
			{
				Scheme: IdentifierMusicBrainz,
				Value:  "recording-id",
			},
		},
		Provider:      "musicbrainz",
		ProviderScore: 1,
	}

	if candidate.Provider != "musicbrainz" {
		t.Errorf(
			"expected provider %q, got %q",
			"musicbrainz",
			candidate.Provider,
		)
	}

	if candidate.ProviderScore != 1 {
		t.Errorf(
			"expected provider score 1, got %f",
			candidate.ProviderScore,
		)
	}

}

func TestEnrichmentMatchPreservesObservedIdentity(t *testing.T) {

	observed := PlaybackItem{
		Type:   PlaybackTrack,
		Artist: "HOZIER",
		Title:  "Too Sweet",
	}

	match := EnrichmentMatch{
		Track: CanonicalTrack{
			Artist: "Hozier",
			Title:  "Too Sweet",
		},
		Provider:   "test",
		Confidence: 1,
	}

	if observed.Artist != "HOZIER" {
		t.Errorf(
			"expected observed artist to remain %q, got %q",
			"HOZIER",
			observed.Artist,
		)
	}

	if match.Track.Artist != "Hozier" {
		t.Errorf(
			"expected canonical artist %q, got %q",
			"Hozier",
			match.Track.Artist,
		)
	}

}
