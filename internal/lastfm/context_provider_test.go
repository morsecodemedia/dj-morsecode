package lastfm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

func TestContextProviderUsesArtistEvidenceWhenTrackTagsEmpty(
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

				switch r.URL.Query().Get(
					"method",
				) {

				case "track.getTopTags":

					_, _ = w.Write(
						[]byte(`{
							"toptags": {
								"tag": []
							}
						}`),
					)

				case "artist.getTopTags":

					_, _ = w.Write(
						[]byte(`{
							"toptags": {
								"tag": [
									{
										"name": "hard rock",
										"count": 100
									},
									{
										"name": "classic rock",
										"count": 61
									},
									{
										"name": "australian",
										"count": 3
									}
								]
							}
						}`),
					)

				}

			},
		),
	)
	defer server.Close()

	client := NewClient(
		"test-key",
	)

	client.baseURL = server.URL
	client.httpClient = server.Client()

	provider := NewContextProvider(
		client,
	)

	track := metadata.CanonicalTrack{
		Artist: "AC/DC",
		Title:  "Who Made Who",
	}

	result, ok, err := provider.Contextualize(
		context.Background(),
		track,
	)
	if err != nil {

		t.Fatalf(
			"Contextualize returned error: %v",
			err,
		)

	}

	if !ok {
		t.Fatal(
			"expected track context",
		)
	}

	if len(result.Tags) != 3 {

		t.Fatalf(
			"expected three context tags, got %d",
			len(result.Tags),
		)

	}

	// ResolveContextTags sorts by normalized name.
	if result.Tags[0].Name != "australian" {
		t.Errorf(
			"expected first tag %q, got %q",
			"australian",
			result.Tags[0].Name,
		)
	}

	if result.Tags[0].ArtistCount != 3 {
		t.Errorf(
			"expected artist count %d, got %d",
			3,
			result.Tags[0].ArtistCount,
		)
	}

}
