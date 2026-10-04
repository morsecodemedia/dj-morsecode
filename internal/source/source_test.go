package source

import "testing"

func TestSourceValid(
	t *testing.T,
) {

	source := Source{
		Kind: KindM3U,
		URI:  "stations.m3u",
	}

	if !source.Valid() {
		t.Fatal(
			"expected source to be valid",
		)
	}

}

func TestSourceRejectsMissingKind(
	t *testing.T,
) {

	source := Source{
		URI: "stations.m3u",
	}

	if source.Valid() {
		t.Fatal(
			"expected source without kind to be invalid",
		)
	}

}

func TestItemRefValidWithID(
	t *testing.T,
) {

	ref := ItemRef{
		Source: Source{
			Kind: KindSpotify,
		},
		ID: "track-123",
	}

	if !ref.Valid() {
		t.Fatal(
			"expected item reference with ID to be valid",
		)
	}

}

func TestItemRefValidWithURI(
	t *testing.T,
) {

	ref := ItemRef{
		Source: Source{
			Kind: KindM3U,
			URI:  "stations.m3u",
		},
		URI: "https://example.com/live",
	}

	if !ref.Valid() {
		t.Fatal(
			"expected item reference with URI to be valid",
		)
	}

}

func TestItemRefRejectsMissingIdentity(
	t *testing.T,
) {

	ref := ItemRef{
		Source: Source{
			Kind: KindLastFM,
		},
		Name: "Interesting Track",
	}

	if ref.Valid() {
		t.Fatal(
			"expected item reference without ID or URI to be invalid",
		)
	}

}

func TestItemRefRejectsInvalidSource(
	t *testing.T,
) {

	ref := ItemRef{
		ID: "track-123",
	}

	if ref.Valid() {
		t.Fatal(
			"expected item reference with invalid source to be invalid",
		)
	}

}
