package metadata

import (
	"testing"
	"time"
)

func TestMatchEnrichment(t *testing.T) {

	observed := PlaybackItem{
		Type:     PlaybackTrack,
		Artist:   "HOZIER",
		Title:    "Too Sweet",
		Duration: 251 * time.Second,
	}

	candidates := []EnrichmentCandidate{
		{
			Artist:   "Hozier",
			Title:    "Too Sweet",
			Duration: 250 * time.Second,

			Identifiers: []Identifier{
				{
					Scheme: IdentifierMusicBrainz,
					Value:  "recording-id",
				},
			},

			Provider:      "musicbrainz",
			ProviderScore: 1,
		},
	}

	match, ok := MatchEnrichment(
		observed,
		candidates,
	)
	if !ok {
		t.Fatal("expected enrichment match")
	}

	if match.Track.Artist != "Hozier" {
		t.Errorf(
			"expected canonical artist %q, got %q",
			"Hozier",
			match.Track.Artist,
		)
	}

	if match.Track.Title != "Too Sweet" {
		t.Errorf(
			"expected canonical title %q, got %q",
			"Too Sweet",
			match.Track.Title,
		)
	}

	if match.Provider != "musicbrainz" {
		t.Errorf(
			"expected provider %q, got %q",
			"musicbrainz",
			match.Provider,
		)
	}

}

func TestMatchEnrichmentRejectsWrongArtist(t *testing.T) {

	observed := PlaybackItem{
		Type:   PlaybackTrack,
		Artist: "Queen",
		Title:  "Somebody to Love",
	}

	candidates := []EnrichmentCandidate{
		{
			Artist:        "Jefferson Airplane",
			Title:         "Somebody to Love",
			Provider:      "musicbrainz",
			ProviderScore: 1,
		},
	}

	_, ok := MatchEnrichment(
		observed,
		candidates,
	)

	if ok {
		t.Fatal("expected mismatched artist to be rejected")
	}

}

func TestMatchEnrichmentRejectsWrongTitle(t *testing.T) {

	observed := PlaybackItem{
		Type:   PlaybackTrack,
		Artist: "Hozier",
		Title:  "Too Sweet",
	}

	candidates := []EnrichmentCandidate{
		{
			Artist:        "Hozier",
			Title:         "Take Me to Church",
			Provider:      "musicbrainz",
			ProviderScore: 1,
		},
	}

	_, ok := MatchEnrichment(
		observed,
		candidates,
	)

	if ok {
		t.Fatal("expected mismatched title to be rejected")
	}

}

func TestMatchEnrichmentRejectsLowProviderScore(t *testing.T) {

	observed := PlaybackItem{
		Type:   PlaybackTrack,
		Artist: "Hozier",
		Title:  "Too Sweet",
	}

	candidates := []EnrichmentCandidate{
		{
			Artist:        "Hozier",
			Title:         "Too Sweet",
			Provider:      "musicbrainz",
			ProviderScore: 0.75,
		},
	}

	_, ok := MatchEnrichment(
		observed,
		candidates,
	)

	if ok {
		t.Fatal("expected low-score candidate to be rejected")
	}

}

func TestMatchEnrichmentRejectsDurationMismatch(t *testing.T) {

	observed := PlaybackItem{
		Type:     PlaybackTrack,
		Artist:   "Hozier",
		Title:    "Too Sweet",
		Duration: 251 * time.Second,
	}

	candidates := []EnrichmentCandidate{
		{
			Artist:        "Hozier",
			Title:         "Too Sweet",
			Duration:      300 * time.Second,
			Provider:      "musicbrainz",
			ProviderScore: 1,
		},
	}

	_, ok := MatchEnrichment(
		observed,
		candidates,
	)

	if ok {
		t.Fatal("expected duration mismatch to be rejected")
	}

}

func TestMatchEnrichmentAllowsMissingDuration(t *testing.T) {

	observed := PlaybackItem{
		Type:   PlaybackTrack,
		Artist: "Hozier",
		Title:  "Too Sweet",
	}

	candidates := []EnrichmentCandidate{
		{
			Artist:        "Hozier",
			Title:         "Too Sweet",
			Duration:      251 * time.Second,
			Provider:      "musicbrainz",
			ProviderScore: 1,
		},
	}

	_, ok := MatchEnrichment(
		observed,
		candidates,
	)

	if !ok {
		t.Fatal(
			"expected missing observed duration not to prevent match",
		)
	}

}

func TestMatchEnrichmentRejectsNonTrack(t *testing.T) {

	observed := PlaybackItem{
		Type:  PlaybackStationID,
		Title: "Z100",
	}

	candidates := []EnrichmentCandidate{
		{
			Artist:        "Z100",
			Title:         "Z100",
			Provider:      "musicbrainz",
			ProviderScore: 1,
		},
	}

	_, ok := MatchEnrichment(
		observed,
		candidates,
	)

	if ok {
		t.Fatal(
			"expected non-track playback item to be rejected",
		)
	}

}
func TestMatchEnrichmentRejectsAmbiguousCandidates(
	t *testing.T,
) {

	observed := PlaybackItem{
		Type:     PlaybackTrack,
		Artist:   "Hozier",
		Title:    "Too Sweet",
		Duration: 251 * time.Second,
	}

	candidates := []EnrichmentCandidate{
		{
			Artist:        "Hozier",
			Title:         "Too Sweet",
			Duration:      251*time.Second + 424*time.Millisecond,
			Provider:      "musicbrainz",
			ProviderScore: 1,
			Identifiers: []Identifier{
				{
					Scheme: IdentifierMusicBrainz,
					Value:  "recording-a",
				},
			},
		},
		{
			Artist:        "Hozier",
			Title:         "Too Sweet",
			Duration:      251 * time.Second,
			Provider:      "musicbrainz",
			ProviderScore: 1,
			Identifiers: []Identifier{
				{
					Scheme: IdentifierMusicBrainz,
					Value:  "recording-b",
				},
			},
		},
	}

	_, ok := MatchEnrichment(
		observed,
		candidates,
	)

	if ok {
		t.Fatal(
			"expected ambiguous candidates to be rejected",
		)
	}

}

func TestMatchEnrichmentPrefersSingleUnqualifiedCandidate(
	t *testing.T,
) {

	observed := PlaybackItem{
		Type:     PlaybackTrack,
		Artist:   "Hozier",
		Title:    "Too Sweet",
		Duration: 251 * time.Second,
	}

	candidates := []EnrichmentCandidate{
		{
			Artist:        "Hozier",
			Title:         "Too Sweet",
			Duration:      251*time.Second + 424*time.Millisecond,
			Variant:       "Dolby Atmos mix",
			Provider:      "musicbrainz",
			ProviderScore: 1,
			Identifiers: []Identifier{
				{
					Scheme: IdentifierMusicBrainz,
					Value:  "atmos-recording",
				},
			},
		},
		{
			Artist:        "Hozier",
			Title:         "Too Sweet",
			Duration:      251 * time.Second,
			Provider:      "musicbrainz",
			ProviderScore: 1,
			Identifiers: []Identifier{
				{
					Scheme: IdentifierMusicBrainz,
					Value:  "standard-recording",
				},
			},
		},
	}

	match, ok := MatchEnrichment(
		observed,
		candidates,
	)
	if !ok {
		t.Fatal(
			"expected unqualified candidate to resolve ambiguity",
		)
	}

	if len(match.Track.Identifiers) != 1 {
		t.Fatalf(
			"expected one identifier, got %d",
			len(match.Track.Identifiers),
		)
	}

	if match.Track.Identifiers[0].Value !=
		"standard-recording" {

		t.Errorf(
			"expected standard recording, got %q",
			match.Track.Identifiers[0].Value,
		)

	}

}

func TestMatchEnrichmentRejectsMultipleQualifiedCandidates(
	t *testing.T,
) {

	observed := PlaybackItem{
		Type:   PlaybackTrack,
		Artist: "Hozier",
		Title:  "Too Sweet",
	}

	candidates := []EnrichmentCandidate{
		{
			Artist:        "Hozier",
			Title:         "Too Sweet",
			Variant:       "live",
			Provider:      "musicbrainz",
			ProviderScore: 1,
		},
		{
			Artist:        "Hozier",
			Title:         "Too Sweet",
			Variant:       "remix",
			Provider:      "musicbrainz",
			ProviderScore: 1,
		},
	}

	_, ok := MatchEnrichment(
		observed,
		candidates,
	)

	if ok {
		t.Fatal(
			"expected qualified candidates to remain ambiguous",
		)
	}

}
