package musicbrainz

import (
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

func TestReleaseEvidence(t *testing.T) {

	release := Release{
		ID:      "release-id",
		Title:   "Unheard",
		Date:    "2024-03-22",
		Country: "XW",

		ReleaseGroup: ReleaseGroup{
			ID:               "release-group-id",
			Title:            "Unheard",
			FirstReleaseDate: "2024-03-22",
			PrimaryType:      "EP",
			SecondaryTypes: []string{
				"Compilation",
			},
		},
	}

	evidence := ReleaseEvidence(
		release,
	)

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

	if evidence.GroupID !=
		"release-group-id" {

		t.Errorf(
			"expected group ID %q, got %q",
			"release-group-id",
			evidence.GroupID,
		)

	}

	if evidence.GroupTitle != "Unheard" {
		t.Errorf(
			"expected group title %q, got %q",
			"Unheard",
			evidence.GroupTitle,
		)
	}

	if evidence.GroupFirstReleaseDate !=
		"2024-03-22" {

		t.Errorf(
			"expected first release date %q, got %q",
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

	if evidence.Provider != providerName {
		t.Errorf(
			"expected provider %q, got %q",
			providerName,
			evidence.Provider,
		)
	}

}

func TestRecordingReleaseEvidence(
	t *testing.T,
) {

	recording := Recording{
		ID:     "recording-id",
		Artist: "Hozier",
		Title:  "Too Sweet",

		Releases: []Release{
			{
				ID:    "unheard-release",
				Title: "Unheard",
				Date:  "2024-03-22",

				ReleaseGroup: ReleaseGroup{
					ID:               "unheard",
					Title:            "Unheard",
					FirstReleaseDate: "2024-03-22",
					PrimaryType:      "EP",
				},
			},
			{
				ID:    "album-release",
				Title: "Unreal Unearth",
				Date:  "2024-03-22",

				ReleaseGroup: ReleaseGroup{
					ID:               "unreal-unearth",
					Title:            "Unreal Unearth",
					FirstReleaseDate: "2023-08-18",
					PrimaryType:      "Album",
				},
			},
		},
	}

	evidence := RecordingReleaseEvidence(
		recording,
	)

	if len(evidence) != 2 {

		t.Fatalf(
			"expected two release evidence items, got %d",
			len(evidence),
		)

	}

	if evidence[0].GroupTitle != "Unheard" {
		t.Errorf(
			"expected first group %q, got %q",
			"Unheard",
			evidence[0].GroupTitle,
		)
	}

	if evidence[1].GroupTitle !=
		"Unreal Unearth" {

		t.Errorf(
			"expected second group %q, got %q",
			"Unreal Unearth",
			evidence[1].GroupTitle,
		)

	}

}

func TestRecordingReleaseEvidenceSelectsContext(
	t *testing.T,
) {

	recording := Recording{
		ID:     "recording-id",
		Artist: "Hozier",
		Title:  "Too Sweet",

		Releases: []Release{
			{
				Title: "Unheard",
				Date:  "2024-03-22",

				ReleaseGroup: ReleaseGroup{
					ID:               "unheard",
					Title:            "Unheard",
					FirstReleaseDate: "2024-03-22",
					PrimaryType:      "EP",
				},
			},
			{
				Title: "Unreal Unearth",
				Date:  "2024-03-22",

				ReleaseGroup: ReleaseGroup{
					ID:               "unreal-unearth",
					Title:            "Unreal Unearth",
					FirstReleaseDate: "2023-08-18",
					PrimaryType:      "Album",
				},
			},
			{
				Title: "Bravo Hits 125",
				Date:  "2024-04-26",

				ReleaseGroup: ReleaseGroup{
					ID:               "bravo",
					Title:            "Bravo Hits 125",
					FirstReleaseDate: "2024-04-26",
					PrimaryType:      "Album",
					SecondaryTypes: []string{
						"Compilation",
					},
				},
			},
		},
	}

	evidence := RecordingReleaseEvidence(
		recording,
	)

	context, ok := metadata.SelectReleaseContext(
		evidence,
	)
	if !ok {

		t.Fatal(
			"expected release context",
		)

	}

	if context.Title != "Unheard" {

		t.Errorf(
			"expected context title %q, got %q",
			"Unheard",
			context.Title,
		)

	}

	if context.Date != "2024-03-22" {

		t.Errorf(
			"expected context date %q, got %q",
			"2024-03-22",
			context.Date,
		)

	}

}
