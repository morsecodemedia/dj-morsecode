package source

import "testing"

func TestMediaItemValid(
	t *testing.T,
) {

	item := MediaItem{
		Kind: MediaTrack,
		Ref: ItemRef{
			Source: Source{
				Kind: KindLastFM,
			},
			ID: "track-123",
		},
		Artist: "Stone Temple Pilots",
		Title:  "Interstate Love Song",
	}

	if !item.Valid() {
		t.Fatal(
			"expected media item to be valid",
		)
	}

}

func TestMediaItemRejectsUnknownKind(
	t *testing.T,
) {

	item := MediaItem{
		Ref: ItemRef{
			Source: Source{
				Kind: KindSpotify,
			},
			ID: "track-123",
		},
	}

	if item.Valid() {
		t.Fatal(
			"expected unknown media kind to be invalid",
		)
	}

}

func TestMediaItemRejectsInvalidReference(
	t *testing.T,
) {

	item := MediaItem{
		Kind: MediaTrack,
	}

	if item.Valid() {
		t.Fatal(
			"expected media item without provenance to be invalid",
		)
	}

}
