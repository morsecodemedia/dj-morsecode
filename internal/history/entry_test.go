package history

import (
	"testing"
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

func TestStationEntry(
	t *testing.T,
) {

	playedAt := time.Date(
		2026,
		time.October,
		7,
		22,
		0,
		0,
		0,
		time.UTC,
	)

	entry := Entry{
		Kind:      KindStation,
		StationID: "hardrockradiofm",
		PlayedAt:  playedAt,
	}

	if entry.Kind != KindStation {
		t.Errorf(
			"expected station kind, got %q",
			entry.Kind,
		)
	}

	if entry.StationID !=
		"hardrockradiofm" {

		t.Errorf(
			"unexpected station ID %q",
			entry.StationID,
		)
	}

	if !entry.PlayedAt.Equal(
		playedAt,
	) {

		t.Errorf(
			"unexpected played time %v",
			entry.PlayedAt,
		)

	}

}

func TestMediaEntryPreservesSourceItem(
	t *testing.T,
) {

	item := source.MediaItem{
		Kind: source.MediaVideo,

		Ref: source.ItemRef{
			Source: source.Source{
				Kind: source.KindYouTube,
			},

			ID:   "qORYO0atB6g",
			URI:  "https://www.youtube.com/watch?v=qORYO0atB6g",
			Name: "Beastie Boys - Intergalactic",
		},

		Artist: "Beastie Boys",
		Title:  "Intergalactic",
	}

	entry := Entry{
		Kind:  KindMedia,
		Media: item,
	}

	if entry.Kind != KindMedia {
		t.Errorf(
			"expected media kind, got %q",
			entry.Kind,
		)
	}

	if entry.Media.Ref.ID !=
		item.Ref.ID {

		t.Errorf(
			"unexpected media ID %q",
			entry.Media.Ref.ID,
		)

	}

	if entry.Media.Ref.Source.Kind !=
		source.KindYouTube {

		t.Errorf(
			"unexpected source kind %q",
			entry.Media.Ref.Source.Kind,
		)

	}

	if entry.Media.Artist !=
		"Beastie Boys" {

		t.Errorf(
			"unexpected artist %q",
			entry.Media.Artist,
		)

	}

	if entry.Media.Title !=
		"Intergalactic" {

		t.Errorf(
			"unexpected title %q",
			entry.Media.Title,
		)

	}

}
