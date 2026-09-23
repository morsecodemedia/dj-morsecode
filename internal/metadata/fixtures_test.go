package metadata

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func loadMetadataFixture(
	t *testing.T,
	name string,
) map[string]string {

	t.Helper()

	path := filepath.Join(
		"testdata",
		name,
	)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf(
			"read fixture %q: %v",
			name,
			err,
		)
	}

	var fields map[string]string

	if err := json.Unmarshal(
		data,
		&fields,
	); err != nil {

		t.Fatalf(
			"decode fixture %q: %v",
			name,
			err,
		)

	}

	return fields
}

func TestMetadataFixtures(t *testing.T) {

	fixtures := []string{
		"ahfm-program.json",
		"radioparadise-track.json",
		"skafari-ad-domain.json",
		"skafari-ad-progress.json",
		"skafari-ad.json",
		"skafari-track.json",
		"somafm-track.json",
		"wxpn-program.json",
		"z100-no-item.json",
		"z100-spot-block-end.json",
		"z100-spot.json",
		"z100-station-id.json",
		"z100-track-fields-2.json",
		"z100-track-fields.json",
		"z100-track.json",
	}

	for _, fixture := range fixtures {

		t.Run(fixture, func(t *testing.T) {

			fields := loadMetadataFixture(
				t,
				fixture,
			)

			if len(fields) == 0 {
				t.Fatal("expected metadata fields")
			}

		})

	}

}
