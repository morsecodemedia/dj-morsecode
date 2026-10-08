package history

import (
	"testing"
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/observation"
	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

func historyTime(
	minute int,
) time.Time {

	return time.Date(
		2026,
		time.October,
		7,
		22,
		minute,
		0,
		0,
		time.UTC,
	)

}

func mediaObservation(
	id string,
	artist string,
	title string,
	minute int,
) observation.MediaObservation {

	return observation.MediaObservation{
		Item: source.MediaItem{
			Kind: source.MediaVideo,

			Ref: source.ItemRef{
				Source: source.Source{
					Kind: source.KindYouTube,
				},

				ID: id,
			},

			Artist: artist,
			Title:  title,
		},

		ObservedAt: historyTime(
			minute,
		),
	}

}

func TestDeriveIncludesOnlyConfirmedStations(
	t *testing.T,
) {

	stations := []observation.StationObservation{
		{
			Kind: observation.StationTuneRequested,

			StationID: "z100",
			ObservedAt: historyTime(
				0,
			),
		},
		{
			Kind: observation.StationTuneConfirmed,

			StationID: "z100",
			ObservedAt: historyTime(
				1,
			),
		},
		{
			Kind: observation.StationTuneFailed,

			StationID: "wxpn",
			ObservedAt: historyTime(
				2,
			),
		},
	}

	entries := Derive(
		stations,
		nil,
	)

	if len(entries) != 1 {
		t.Fatalf(
			"expected one history entry, got %d",
			len(entries),
		)
	}

	if entries[0].Kind !=
		KindStation {

		t.Errorf(
			"unexpected kind %q",
			entries[0].Kind,
		)

	}

	if entries[0].StationID !=
		"z100" {

		t.Errorf(
			"unexpected station %q",
			entries[0].StationID,
		)

	}

}

func TestDeriveIncludesMedia(
	t *testing.T,
) {

	media := []observation.MediaObservation{
		mediaObservation(
			"a",
			"Beastie Boys",
			"Intergalactic",
			0,
		),
		mediaObservation(
			"b",
			"A Tribe Called Quest",
			"Can I Kick It?",
			1,
		),
	}

	entries := Derive(
		nil,
		media,
	)

	if len(entries) != 2 {
		t.Fatalf(
			"expected two entries, got %d",
			len(entries),
		)
	}

	if entries[0].Media.Ref.ID != "a" {
		t.Errorf(
			"unexpected first media %q",
			entries[0].Media.Ref.ID,
		)
	}

	if entries[1].Media.Ref.ID != "b" {
		t.Errorf(
			"unexpected second media %q",
			entries[1].Media.Ref.ID,
		)
	}

}

func TestDeriveOrdersMixedHistoryChronologically(
	t *testing.T,
) {

	stations :=
		[]observation.StationObservation{
			{
				Kind: observation.StationTuneConfirmed,

				StationID: "hardrockradiofm",

				ObservedAt: historyTime(2),
			},
		}

	media :=
		[]observation.MediaObservation{
			mediaObservation(
				"b",
				"A Tribe Called Quest",
				"Can I Kick It?",
				3,
			),
			mediaObservation(
				"a",
				"Beastie Boys",
				"Intergalactic",
				1,
			),
		}

	entries := Derive(
		stations,
		media,
	)

	if len(entries) != 3 {
		t.Fatalf(
			"expected three entries, got %d",
			len(entries),
		)
	}

	if entries[0].Media.Ref.ID != "a" {
		t.Errorf(
			"unexpected first entry",
		)
	}

	if entries[1].StationID !=
		"hardrockradiofm" {

		t.Errorf(
			"unexpected second entry",
		)

	}

	if entries[2].Media.Ref.ID != "b" {
		t.Errorf(
			"unexpected third entry",
		)
	}

}

func TestDerivePreservesStationRevisits(
	t *testing.T,
) {

	stations :=
		[]observation.StationObservation{
			{
				Kind: observation.StationTuneConfirmed,

				StationID:  "z100",
				ObservedAt: historyTime(0),
			},
			{
				Kind: observation.StationTuneConfirmed,

				StationID: "hardrockradiofm",

				ObservedAt: historyTime(1),
			},
			{
				Kind: observation.StationTuneConfirmed,

				StationID:  "z100",
				ObservedAt: historyTime(2),
			},
		}

	entries := Derive(
		stations,
		nil,
	)

	if len(entries) != 3 {
		t.Fatalf(
			"expected three entries, got %d",
			len(entries),
		)
	}

	expected :=
		[]string{
			"z100",
			"hardrockradiofm",
			"z100",
		}

	for i, stationID := range expected {

		if entries[i].StationID !=
			stationID {

			t.Errorf(
				"entry %d: expected %q, got %q",
				i,
				stationID,
				entries[i].StationID,
			)

		}

	}

}

func TestDerivePreservesMediaRevisitAcrossStation(
	t *testing.T,
) {

	stations :=
		[]observation.StationObservation{
			{
				Kind: observation.StationTuneConfirmed,

				StationID: "hardrockradiofm",

				ObservedAt: historyTime(1),
			},
		}

	media :=
		[]observation.MediaObservation{
			mediaObservation(
				"a",
				"Beastie Boys",
				"Intergalactic",
				0,
			),
			mediaObservation(
				"a",
				"Beastie Boys",
				"Intergalactic",
				2,
			),
		}

	entries := Derive(
		stations,
		media,
	)

	if len(entries) != 3 {
		t.Fatalf(
			"expected three entries, got %d",
			len(entries),
		)
	}

	if entries[0].Media.Ref.ID != "a" ||
		entries[1].StationID !=
			"hardrockradiofm" ||
		entries[2].Media.Ref.ID != "a" {

		t.Fatal(
			"expected media → station → media chronology",
		)

	}

}
