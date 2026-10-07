package youtube

import (
	"context"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

type Catalog struct {
	entries []CatalogEntry
}

type CatalogEntry struct {
	Name string

	Artist string
	Title  string

	URL string
}

func NewCatalog(
	entries []CatalogEntry,
) *Catalog {

	return &Catalog{
		entries: entries,
	}

}

var _ source.CatalogProvider = (*Catalog)(nil)

func (c *Catalog) Items(
	ctx context.Context,
) ([]source.MediaItem, error) {

	_ = ctx

	items := make(
		[]source.MediaItem,
		0,
		len(c.entries),
	)

	for _, entry := range c.entries {

		ref, ok := ParseURL(
			entry.URL,
		)
		if !ok {
			continue
		}

		item, ok := MediaItem(ref)
		if !ok {
			continue
		}

		item.Ref.Name =
			entry.Name

		item.Artist =
			entry.Artist

		item.Title =
			entry.Title

		if item.Title == "" {
			item.Title =
				entry.Name
		}

		items = append(
			items,
			item,
		)

	}

	return items, nil

}
