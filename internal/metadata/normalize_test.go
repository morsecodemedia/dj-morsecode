package metadata

import "testing"

func TestNormalizeSomaFMTrack(t *testing.T) {

	fields := loadMetadataFixture(
		t,
		"somafm-track.json",
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

	if item.Artist != "Carbon Based Lifeforms" {
		t.Errorf(
			"expected artist %q, got %q",
			"Carbon Based Lifeforms",
			item.Artist,
		)
	}

	if item.Title != "Dreamshore Forest [Analog Remake]" {
		t.Errorf(
			"expected title %q, got %q",
			"Dreamshore Forest [Analog Remake]",
			item.Title,
		)
	}

}

func TestNormalizeRadioParadiseTrack(t *testing.T) {

	fields := loadMetadataFixture(
		t,
		"radioparadise-track.json",
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

	if item.Artist != "Fleetwood Mac" {
		t.Errorf(
			"expected artist %q, got %q",
			"Fleetwood Mac",
			item.Artist,
		)
	}

	if item.Title != "Coming Your Way" {
		t.Errorf(
			"expected title %q, got %q",
			"Coming Your Way",
			item.Title,
		)
	}

}

func TestNormalizeEmptyMetadata(t *testing.T) {

	item := Normalize(
		map[string]string{},
	)

	if item.Type != PlaybackUnknown {
		t.Fatalf(
			"expected unknown, got %q",
			item.Type,
		)
	}

	if item.IsTrack() {
		t.Fatal("expected empty metadata not to be a track")
	}

}
