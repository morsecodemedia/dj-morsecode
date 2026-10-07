package source

import (
	"context"
	"testing"
)

type catalogStub struct {
	items []MediaItem
}

func (s catalogStub) Items(
	ctx context.Context,
) ([]MediaItem, error) {

	return s.items, nil
}

type discovererStub struct {
	results []DiscoveryResult
}

func (s discovererStub) Discover(
	ctx context.Context,
	request DiscoveryRequest,
) ([]DiscoveryResult, error) {

	return s.results, nil
}

func TestCatalogProviderContract(
	t *testing.T,
) {

	var provider CatalogProvider = catalogStub{
		items: []MediaItem{
			{
				Kind: MediaStation,
				Ref: ItemRef{
					Source: Source{
						Kind: KindM3U,
					},
					URI: "https://example.com/live",
				},
			},
		},
	}

	items, err := provider.Items(
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

}

func TestDiscovererContract(
	t *testing.T,
) {

	var provider Discoverer = discovererStub{
		results: []DiscoveryResult{
			{
				Item: MediaItem{
					Kind: MediaTrack,
					Ref: ItemRef{
						Source: Source{
							Kind: KindLastFM,
						},
						ID: "track-123",
					},
				},
				Score: 0.91,
			},
		},
	}

	results, err := provider.Discover(
		context.Background(),
		DiscoveryRequest{
			Kind:  DiscoverySearch,
			Query: "Stone Temple Pilots",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"expected 1 item, got %d",
			len(results),
		)
	}

}
