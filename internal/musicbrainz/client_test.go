package musicbrainz

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
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

func TestSearchRecordingsThrottlesRequests(
	t *testing.T,
) {

	var mu sync.Mutex
	var requests []time.Time

	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				mu.Lock()
				requests = append(
					requests,
					time.Now(),
				)
				mu.Unlock()

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				_, _ = w.Write(
					[]byte(`{
						"recordings": []
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
	client.requestInterval = 25 * time.Millisecond

	ctx := context.Background()

	if _, err := client.SearchRecordings(
		ctx,
		"Hozier",
		"Too Sweet",
	); err != nil {

		t.Fatalf(
			"first search returned error: %v",
			err,
		)

	}

	if _, err := client.SearchRecordings(
		ctx,
		"Hozier",
		"Too Sweet",
	); err != nil {

		t.Fatalf(
			"second search returned error: %v",
			err,
		)

	}

	mu.Lock()
	defer mu.Unlock()

	if len(requests) != 2 {
		t.Fatalf(
			"expected two requests, got %d",
			len(requests),
		)
	}

	elapsed := requests[1].Sub(
		requests[0],
	)

	if elapsed < client.requestInterval {

		t.Errorf(
			"expected requests at least %s apart, got %s",
			client.requestInterval,
			elapsed,
		)

	}

}
func TestSearchRecordingsThrottleRespectsContext(
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
						"recordings": []
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
	client.requestInterval = time.Second

	if _, err := client.SearchRecordings(
		context.Background(),
		"Hozier",
		"Too Sweet",
	); err != nil {

		t.Fatalf(
			"first search returned error: %v",
			err,
		)

	}

	ctx, cancel := context.WithCancel(
		context.Background(),
	)
	cancel()

	_, err := client.SearchRecordings(
		ctx,
		"Hozier",
		"Too Sweet",
	)

	if err == nil {
		t.Fatal(
			"expected canceled search to return error",
		)
	}

}
func TestLookupRecording(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				if r.URL.Path !=
					"/recording/abc45ef4-a3a9-42bb-a14c-09c77f26936d" {

					t.Errorf(
						"unexpected path %q",
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
					"inc",
				); got != "artist-credits+isrcs+releases" {

					t.Errorf(
						"unexpected inc %q",
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
						"id": "abc45ef4-a3a9-42bb-a14c-09c77f26936d",
						"title": "Too Sweet",
						"length": 251424,
						"disambiguation": "",
						"artist-credit": [
							{
								"name": "Hozier",
								"artist": {
									"name": "Hozier"
								}
							}
						],
						"isrcs": [
							"IEACJ2400038"
						],
						"releases": [
							{
								"id": "release-id",
								"title": "Unheard",
								"date": "2024-03-22",
								"country": "XW"
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

	recording, err := client.LookupRecording(
		context.Background(),
		"abc45ef4-a3a9-42bb-a14c-09c77f26936d",
	)
	if err != nil {

		t.Fatalf(
			"LookupRecording returned error: %v",
			err,
		)

	}

	if recording.ID !=
		"abc45ef4-a3a9-42bb-a14c-09c77f26936d" {

		t.Errorf(
			"unexpected recording ID %q",
			recording.ID,
		)

	}

	if recording.Artist != "Hozier" {

		t.Errorf(
			"expected artist %q, got %q",
			"Hozier",
			recording.Artist,
		)

	}

	if recording.Title != "Too Sweet" {

		t.Errorf(
			"expected title %q, got %q",
			"Too Sweet",
			recording.Title,
		)

	}

	if recording.Duration !=
		251424*time.Millisecond {

		t.Errorf(
			"expected duration %s, got %s",
			251424*time.Millisecond,
			recording.Duration,
		)

	}

	if len(recording.ISRCs) != 1 {

		t.Fatalf(
			"expected one ISRC, got %d",
			len(recording.ISRCs),
		)

	}

	if recording.ISRCs[0] != "IEACJ2400038" {

		t.Errorf(
			"expected ISRC %q, got %q",
			"IEACJ2400038",
			recording.ISRCs[0],
		)

	}

	if len(recording.Releases) != 1 {

		t.Fatalf(
			"expected one release, got %d",
			len(recording.Releases),
		)

	}

	release := recording.Releases[0]

	if release.Title != "Unheard" {

		t.Errorf(
			"expected release title %q, got %q",
			"Unheard",
			release.Title,
		)

	}

	if release.Date != "2024-03-22" {

		t.Errorf(
			"expected release date %q, got %q",
			"2024-03-22",
			release.Date,
		)

	}

	if release.Country != "XW" {

		t.Errorf(
			"expected release country %q, got %q",
			"XW",
			release.Country,
		)

	}

}
