package musicbrainz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const defaultBaseURL = "https://musicbrainz.org/ws/2"

const defaultLimit = 10
const minimumRequestInterval = time.Second

type Client struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string

	requestInterval time.Duration

	mu          sync.Mutex
	lastRequest time.Time
}

func NewClient(
	userAgent string,
) *Client {

	return &Client{
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		userAgent:       userAgent,
		requestInterval: minimumRequestInterval,
	}

}

func (c *Client) waitForRequest(
	ctx context.Context,
) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lastRequest.IsZero() {

		c.lastRequest = time.Now()
		return nil

	}

	wait := c.requestInterval -
		time.Since(c.lastRequest)

	if wait > 0 {

		timer := time.NewTimer(
			wait,
		)
		defer timer.Stop()

		select {

		case <-ctx.Done():
			return ctx.Err()

		case <-timer.C:

		}

	}

	c.lastRequest = time.Now()

	return nil

}

type searchResponse struct {
	Recordings []recordingResponse `json:"recordings"`
}

type releaseResponse struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Date    string `json:"date"`
	Country string `json:"country"`
}

type recordingResponse struct {
	ID     string `json:"id"`
	Score  int    `json:"score"`
	Title  string `json:"title"`
	Length int64  `json:"length"`

	ArtistCredit []artistCreditResponse `json:"artist-credit"`

	ISRCs []string `json:"isrcs"`

	Disambiguation string `json:"disambiguation"`

	Releases []releaseResponse `json:"releases"`
}

type artistCreditResponse struct {
	Name       string `json:"name"`
	JoinPhrase string `json:"joinphrase"`

	Artist struct {
		Name string `json:"name"`
	} `json:"artist"`
}

func (c *Client) SearchRecordings(
	ctx context.Context,
	artist string,
	title string,
) ([]Recording, error) {

	query := recordingQuery(
		artist,
		title,
	)

	endpoint, err := url.Parse(
		c.baseURL + "/recording",
	)
	if err != nil {
		return nil, err
	}

	values := endpoint.Query()
	values.Set("query", query)
	values.Set("fmt", "json")
	values.Set(
		"limit",
		strconv.Itoa(defaultLimit),
	)

	endpoint.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set(
		"User-Agent",
		c.userAgent,
	)

	if err := c.waitForRequest(
		ctx,
	); err != nil {

		return nil, err

	}

	response, err := c.httpClient.Do(
		request,
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {

		return nil, fmt.Errorf(
			"musicbrainz search failed: %s",
			response.Status,
		)

	}

	var result searchResponse

	if err := json.NewDecoder(
		response.Body,
	).Decode(&result); err != nil {

		return nil, err
	}

	recordings := make(
		[]Recording,
		0,
		len(result.Recordings),
	)

	for _, response := range result.Recordings {

		recordings = append(
			recordings,
			recordingFromResponse(
				response,
			),
		)

	}

	return recordings, nil

}

func recordingQuery(
	artist string,
	title string,
) string {

	return fmt.Sprintf(
		`artist:"%s" AND recording:"%s"`,
		escapeQueryValue(artist),
		escapeQueryValue(title),
	)

}

func escapeQueryValue(
	value string,
) string {

	replacer := strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
	)

	return replacer.Replace(
		strings.TrimSpace(value),
	)

}

func artistCreditName(
	credits []artistCreditResponse,
) string {

	var builder strings.Builder

	for _, credit := range credits {

		name := credit.Name

		if name == "" {
			name = credit.Artist.Name
		}

		builder.WriteString(name)
		builder.WriteString(
			credit.JoinPhrase,
		)

	}

	return strings.TrimSpace(
		builder.String(),
	)

}

func recordingFromResponse(
	response recordingResponse,
) Recording {

	releases := make(
		[]Release,
		0,
		len(response.Releases),
	)

	for _, release := range response.Releases {

		releases = append(
			releases,
			Release{
				ID:      release.ID,
				Title:   release.Title,
				Date:    release.Date,
				Country: release.Country,
			},
		)

	}

	return Recording{
		ID: response.ID,

		Artist: artistCreditName(
			response.ArtistCredit,
		),
		Title: response.Title,

		Score: response.Score,
		Duration: time.Duration(
			response.Length,
		) * time.Millisecond,

		ISRCs: response.ISRCs,

		Disambiguation: response.Disambiguation,
		Releases:       releases,
	}

}

func (c *Client) LookupRecording(
	ctx context.Context,
	id string,
) (Recording, error) {

	endpoint, err := url.Parse(
		c.baseURL + "/recording/" + url.PathEscape(id),
	)
	if err != nil {
		return Recording{}, err
	}

	values := endpoint.Query()
	values.Set("fmt", "json")
	values.Set(
		"inc",
		"artist-credits+isrcs+releases",
	)

	endpoint.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint.String(),
		nil,
	)
	if err != nil {
		return Recording{}, err
	}

	request.Header.Set(
		"User-Agent",
		c.userAgent,
	)

	if err := c.waitForRequest(
		ctx,
	); err != nil {

		return Recording{}, err
	}

	response, err := c.httpClient.Do(
		request,
	)
	if err != nil {
		return Recording{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {

		return Recording{}, fmt.Errorf(
			"musicbrainz recording lookup failed: %s",
			response.Status,
		)

	}

	var result recordingResponse

	if err := json.NewDecoder(
		response.Body,
	).Decode(&result); err != nil {

		return Recording{}, err
	}

	return recordingFromResponse(
		result,
	), nil

}
