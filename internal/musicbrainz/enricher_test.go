package musicbrainz

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

func TestEnricher(t *testing.T) {

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
						"recordings": [
							{
								"id": "standard-recording",
								"score": 100,
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
								"isrcs": [
									"USSM12401865"
								]
							},
							{
								"id": "atmos-recording",
								"score": 99,
								"title": "Too Sweet",
								"length": 251424,
								"disambiguation": "Dolby Atmos mix",
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

	enricher := NewEnricher(
		client,
	)

	observed := metadata.PlaybackItem{
		Type:     metadata.PlaybackTrack,
		Artist:   "Hozier",
		Title:    "Too Sweet",
		Duration: 251 * time.Second,
	}

	match, status, err := enricher.Enrich(
		context.Background(),
		observed,
	)
	if err != nil {

		t.Fatalf(
			"Enrich returned error: %v",
			err,
		)

	}

	if status != metadata.MatchAccepted {

		t.Fatalf(
			"expected accepted match, got %v",
			status,
		)

	}

	if match.Track.Artist != "Hozier" {
		t.Errorf(
			"expected artist %q, got %q",
			"Hozier",
			match.Track.Artist,
		)
	}

	if match.Track.Title != "Too Sweet" {
		t.Errorf(
			"expected title %q, got %q",
			"Too Sweet",
			match.Track.Title,
		)
	}

	if len(match.Track.Identifiers) == 0 {
		t.Fatal(
			"expected canonical identifiers",
		)
	}

	if match.Track.Identifiers[0].Value !=
		"standard-recording" {

		t.Errorf(
			"expected standard recording, got %q",
			match.Track.Identifiers[0].Value,
		)

	}

}
func TestEnricherSkipsNonTrack(
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

	enricher := NewEnricher(
		client,
	)

	observed := metadata.PlaybackItem{
		Type:  metadata.PlaybackStationID,
		Title: "Z100",
	}

	_, status, err := enricher.Enrich(
		context.Background(),
		observed,
	)
	if err != nil {

		t.Fatalf(
			"Enrich returned error: %v",
			err,
		)

	}

	if status != metadata.MatchNone {

		t.Fatalf(
			"expected no match, got %v",
			status,
		)

	}

	if requested {

		t.Fatal(
			"expected non-track enrichment not to make request",
		)

	}

}
