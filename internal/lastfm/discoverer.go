package lastfm

import (
	"context"
	"fmt"

	"github.com/morsecodemedia/dj-morsecode/internal/source"
)

type Discoverer struct {
	client *Client
}

func NewDiscoverer(
	client *Client,
) *Discoverer {

	return &Discoverer{
		client: client,
	}

}

var _ source.Discoverer = (*Discoverer)(nil)

func (d *Discoverer) Discover(
	ctx context.Context,
	request source.DiscoveryRequest,
) ([]source.DiscoveryResult, error) {

	if request.Kind !=
		source.DiscoverySimilarTrack {

		return nil, fmt.Errorf(
			"last.fm discovery does not support %q",
			request.Kind,
		)

	}

	if request.Artist == "" ||
		request.Title == "" {

		return nil, fmt.Errorf(
			"last.fm similar-track discovery requires artist and title",
		)

	}

	tracks, err := d.client.TrackSimilar(
		ctx,
		request.Artist,
		request.Title,
		request.Limit,
	)
	if err != nil {
		return nil, err
	}

	results := make(
		[]source.DiscoveryResult,
		0,
		len(tracks),
	)

	for _, candidate := range tracks {

		item, ok := MediaItem(
			candidate.Track,
		)

		if !ok {
			continue
		}

		results = append(
			results,
			source.DiscoveryResult{
				Item:  item,
				Score: candidate.Match,
			},
		)

	}

	return results, nil

}
