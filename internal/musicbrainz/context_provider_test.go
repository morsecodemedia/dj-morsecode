package musicbrainz

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

func TestContextProviderSkipsTrackWithoutMusicBrainzID(
	t *testing.T,
) {

	requested := false

	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				requested = true
			},
		),
	)
	defer server.Close()

	client := NewClient(
		"dj-morsecode/test",
	)

	client.baseURL = server.URL
	client.httpClient = server.Client()

	provider := NewContextProvider(
		client,
	)

	track := metadata.CanonicalTrack{
		Artist: "Hozier",
		Title:  "Too Sweet",
	}

	_, ok, err := provider.Contextualize(
		context.Background(),
		track,
	)
	if err != nil {

		t.Fatalf(
			"Contextualize returned error: %v",
			err,
		)

	}

	if ok {

		t.Fatal(
			"expected no context without MusicBrainz ID",
		)

	}

	if requested {

		t.Fatal(
			"expected no MusicBrainz request",
		)

	}

}

func TestContextProviderSelectsRelease(
	t *testing.T,
) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				if r.URL.Path !=
					"/recording/272be1e2-9599-4b8d-b872-5e879881d103" {

					t.Errorf(
						"unexpected path %q",
						r.URL.Path,
					)

				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				_, _ = w.Write(
					[]byte(`{
						"id": "272be1e2-9599-4b8d-b872-5e879881d103",
						"title": "Too Sweet",
						"length": 251000,
						"artist-credit": [
							{
								"name": "Hozier",
								"artist": {
									"name": "Hozier"
								}
							}
						],
						"releases": [
							{
								"id": "unheard-release",
								"title": "Unheard",
								"date": "2024-03-22",
								"release-group": {
									"id": "unheard",
									"title": "Unheard",
									"first-release-date": "2024-03-22",
									"primary-type": "EP",
									"secondary-types": []
								}
							},
							{
								"id": "unreal-release",
								"title": "Unreal Unearth",
								"date": "2024-03-22",
								"release-group": {
									"id": "unreal-unearth",
									"title": "Unreal Unearth",
									"first-release-date": "2023-08-18",
									"primary-type": "Album",
									"secondary-types": []
								}
							},
							{
								"id": "bravo-release",
								"title": "Bravo Hits 125",
								"date": "2024-04-26",
								"release-group": {
									"id": "bravo",
									"title": "Bravo Hits 125",
									"first-release-date": "2024-04-26",
									"primary-type": "Album",
									"secondary-types": [
										"Compilation"
									]
								}
							}
						]
					}`),
				)

			},
		),
	)
	defer server.Close()

	client := NewClient(
		"dj-morsecode/test",
	)

	client.baseURL = server.URL
	client.httpClient = server.Client()

	provider := NewContextProvider(
		client,
	)

	track := metadata.CanonicalTrack{
		Artist: "Hozier",
		Title:  "Too Sweet",

		Identifiers: []metadata.Identifier{
			{
				Scheme: metadata.IdentifierMusicBrainz,
				Value:  "272be1e2-9599-4b8d-b872-5e879881d103",
			},
		},
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

	if result.Release.Title != "Unheard" {

		t.Errorf(
			"expected release title %q, got %q",
			"Unheard",
			result.Release.Title,
		)

	}

	if result.Release.Date != "2024-03-22" {

		t.Errorf(
			"expected release date %q, got %q",
			"2024-03-22",
			result.Release.Date,
		)

	}

	if len(result.Genres) != 0 {

		t.Errorf(
			"expected genres to remain empty, got %v",
			result.Genres,
		)

	}

}

func TestContextProviderRejectsAmbiguousReleaseContext(
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
						"id": "recording-id",
						"title": "Example",
						"releases": [
							{
								"title": "Album",
								"date": "2000",
								"release-group": {
									"id": "album",
									"title": "Album",
									"first-release-date": "2000",
									"primary-type": "Album",
									"secondary-types": []
								}
							},
							{
								"title": "Single",
								"date": "2000",
								"release-group": {
									"id": "single",
									"title": "Single",
									"first-release-date": "2000",
									"primary-type": "Single",
									"secondary-types": []
								}
							}
						]
					}`),
				)

			},
		),
	)
	defer server.Close()

	client := NewClient(
		"dj-morsecode/test",
	)

	client.baseURL = server.URL
	client.httpClient = server.Client()

	provider := NewContextProvider(
		client,
	)

	track := metadata.CanonicalTrack{
		Artist: "Artist",
		Title:  "Example",
		Identifiers: []metadata.Identifier{
			{
				Scheme: metadata.IdentifierMusicBrainz,
				Value:  "recording-id",
			},
		},
	}

	_, ok, err := provider.Contextualize(
		context.Background(),
		track,
	)
	if err != nil {

		t.Fatalf(
			"Contextualize returned error: %v",
			err,
		)

	}

	if ok {

		t.Fatal(
			"expected ambiguous release context to remain unresolved",
		)

	}

}
