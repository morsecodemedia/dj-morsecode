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
		[]string{
			"https://youtu.be/video-one",
			"https://www.youtube.com/playlist?list=PL123",
			"https://example.com/not-youtube",
			"https://www.youtube.com/shorts/video-two",
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

}
