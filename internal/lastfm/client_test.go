package lastfm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrackTopTags(
	t *testing.T,
) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				query := r.URL.Query()

				if got := query.Get(
					"method",
				); got != "track.getTopTags" {

					t.Errorf(
						"unexpected method %q",
						got,
					)

				}

				if got := query.Get(
					"api_key",
				); got != "test-key" {

					t.Errorf(
						"unexpected API key %q",
						got,
					)

				}

				if got := query.Get(
					"artist",
				); got != "Hozier" {

					t.Errorf(
						"unexpected artist %q",
						got,
					)

				}

				if got := query.Get(
					"track",
				); got != "Too Sweet" {

					t.Errorf(
						"unexpected track %q",
						got,
					)

				}

				if got := query.Get(
					"autocorrect",
				); got != "1" {

					t.Errorf(
						"unexpected autocorrect %q",
						got,
					)

				}

				if got := query.Get(
					"format",
				); got != "json" {

					t.Errorf(
						"unexpected format %q",
						got,
					)

				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				_, _ = w.Write(
					[]byte(`{
						"toptags": {
							"tag": [
								{
									"name": "indie",
									"count": 100
								},
								{
									"name": "soul",
									"count": 72
								},
								{
									"name": "2024",
									"count": 19
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

	client.baseURL = server.URL
	client.httpClient = server.Client()

	tags, err := client.TrackTopTags(
		context.Background(),
		"Hozier",
		"Too Sweet",
	)
	if err != nil {

		t.Fatalf(
			"TrackTopTags returned error: %v",
			err,
		)

	}

	if len(tags) != 3 {

		t.Fatalf(
			"expected three tags, got %d",
			len(tags),
		)

	}

	if tags[0].Name != "indie" {

		t.Errorf(
			"expected first tag %q, got %q",
			"indie",
			tags[0].Name,
		)

	}

	if tags[0].Count != 100 {

		t.Errorf(
			"expected first count %d, got %d",
			100,
			tags[0].Count,
		)

	}

}

func TestTrackTopTagsRejectsHTTPError(
	t *testing.T,
) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				http.Error(
					w,
					"nope",
					http.StatusInternalServerError,
				)

			},
		),
	)
	defer server.Close()

	client := NewClient(
		"test-key",
	)

	client.baseURL = server.URL
	client.httpClient = server.Client()

	_, err := client.TrackTopTags(
		context.Background(),
		"Hozier",
		"Too Sweet",
	)

	if err == nil {
		t.Fatal(
			"expected HTTP error",
		)
	}

}

func TestTrackInfo(
	t *testing.T,
) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				query := r.URL.Query()

				if got := query.Get(
					"method",
				); got != "track.getInfo" {

					t.Errorf(
						"unexpected method %q",
						got,
					)

				}

				if got := query.Get(
					"artist",
				); got != "Tycho" {

					t.Errorf(
						"unexpected artist %q",
						got,
					)

				}

				if got := query.Get(
					"track",
				); got != "A Walk" {

					t.Errorf(
						"unexpected track %q",
						got,
					)

				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				_, _ = w.Write(
					[]byte(`{
						"track": {
							"name": "A Walk",
							"mbid": "recording-id",
							"url": "https://www.last.fm/music/Tycho/_/A+Walk",
							"listeners": "123456",
							"playcount": "789012",
							"artist": {
								"name": "Tycho"
							},
							"toptags": {
								"tag": [
									{
										"name": "ambient"
									},
									{
										"name": "electronic"
									}
								]
							}
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

	client.baseURL = server.URL
	client.httpClient = server.Client()

	info, err := client.TrackInfo(
		context.Background(),
		"Tycho",
		"A Walk",
	)
	if err != nil {

		t.Fatalf(
			"TrackInfo returned error: %v",
			err,
		)

	}

	if info.Artist != "Tycho" {
		t.Errorf(
			"expected artist %q, got %q",
			"Tycho",
			info.Artist,
		)
	}

	if info.Title != "A Walk" {
		t.Errorf(
			"expected title %q, got %q",
			"A Walk",
			info.Title,
		)
	}

	if info.MBID != "recording-id" {
		t.Errorf(
			"expected MBID %q, got %q",
			"recording-id",
			info.MBID,
		)
	}

	if info.Listeners != 123456 {
		t.Errorf(
			"expected listeners %d, got %d",
			123456,
			info.Listeners,
		)
	}

	if info.Playcount != 789012 {
		t.Errorf(
			"expected playcount %d, got %d",
			789012,
			info.Playcount,
		)
	}

	if len(info.Tags) != 2 {
		t.Fatalf(
			"expected two tags, got %d",
			len(info.Tags),
		)
	}

}

func TestArtistTopTags(
	t *testing.T,
) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				query := r.URL.Query()

				if got := query.Get(
					"method",
				); got != "artist.getTopTags" {

					t.Errorf(
						"unexpected method %q",
						got,
					)

				}

				if got := query.Get(
					"artist",
				); got != "Tycho" {

					t.Errorf(
						"unexpected artist %q",
						got,
					)

				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				_, _ = w.Write(
					[]byte(`{
						"toptags": {
							"tag": [
								{
									"name": "electronic",
									"count": 100
								},
								{
									"name": "ambient",
									"count": 91
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

	client.baseURL = server.URL
	client.httpClient = server.Client()

	tags, err := client.ArtistTopTags(
		context.Background(),
		"Tycho",
	)
	if err != nil {

		t.Fatalf(
			"ArtistTopTags returned error: %v",
			err,
		)

	}

	if len(tags) != 2 {

		t.Fatalf(
			"expected two tags, got %d",
			len(tags),
		)

	}

	if tags[0].Name != "electronic" ||
		tags[0].Count != 100 {

		t.Errorf(
			"unexpected first tag: %+v",
			tags[0],
		)

	}

}
