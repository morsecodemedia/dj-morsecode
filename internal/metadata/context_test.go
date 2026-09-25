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
