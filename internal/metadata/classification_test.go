package metadata

import "testing"

func TestNormalizeFixtureClassification(t *testing.T) {

	tests := []struct {
		name     string
		fixture  string
		expected PlaybackItemType
	}{
		{
			name:     "Z100 track",
			fixture:  "z100-track.json",
			expected: PlaybackTrack,
		},
		{
			name:     "Z100 station ID",
			fixture:  "z100-station-id.json",
			expected: PlaybackStationID,
		},
		{
			name:     "SomaFM track",
			fixture:  "somafm-track.json",
			expected: PlaybackTrack,
		},
		{
			name:     "SKAfari track",
			fixture:  "skafari-track.json",
			expected: PlaybackTrack,
		},
		{
			name:     "SKAfari advertisement",
			fixture:  "skafari-ad.json",
			expected: PlaybackUnknown,
		},
		{
			name:     "WXPN program",
			fixture:  "wxpn-program.json",
			expected: PlaybackUnknown,
		},
		{
			name:     "AH.FM program",
			fixture:  "ahfm-program.json",
			expected: PlaybackUnknown,
		},
		{
			name:     "Radio Paradise track",
			fixture:  "radioparadise-track.json",
			expected: PlaybackTrack,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			fields := loadMetadataFixture(
				t,
				test.fixture,
			)

			item := Normalize(
				fields,
			)

			if item.Type != test.expected {

				t.Errorf(
					"expected %q, got %q for raw title %q",
					test.expected,
					item.Type,
					item.RawTitle,
				)

			}

		})

	}

}
