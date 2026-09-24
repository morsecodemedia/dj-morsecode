package musicbrainz

import (
	"context"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

type Enricher struct {
	client *Client
}

func NewEnricher(
	client *Client,
) *Enricher {

	return &Enricher{
		client: client,
	}

}

func (e *Enricher) Enrich(
	ctx context.Context,
	observed metadata.PlaybackItem,
) (
	metadata.EnrichmentMatch,
	metadata.MatchStatus,
	error,
) {

	if !observed.IsTrack() {

		return metadata.EnrichmentMatch{},
			metadata.MatchNone,
			nil

	}

	recordings, err := e.client.SearchRecordings(
		ctx,
		observed.Artist,
		observed.Title,
	)
	if err != nil {

		return metadata.EnrichmentMatch{},
			metadata.MatchNone,
			err

	}

	candidates := make(
		[]metadata.EnrichmentCandidate,
		0,
		len(recordings),
	)

	for _, recording := range recordings {

		candidates = append(
			candidates,
			Candidate(
				recording,
			),
		)

	}

	match, status := metadata.MatchEnrichment(
		observed,
		candidates,
	)

	return match, status, nil

}
