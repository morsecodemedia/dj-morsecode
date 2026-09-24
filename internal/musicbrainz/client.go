package musicbrainz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const defaultBaseURL = "https://musicbrainz.org/ws/2"

const defaultLimit = 10

type Client struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string
}

func NewClient(
	userAgent string,
) *Client {

	return &Client{
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		userAgent: userAgent,
	}

}

type searchResponse struct {
	Recordings []recordingResponse `json:"recordings"`
}

type recordingResponse struct {
	ID     string `json:"id"`
	Score  int    `json:"score"`
	Title  string `json:"title"`
	Length int64  `json:"length"`

	ArtistCredit []artistCreditResponse `json:"artist-credit"`

	ISRCs []string `json:"isrcs"`
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

	for _, candidate := range result.Recordings {

		recordings = append(
			recordings,
			Recording{
				ID: candidate.ID,

				Artist: artistCreditName(
					candidate.ArtistCredit,
				),
				Title: candidate.Title,

				Score: candidate.Score,
				Duration: time.Duration(
					candidate.Length,
				) * time.Millisecond,

				ISRCs: candidate.ISRCs,
			},
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
