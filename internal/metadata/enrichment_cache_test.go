package metadata

import "testing"

func TestEnrichmentCacheKeyIsDeterministic(
	t *testing.T,
) {

	first := EnrichmentCacheKey(
		"Hozier",
		"Too Sweet",
	)

	second := EnrichmentCacheKey(
		"Hozier",
		"Too Sweet",
	)

	if first != second {

		t.Errorf(
			"expected identical keys, got %q and %q",
			first,
			second,
		)

	}

}

func TestEnrichmentCacheKeyNormalizesIdentity(
	t *testing.T,
) {

	first := EnrichmentCacheKey(
		"Hozier",
		"Too Sweet",
	)

	second := EnrichmentCacheKey(
		"  HOZIER ",
		" TOO SWEET  ",
	)

	if first != second {

		t.Errorf(
			"expected normalized identities to match, got %q and %q",
			first,
			second,
		)

	}

}

func TestEnrichmentCacheKeySeparatesArtistAndTitle(
	t *testing.T,
) {

	first := EnrichmentCacheKey(
		"ab",
		"c",
	)

	second := EnrichmentCacheKey(
		"a",
		"bc",
	)

	if first == second {

		t.Errorf(
			"expected different keys, both were %q",
			first,
		)

	}

}

func TestEnrichmentCacheKeyLength(
	t *testing.T,
) {

	key := EnrichmentCacheKey(
		"Hozier",
		"Too Sweet",
	)

	if len(key) != 64 {

		t.Errorf(
			"expected 64-character SHA-256 key, got %d",
			len(key),
		)

	}

}
