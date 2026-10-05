package lastfm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

func TestDiscovererSimilarTrack(
	t *testing.T,
) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				_, _ = w.Write(
					[]byte(`{
						"similartracks": {
							"track": [
								{
									"name": "Would?",
									"mbid": "would-recording",
									"match": "0.9234",
									"url": "https://www.last.fm/example/would",
									"artist": {
										"name": "Alice in Chains"
									}
								}
							]
						}
					}`),
				)

			},
		),
	)
	defer server.Close()

	client := NewClient(
		"test-key",
	)

	client.baseURL =
		server.URL

	client.httpClient =
		server.Client()

	discoverer := NewDiscoverer(
		client,
	)

	results, err := discoverer.Discover(
		context.Background(),
		source.DiscoveryRequest{
			Kind: source.DiscoverySimilarTrack,

			Artist: "Stone Temple Pilots",

			Title: "Interstate Love Song",

			Limit: 10,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {

		t.Fatalf(
			"expected 1 discovery result, got %d",
			len(results),
		)

	}

	result := results[0]

	if result.Item.Kind !=
		source.MediaTrack {

		t.Errorf(
			"expected track, got %q",
			result.Item.Kind,
		)

	}

	if result.Item.Artist !=
		"Alice in Chains" {

		t.Errorf(
			"unexpected artist %q",
			result.Item.Artist,
		)

	}

	if result.Item.Title !=
		"Would?" {

		t.Errorf(
			"unexpected title %q",
			result.Item.Title,
		)

	}

	if result.Item.Ref.Source.Kind !=
		source.KindLastFM {

		t.Errorf(
			"expected Last.fm provenance, got %q",
			result.Item.Ref.Source.Kind,
		)

	}

	if result.Score != 0.9234 {

		t.Errorf(
			"expected score 0.9234, got %f",
			result.Score,
		)

	}

}

func TestDiscovererRejectsUnsupportedDiscovery(
	t *testing.T,
) {

	discoverer := NewDiscoverer(
		NewClient("test-key"),
	)

	_, err := discoverer.Discover(
		context.Background(),
		source.DiscoveryRequest{
			Kind:  source.DiscoverySearch,
			Query: "Stone Temple Pilots",
		},
	)

	if err == nil {
		t.Fatal(
			"expected unsupported discovery error",
		)
	}

}

func TestDiscovererRequiresTrackIdentity(
	t *testing.T,
) {

	discoverer := NewDiscoverer(
		NewClient("test-key"),
	)

	_, err := discoverer.Discover(
		context.Background(),
		source.DiscoveryRequest{
			Kind:   source.DiscoverySimilarTrack,
			Artist: "Stone Temple Pilots",
		},
	)

	if err == nil {
		t.Fatal(
			"expected missing track identity error",
		)
	}

}
