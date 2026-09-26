package lastfm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const defaultBaseURL = "https://ws.audioscrobbler.com/2.0/"

type Client struct {
	apiKey string

	baseURL    string
	httpClient *http.Client
}

type Tag struct {
	Name  string
	Count int
}

type topTagsResponse struct {
	TopTags struct {
		Tags []struct {
			Name  string `json:"name"`
			Count int    `json:"count"`
		} `json:"tag"`
	} `json:"toptags"`
}

func NewClient(
	apiKey string,
) *Client {

	return &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: http.DefaultClient,
	}

}

func (c *Client) TrackTopTags(
	ctx context.Context,
	artist string,
	track string,
) ([]Tag, error) {

	endpoint, err := url.Parse(
		c.baseURL,
	)
	if err != nil {
		return nil, err
	}

	values := endpoint.Query()

	values.Set(
		"method",
		"track.getTopTags",
	)

	values.Set(
		"api_key",
		c.apiKey,
	)

	values.Set(
		"artist",
		artist,
	)

	values.Set(
		"track",
		track,
	)

	values.Set(
		"autocorrect",
		"1",
	)

	values.Set(
		"format",
		"json",
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

	response, err := c.httpClient.Do(
		request,
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {

		return nil, fmt.Errorf(
			"last.fm track tags failed: %s",
			response.Status,
		)

	}

	var result topTagsResponse

	if err := json.NewDecoder(
		response.Body,
	).Decode(&result); err != nil {

		return nil, err
	}

	tags := make(
		[]Tag,
		0,
		len(result.TopTags.Tags),
	)

	for _, tag := range result.TopTags.Tags {

		if tag.Name == "" {
			continue
		}

		tags = append(
			tags,
			Tag{
				Name:  tag.Name,
				Count: tag.Count,
			},
		)

	}

	return tags, nil

}
