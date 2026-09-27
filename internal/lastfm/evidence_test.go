package lastfm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

func TestTagEvidence(t *testing.T) {

	tags := []Tag{
		{
			Name:  "hard rock",
			Count: 100,
		},
		{
			Name:  "australian",
			Count: 3,
		},
	}

	evidence := TagEvidence(
		tags,
		metadata.TagScopeArtist,
	)

	if len(evidence) != 2 {

		t.Fatalf(
			"expected two evidence items, got %d",
			len(evidence),
		)

	}

	first := evidence[0]

	if first.Name != "hard rock" {
		t.Errorf(
			"expected name %q, got %q",
			"hard rock",
			first.Name,
		)
	}

	if first.Count != 100 {
		t.Errorf(
			"expected count %d, got %d",
			100,
			first.Count,
		)
	}

	if first.Scope != metadata.TagScopeArtist {
		t.Errorf(
			"expected artist scope, got %q",
			first.Scope,
		)
	}

	if first.Provider != providerName {
		t.Errorf(
			"expected provider %q, got %q",
			providerName,
			first.Provider,
		)
	}

}

func TestTagEvidencePreservesTrackScope(
	t *testing.T,
) {

	evidence := TagEvidence(
		[]Tag{
			{
				Name:  "pop soul",
				Count: 100,
			},
		},
		metadata.TagScopeTrack,
	)

	if len(evidence) != 1 {
		t.Fatalf(
			"expected one evidence item, got %d",
			len(evidence),
		)
	}

	if evidence[0].Scope !=
		metadata.TagScopeTrack {

		t.Errorf(
			"expected track scope, got %q",
			evidence[0].Scope,
		)

	}

}

func TestTrackTagEvidenceCombinesScopes(
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
								"tag": [
									{
										"name": "pop soul",
										"count": 100
									}
								]
							}
						}`),
					)

				case "artist.getTopTags":

					_, _ = w.Write(
						[]byte(`{
							"toptags": {
								"tag": [
									{
										"name": "blues",
										"count": 100
									}
								]
							}
						}`),
					)

				default:

					http.Error(
						w,
						"unexpected method",
						http.StatusBadRequest,
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

	result := client.TrackTagEvidence(
		context.Background(),
		"Hozier",
		"Too Sweet",
	)

	if result.Failed() {
		t.Fatal(
			"expected evidence request to succeed",
		)
	}

	if !result.HasEvidence() {
		t.Fatal(
			"expected tag evidence",
		)
	}

	if len(result.Evidence) != 2 {
		t.Fatalf(
			"expected two evidence items, got %d",
			len(result.Evidence),
		)
	}

	if result.Evidence[0].Scope !=
		metadata.TagScopeTrack {

		t.Errorf(
			"expected first evidence to have track scope",
		)

	}

	if result.Evidence[1].Scope !=
		metadata.TagScopeArtist {

		t.Errorf(
			"expected second evidence to have artist scope",
		)

	}

}
