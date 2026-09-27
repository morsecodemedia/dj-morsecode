package metadata

import "testing"

func TestSelectReleaseContextBobJames(
	t *testing.T,
) {

	evidence := []ReleaseEvidence{
		{
			ReleaseTitle: "Touchdown",
			ReleaseDate:  "1995-10-17",

			GroupID:               "touchdown",
			GroupTitle:            "Touchdown",
			GroupFirstReleaseDate: "1978",

			PrimaryType: "Album",
			Provider:    "musicbrainz",
		},
		{
			ReleaseTitle: "The Best of Bob James",
			ReleaseDate:  "1996-04-01",

			GroupID:               "best-of",
			GroupTitle:            "The Best of Bob James",
			GroupFirstReleaseDate: "1996-04-01",

			PrimaryType: "Album",
			SecondaryTypes: []string{
				"Compilation",
			},

			Provider: "musicbrainz",
		},
		{
			ReleaseTitle: "Touchdown",
			ReleaseDate:  "1978",

			GroupID:               "touchdown",
			GroupTitle:            "Touchdown",
			GroupFirstReleaseDate: "1978",

			PrimaryType: "Album",
			Provider:    "musicbrainz",
		},
	}

	context, ok := SelectReleaseContext(
		evidence,
	)
	if !ok {

		t.Fatal(
			"expected release context",
		)

	}

	if context.Title != "Touchdown" {

		t.Errorf(
			"expected title %q, got %q",
			"Touchdown",
			context.Title,
		)

	}

	if context.Date != "1978" {

		t.Errorf(
			"expected date %q, got %q",
			"1978",
			context.Date,
		)

	}

}

func TestSelectReleaseContextHozier(
	t *testing.T,
) {

	evidence := []ReleaseEvidence{
		{
			ReleaseTitle: "Unheard",
			ReleaseDate:  "2024-03-22",

			GroupID:               "unheard",
			GroupTitle:            "Unheard",
			GroupFirstReleaseDate: "2024-03-22",

			PrimaryType: "EP",
			Provider:    "musicbrainz",
		},
		{
			ReleaseTitle: "Unreal Unearth",
			ReleaseDate:  "2024-03-22",

			GroupID:               "unreal-unearth",
			GroupTitle:            "Unreal Unearth",
			GroupFirstReleaseDate: "2023-08-18",

			PrimaryType: "Album",
			Provider:    "musicbrainz",
		},
		{
			ReleaseTitle: "Bravo Hits 125",
			ReleaseDate:  "2024-04-26",

			GroupID:               "bravo-hits",
			GroupTitle:            "Bravo Hits 125",
			GroupFirstReleaseDate: "2024-04-26",

			PrimaryType: "Album",
			SecondaryTypes: []string{
				"Compilation",
			},

			Provider: "musicbrainz",
		},
	}

	context, ok := SelectReleaseContext(
		evidence,
	)
	if !ok {

		t.Fatal(
			"expected release context",
		)

	}

	if context.Title != "Unheard" {

		t.Errorf(
			"expected title %q, got %q",
			"Unheard",
			context.Title,
		)

	}

	if context.Date != "2024-03-22" {

		t.Errorf(
			"expected date %q, got %q",
			"2024-03-22",
			context.Date,
		)

	}

}

func TestSelectReleaseContextRejectsAmbiguousOriginals(
	t *testing.T,
) {

	evidence := []ReleaseEvidence{
		{
			ReleaseTitle: "Original Album",
			ReleaseDate:  "2000-01-01",

			GroupID:               "album",
			GroupTitle:            "Original Album",
			GroupFirstReleaseDate: "2000-01-01",

			PrimaryType: "Album",
		},
		{
			ReleaseTitle: "Original Single",
			ReleaseDate:  "2000-01-01",

			GroupID:               "single",
			GroupTitle:            "Original Single",
			GroupFirstReleaseDate: "2000-01-01",

			PrimaryType: "Single",
		},
	}

	_, ok := SelectReleaseContext(
		evidence,
	)

	if ok {

		t.Fatal(
			"expected ambiguous release context to be rejected",
		)

	}

}

func TestSelectReleaseContextRejectsCompilationOnly(
	t *testing.T,
) {

	evidence := []ReleaseEvidence{
		{
			ReleaseTitle: "Greatest Hits Forever",
			ReleaseDate:  "2010",

			GroupID:               "compilation",
			GroupTitle:            "Greatest Hits Forever",
			GroupFirstReleaseDate: "2010",

			PrimaryType: "Album",
			SecondaryTypes: []string{
				"Compilation",
			},
		},
	}

	_, ok := SelectReleaseContext(
		evidence,
	)

	if ok {

		t.Fatal(
			"expected compilation-only evidence to be rejected",
		)

	}

}

func TestSelectReleaseContextRejectsMissingDateEvidence(
	t *testing.T,
) {

	evidence := []ReleaseEvidence{
		{
			ReleaseTitle: "Mystery Album",

			GroupID:    "mystery",
			GroupTitle: "Mystery Album",

			PrimaryType: "Album",
		},
	}

	_, ok := SelectReleaseContext(
		evidence,
	)

	if ok {

		t.Fatal(
			"expected incomplete release evidence to be rejected",
		)

	}

}
