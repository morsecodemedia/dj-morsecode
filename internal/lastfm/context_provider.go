package lastfm

import (
	"context"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

type ContextProvider struct {
	client *Client
}

func NewContextProvider(
	client *Client,
) *ContextProvider {

	return &ContextProvider{
		client: client,
	}

}

func (p *ContextProvider) Contextualize(
	ctx context.Context,
	track metadata.CanonicalTrack,
) (
	metadata.TrackContext,
	bool,
	error,
) {

	if track.Artist == "" ||
		track.Title == "" {

		return metadata.TrackContext{},
			false,
			nil

	}

	result := p.client.TrackTagEvidence(
		ctx,
		track.Artist,
		track.Title,
	)

	if result.Failed() {

		return metadata.TrackContext{},
			false,
			result.Error()

	}

	tags := metadata.ResolveContextTags(
		result.Evidence,
	)

	if len(tags) == 0 {

		return metadata.TrackContext{},
			false,
			nil

	}

	return metadata.TrackContext{
		Tags: tags,
	}, true, nil

}
