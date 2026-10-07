package m3u

import (
	"context"
	"io"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

type Provider struct {
	reader io.Reader
	origin source.Source
}

var _ source.CatalogProvider = (*Provider)(nil)

func NewProvider(
	reader io.Reader,
	origin source.Source,
) *Provider {

	return &Provider{
		reader: reader,
		origin: origin,
	}

}

func (p *Provider) Items(
	ctx context.Context,
) ([]source.MediaItem, error) {

	_ = ctx

	entries, err := Parse(
		p.reader,
		p.origin,
	)
	if err != nil {
		return nil, err
	}

	items := make(
		[]source.MediaItem,
		0,
		len(entries),
	)

	for _, entry := range entries {

		item, ok :=
			StationMediaItem(
				entry,
			)

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
