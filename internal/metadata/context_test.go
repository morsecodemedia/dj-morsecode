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
		Tags: []ContextTag{
			{
				Name:        "soul",
				TrackCount:  16,
				ArtistCount: 56,
				Providers: []string{
					"lastfm",
				},
			},
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
		ReleaseTitle: "Unheard",
		ReleaseDate:  "2024-03-22",

		GroupID:               "release-group-id",
		GroupTitle:            "Unheard",
		GroupFirstReleaseDate: "2024-03-22",

		PrimaryType: "EP",
		SecondaryTypes: []string{
			"Compilation",
		},

		Provider: "musicbrainz",
	}

	if evidence.ReleaseTitle != "Unheard" {
		t.Errorf(
			"expected release title %q, got %q",
			"Unheard",
			evidence.ReleaseTitle,
		)
	}

	if evidence.ReleaseDate != "2024-03-22" {
		t.Errorf(
			"expected release date %q, got %q",
			"2024-03-22",
			evidence.ReleaseDate,
		)
	}

	if evidence.GroupID != "release-group-id" {
		t.Errorf(
			"expected release group ID %q, got %q",
			"release-group-id",
			evidence.GroupID,
		)
	}

	if evidence.GroupTitle != "Unheard" {
		t.Errorf(
			"expected release group title %q, got %q",
			"Unheard",
			evidence.GroupTitle,
		)
	}

	if evidence.GroupFirstReleaseDate !=
		"2024-03-22" {

		t.Errorf(
			"expected group first release date %q, got %q",
			"2024-03-22",
			evidence.GroupFirstReleaseDate,
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

func TestTagEvidence(t *testing.T) {

	evidence := TagEvidence{
		Name:     "hard rock",
		Count:    100,
		Provider: "lastfm",
	}

	if evidence.Name != "hard rock" {
		t.Errorf(
			"expected tag name %q, got %q",
			"hard rock",
			evidence.Name,
		)
	}

	if evidence.Count != 100 {
		t.Errorf(
			"expected count %d, got %d",
			100,
			evidence.Count,
		)
	}

	if evidence.Provider != "lastfm" {
		t.Errorf(
			"expected provider %q, got %q",
			"lastfm",
			evidence.Provider,
		)
	}

}

func TestContextTag(t *testing.T) {

	tag := ContextTag{
		Name:        "hard rock",
		TrackCount:  0,
		ArtistCount: 100,
		Providers: []string{
			"lastfm",
		},
	}

	if tag.Name != "hard rock" {
		t.Errorf(
			"expected name %q, got %q",
			"hard rock",
			tag.Name,
		)
	}

	if tag.TrackCount != 0 {
		t.Errorf(
			"expected track count 0, got %d",
			tag.TrackCount,
		)
	}

	if tag.ArtistCount != 100 {
		t.Errorf(
			"expected artist count 100, got %d",
			tag.ArtistCount,
		)
	}

	if len(tag.Providers) != 1 {
		t.Fatalf(
			"expected one provider, got %d",
			len(tag.Providers),
		)
	}

}
