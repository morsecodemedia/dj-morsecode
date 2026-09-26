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
