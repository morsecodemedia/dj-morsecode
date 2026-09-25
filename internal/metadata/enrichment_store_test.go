package metadata

import "testing"

func TestEnrichmentStorePutAndGet(
	t *testing.T,
) {

	store := NewEnrichmentStore()

	match := EnrichmentMatch{
		Track: CanonicalTrack{
			Artist: "Hozier",
			Title:  "Too Sweet",
			Identifiers: []Identifier{
				{
					Scheme: IdentifierMusicBrainz,
					Value:  "recording-id",
				},
			},
		},
		Provider:   "musicbrainz",
		Confidence: 1,
	}

	store.Put(
		"Hozier",
		"Too Sweet",
		match,
	)

	actual, ok := store.Get(
		"HOZIER",
		" too sweet ",
	)
	if !ok {

		t.Fatal(
			"expected cached enrichment match",
		)

	}

	if actual.Track.Artist != "Hozier" {

		t.Errorf(
			"expected artist %q, got %q",
			"Hozier",
			actual.Track.Artist,
		)

	}

	if len(actual.Track.Identifiers) != 1 {

		t.Fatalf(
			"expected one identifier, got %d",
			len(actual.Track.Identifiers),
		)

	}

}

func TestEnrichmentStoreMiss(
	t *testing.T,
) {

	store := NewEnrichmentStore()

	_, ok := store.Get(
		"Hozier",
		"Too Sweet",
	)

	if ok {

		t.Fatal(
			"expected enrichment cache miss",
		)

	}

}
