package metadata

import "testing"

func TestTrackContext(t *testing.T) {

	context := TrackContext{
		Release: ReleaseContext{
			Title: "Unheard",
			Date:  "2024-03-22",
		},
		Genres: []string{
			"indie rock",
			"soul",
		},
	}

	if context.Release.Title != "Unheard" {
		t.Errorf(
			"expected release title %q, got %q",
			"Unheard",
			context.Release.Title,
		)
	}

	if context.Release.Date != "2024-03-22" {
		t.Errorf(
			"expected release date %q, got %q",
			"2024-03-22",
			context.Release.Date,
		)
	}

	if len(context.Genres) != 2 {
		t.Fatalf(
			"expected two genres, got %d",
			len(context.Genres),
		)
	}

}

func TestReleaseContextPreservesDatePrecision(
	t *testing.T,
) {

	tests := []string{
		"2024",
		"2024-03",
		"2024-03-22",
	}

	for _, date := range tests {

		t.Run(date, func(t *testing.T) {

			release := ReleaseContext{
				Date: date,
			}

			if release.Date != date {
				t.Errorf(
					"expected date %q, got %q",
					date,
					release.Date,
				)
			}

		})

	}

}

func TestReleaseEvidence(t *testing.T) {

	evidence := ReleaseEvidence{
		Title:          "Unheard",
		Date:           "2024-03-22",
		ReleaseGroupID: "release-group-id",
		PrimaryType:    "EP",
		SecondaryTypes: []string{
			"Compilation",
		},
		Provider: "musicbrainz",
	}

	if evidence.Title != "Unheard" {
		t.Errorf(
			"expected title %q, got %q",
			"Unheard",
			evidence.Title,
		)
	}

	if evidence.ReleaseGroupID !=
		"release-group-id" {

		t.Errorf(
			"expected release group ID %q, got %q",
			"release-group-id",
			evidence.ReleaseGroupID,
		)

	}

	if evidence.PrimaryType != "EP" {
		t.Errorf(
			"expected primary type %q, got %q",
			"EP",
			evidence.PrimaryType,
		)
	}

	if len(evidence.SecondaryTypes) != 1 {
		t.Fatalf(
			"expected one secondary type, got %d",
			len(evidence.SecondaryTypes),
		)
	}

	if evidence.Provider != "musicbrainz" {
		t.Errorf(
			"expected provider %q, got %q",
			"musicbrainz",
			evidence.Provider,
		)
	}

}
