package source

import "context"

type CatalogProvider interface {
	Items(
		ctx context.Context,
	) ([]MediaItem, error)
}

type Discoverer interface {
	Discover(
		ctx context.Context,
		query string,
	) ([]MediaItem, error)
}
