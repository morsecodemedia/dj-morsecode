package lrclib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://lrclib.net"
const durationTolerance = 5 * time.Second

type Result struct {
	ID           int     `json:"id"`
	TrackName    string  `json:"trackName"`
	ArtistName   string  `json:"artistName"`
	AlbumName    string  `json:"albumName"`
	Duration     float64 `json:"duration"`
	Instrumental bool    `json:"instrumental"`
	PlainLyrics  string  `json:"plainLyrics"`
	SyncedLyrics string  `json:"syncedLyrics"`
	Lyricsfile   string  `json:"lyricsfile"`
}

func BestMatch(
	results []Result,
	duration time.Duration,
) (Result, bool) {

	var best Result
	var bestDelta time.Duration
	found := false

	for _, result := range results {

		if result.SyncedLyrics == "" {
			continue
		}

		resultDuration := time.Duration(
			result.Duration * float64(time.Second),
		)

		delta := resultDuration - duration

		if delta < 0 {
			delta = -delta
		}

		if delta > durationTolerance {
			continue
		}

		if !found || delta < bestDelta {

			best = result
			bestDelta = delta
			found = true

		}

	}

	return best, found

}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {

	return &Client{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
	}

}

func (c *Client) Search(artist, title string) ([]Result, error) {

	endpoint, err := url.Parse(c.BaseURL + "/api/search")
	if err != nil {
		return nil, fmt.Errorf("parse LRCLIB URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("track_name", title)

	if artist != "" {
		query.Set("artist_name", artist)
	}

	endpoint.RawQuery = query.Encode()

	request, err := http.NewRequest(
		http.MethodGet,
		endpoint.String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create LRCLIB request: %w", err)
	}

	request.Header.Set(
		"User-Agent",
		"dj-morsecode/0.0.1 (https://github.com/morsecodemedia/dj-morsecode)",
	)

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("search LRCLIB: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"LRCLIB returned status %s",
			response.Status,
		)
	}

	var results []Result

	if err := json.NewDecoder(response.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("decode LRCLIB response: %w", err)
	}

	return results, nil

}
