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
			name:     "SKAfari advertisement progress",
			fixture:  "skafari-ad-progress.json",
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
		{
			name:     "Z100 field track",
			fixture:  "z100-track-fields.json",
			expected: PlaybackTrack,
		},
		{
			name:     "Z100 second field track",
			fixture:  "z100-track-fields-2.json",
			expected: PlaybackTrack,
		},
		{
			name:     "Z100 spot",
			fixture:  "z100-spot.json",
			expected: PlaybackStationID,
		},
		{
			name:     "Z100 spot block end",
			fixture:  "z100-spot-block-end.json",
			expected: PlaybackUnknown,
		},
		{
			name:     "SKAfari domain advertisement",
			fixture:  "skafari-ad-domain.json",
			expected: PlaybackUnknown,
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
