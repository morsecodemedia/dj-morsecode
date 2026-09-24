package musicbrainz

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSearchRecordings(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				if r.URL.Path != "/recording" {
					t.Errorf(
						"expected path %q, got %q",
						"/recording",
						r.URL.Path,
					)
				}

				if got := r.URL.Query().Get(
					"fmt",
				); got != "json" {

					t.Errorf(
						"expected fmt json, got %q",
						got,
					)

				}

				if got := r.URL.Query().Get(
					"limit",
				); got != "10" {

					t.Errorf(
						"expected limit 10, got %q",
						got,
					)

				}

				if got := r.Header.Get(
					"User-Agent",
				); got != "dj-morsecode/test" {

					t.Errorf(
						"expected user agent %q, got %q",
						"dj-morsecode/test",
						got,
					)

				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				_, _ = w.Write(
					[]byte(`{
						"recordings": [
							{
								"id": "026fa041-3917-4c73-9079-ed16e36f20f8",
								"score": 100,
								"title": "Blow Your Mind (Mwah)",
								"length": 178000,
								"artist-credit": [
									{
										"name": "Dua Lipa",
										"artist": {
											"name": "Dua Lipa"
										}
									}
								],
								"isrcs": [
									"GBAHT1600302"
								]
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

	recordings, err := client.SearchRecordings(
		context.Background(),
		"Dua Lipa",
		"Blow Your Mind (Mwah)",
	)
	if err != nil {
		t.Fatalf(
			"SearchRecordings returned error: %v",
			err,
		)
	}

	if len(recordings) != 1 {
		t.Fatalf(
			"expected one recording, got %d",
			len(recordings),
		)
	}

	recording := recordings[0]

	if recording.ID != "026fa041-3917-4c73-9079-ed16e36f20f8" {
		t.Errorf(
			"unexpected recording ID %q",
			recording.ID,
		)
	}

	if recording.Artist != "Dua Lipa" {
		t.Errorf(
			"expected artist %q, got %q",
			"Dua Lipa",
			recording.Artist,
		)
	}

	if recording.Title != "Blow Your Mind (Mwah)" {
		t.Errorf(
			"expected title %q, got %q",
			"Blow Your Mind (Mwah)",
			recording.Title,
		)
	}

	if recording.Score != 100 {
		t.Errorf(
			"expected score 100, got %d",
			recording.Score,
		)
	}

	if recording.Duration != 178*time.Second {
		t.Errorf(
			"expected duration %s, got %s",
			178*time.Second,
			recording.Duration,
		)
	}

}

func TestRecordingQuery(t *testing.T) {

	query := recordingQuery(
		`AC/DC`,
		`It's a Long Way "Home"`,
	)

	expected :=
		`artist:"AC/DC" AND recording:"It's a Long Way \"Home\""`

	if query != expected {
		t.Errorf(
			"expected query %q, got %q",
			expected,
			query,
		)
	}

}
