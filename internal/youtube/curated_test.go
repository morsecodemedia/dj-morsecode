package youtube

import (
	"context"
	"testing"
)

func TestCuratedCatalog(
	t *testing.T,
) {

	catalog := NewCatalog(
		Curated,
	)

	items, err := catalog.Items(
		context.Background(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) == 0 {
		t.Fatal(
			"expected curated YouTube items",
		)
	}

	for i, item := range items {

		if !item.Valid() {

			t.Fatalf(
				"expected item %d to be valid",
				i,
			)

		}

	}

}

func TestCatalogUsesNameAsTitleFallback(
	t *testing.T,
) {

	catalog := NewCatalog(
		[]CatalogEntry{
			{
				Name: "Interesting Video",
				URL:  "https://youtu.be/video-one",
			},
		},
	)

	items, err := catalog.Items(
		context.Background(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 1 {
		t.Fatalf(
			"expected 1 item, got %d",
			len(items),
		)
	}

	if items[0].Title !=
		"Interesting Video" {

		t.Errorf(
			"expected title fallback, got %q",
			items[0].Title,
		)

	}

}
