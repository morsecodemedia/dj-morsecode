package metadata

import "testing"

func TestDetectSource(t *testing.T) {

	tests := []struct {
		name     string
		fixture  string
		expected Source
	}{
		{
			name:     "Z100",
			fixture:  "z100-track.json",
			expected: SourceIHeart,
		},
		{
			name:     "SomaFM",
			fixture:  "somafm-track.json",
			expected: SourceSomaFM,
		},
		{
			name:     "SKAfari",
			fixture:  "skafari-track.json",
			expected: SourceLautFM,
		},
		{
			name:     "WXPN",
			fixture:  "wxpn-program.json",
			expected: SourceWXPN,
		},
		{
			name:     "After Hours FM",
			fixture:  "ahfm-program.json",
			expected: SourceAfterHoursFM,
		},
		{
			name:     "Radio Paradise",
			fixture:  "radioparadise-track.json",
			expected: SourceRadioParadise,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			fields := loadMetadataFixture(
				t,
				test.fixture,
			)

			source := DetectSource(
				fields,
			)

			if source != test.expected {

				t.Errorf(
					"expected source %q, got %q",
					test.expected,
					source,
				)

			}

		})

	}

}

func TestDetectSourceUnknown(t *testing.T) {

	source := DetectSource(
		map[string]string{
			"icy-url": "https://example.com/stream",
		},
	)

	if source != SourceGeneric {
		t.Errorf(
			"expected generic source, got %q",
			source,
		)
	}

}
