package youtube

import (
	"context"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

type Catalog struct {
	urls []string
}

func NewCatalog(
	urls []string,
) *Catalog {

	return &Catalog{
		urls: urls,
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
		len(c.urls),
	)

	for _, raw := range c.urls {

		ref, ok := ParseURL(raw)
		if !ok {
			continue
		}

		item, ok := MediaItem(ref)
		if !ok {
			continue
		}

		items = append(
			items,
			item,
		)

	}

	return items, nil

}
