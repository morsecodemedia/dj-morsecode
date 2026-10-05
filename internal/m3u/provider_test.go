package m3u

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

func TestProviderItems(
	t *testing.T,
) {

	content := strings.NewReader(
		"#EXTM3U\n" +
			"#EXTINF:-1,Z100\n" +
			"https://example.com/z100\n" +
			"#EXTINF:-1,SKA World\n" +
			"https://example.com/ska\n",
	)

	provider := NewProvider(
		content,
		source.Source{
			Kind: source.KindM3U,
			URI:  "memory.m3u",
		},
	)

	items, err := provider.Items(
		context.Background(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 2 {

		t.Fatalf(
			"expected 2 items, got %d",
			len(items),
		)

	}

	if items[0].Kind !=
		source.MediaStation {

		t.Errorf(
			"expected station kind, got %q",
			items[0].Kind,
		)

	}

	if items[0].Title != "Z100" {

		t.Errorf(
			"expected Z100, got %q",
			items[0].Title,
		)

	}

	if items[0].Ref.Source.Kind !=
		source.KindM3U {

		t.Errorf(
			"expected M3U provenance, got %q",
			items[0].Ref.Source.Kind,
		)

	}

}

func TestProviderRealExtendedPlaylist(
	t *testing.T,
) {

	file, err := os.Open(
		"testdata/extended.m3u",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	provider := NewProvider(
		file,
		source.Source{
			Kind: source.KindM3U,
			URI:  "testdata/extended.m3u",
		},
	)

	items, err := provider.Items(
		context.Background(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) == 0 {
		t.Fatal(
			"expected catalog items",
		)
	}

	for i, item := range items {

		if !item.Valid() {

			t.Fatalf(
				"expected item %d to be valid",
				i,
			)

		}

		if item.Kind !=
			source.MediaStation {

			t.Fatalf(
				"expected item %d to be a station, got %q",
				i,
				item.Kind,
			)

		}

	}

}
