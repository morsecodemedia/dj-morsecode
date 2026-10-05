package source

import "context"

type CatalogProvider interface {
	Items(
		ctx context.Context,
	) ([]MediaItem, error)
}

type DiscoveryKind string

const (
	DiscoverySearch        DiscoveryKind = "search"
	DiscoverySimilarTrack  DiscoveryKind = "similar-track"
	DiscoverySimilarArtist DiscoveryKind = "similar-artist"
	DiscoveryTag           DiscoveryKind = "tag"
)

type DiscoveryRequest struct {
	Kind DiscoveryKind

	Query string

	Artist string
	Title  string
	MBID   string

	Limit int
}

type Discoverer interface {
	Discover(
		ctx context.Context,
		request DiscoveryRequest,
	) ([]MediaItem, error)
}
