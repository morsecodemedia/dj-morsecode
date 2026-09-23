package metadata

import "testing"

func TestIsLautFMAdvertisement(t *testing.T) {

	tests := []struct {
		name     string
		rawTitle string
		expected bool
	}{
		{
			name:     "consumer information",
			rawTitle: "hollywoodcasino.com - Ein bisschen Verbraucherinformationen ...",
			expected: true,
		},
		{
			name:     "advertisement progress",
			rawTitle: "hollywoodcasino.com - ... schon fast die Hälfte geschafft ...",
			expected: true,
		},
		{
			name:     "music",
			rawTitle: "No Respect - No Nazi's Friend",
			expected: false,
		},
		{
			name:     "domain advertisement",
			rawTitle: "nissan-global.com - ... schon fast die Hälfte geschafft ...",
			expected: true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			actual := isLautFMAdvertisement(
				test.rawTitle,
			)

			if actual != test.expected {

				t.Errorf(
					"expected %t, got %t for %q",
					test.expected,
					actual,
					test.rawTitle,
				)

			}

		})

	}

}

func TestIsDomainLike(t *testing.T) {

	tests := []struct {
		value    string
		expected bool
	}{
		{
			value:    "hollywoodcasino.com",
			expected: true,
		},
		{
			value:    "nissan-global.com",
			expected: true,
		},
		{
			value:    "No Respect",
			expected: false,
		},
		{
			value:    "De Hardheid",
			expected: false,
		},
	}

	for _, test := range tests {

		t.Run(test.value, func(t *testing.T) {

			actual := isDomainLike(
				test.value,
			)

			if actual != test.expected {

				t.Errorf(
					"expected %t, got %t for %q",
					test.expected,
					actual,
					test.value,
				)

			}

		})

	}

}
