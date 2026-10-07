package youtube

import (
	"context"
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

func TestCatalogItems(
	t *testing.T,
) {

	catalog := NewCatalog(
		[]CatalogEntry{
			{
				Name: "Video One",
				URL:  "https://youtu.be/video-one",
			},
			{
				Name: "Playlist",
				URL:  "https://www.youtube.com/playlist?list=PL123",
			},
			{
				Name: "Not YouTube",
				URL:  "https://example.com/not-youtube",
			},
			{
				Name: "Video Two",
				URL:  "https://www.youtube.com/shorts/video-two",
			},
		},
	)

	items, err := catalog.Items(
		context.Background(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 3 {
		t.Fatalf(
			"expected 3 items, got %d",
			len(items),
		)
	}

	if items[0].Kind != source.MediaVideo {
		t.Errorf(
			"expected first item to be video, got %q",
			items[0].Kind,
		)
	}

	if items[1].Kind != source.MediaPlaylist {
		t.Errorf(
			"expected second item to be playlist, got %q",
			items[1].Kind,
		)
	}

	if items[2].Kind != source.MediaVideo {
		t.Errorf(
			"expected third item to be video, got %q",
			items[2].Kind,
		)
	}

	if items[0].Title != "Video One" {

		t.Errorf(
			"expected curated title, got %q",
			items[0].Title,
		)

	}

}
