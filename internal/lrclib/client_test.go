package lrclib

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearch(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				if r.URL.Path != "/api/search" {
					t.Fatalf(
						"expected /api/search, got %s",
						r.URL.Path,
					)
				}

				if got := r.URL.Query().Get("artist_name"); got != "Aylex" {
					t.Errorf(
						"expected artist_name Aylex, got %q",
						got,
					)
				}

				if got := r.URL.Query().Get("track_name"); got != "Blast" {
					t.Errorf(
						"expected track_name Blast, got %q",
						got,
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				_, _ = w.Write(
					[]byte(`[
						{
							"id": 123,
							"trackName": "Blast",
							"artistName": "Aylex",
							"albumName": "Test Album",
							"duration": 71.5,
							"instrumental": false,
							"plainLyrics": "Hello world",
							"syncedLyrics": "[00:01.00]Hello world",
							"lyricsfile": ""
						}
					]`),
				)

			},
		),
	)
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.HTTPClient = server.Client()

	results, err := client.Search(
		"Aylex",
		"Blast",
	)
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(results),
		)
	}

	result := results[0]

	if result.TrackName != "Blast" {
		t.Errorf(
			"expected track name Blast, got %q",
			result.TrackName,
		)
	}

	if result.ArtistName != "Aylex" {
		t.Errorf(
			"expected artist name Aylex, got %q",
			result.ArtistName,
		)
	}

	if result.SyncedLyrics != "[00:01.00]Hello world" {
		t.Errorf(
			"unexpected synced lyrics %q",
			result.SyncedLyrics,
		)
	}

}
