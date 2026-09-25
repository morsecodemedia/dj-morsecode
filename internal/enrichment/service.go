package enrichment

import (
	"context"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

type Provider interface {
	Enrich(
		context.Context,
		metadata.PlaybackItem,
	) (
		metadata.EnrichmentMatch,
		metadata.MatchStatus,
		error,
	)
}

type Service struct {
	store    *metadata.EnrichmentStore
	provider Provider
}

func NewService(
	store *metadata.EnrichmentStore,
	provider Provider,
) *Service {

	return &Service{
		store:    store,
		provider: provider,
	}

}

func (s *Service) Enrich(
	ctx context.Context,
	item metadata.PlaybackItem,
) (
	metadata.EnrichmentMatch,
	metadata.MatchStatus,
	error,
) {

	if !item.IsTrack() {

		return metadata.EnrichmentMatch{},
			metadata.MatchNone,
			nil

	}

	if match, ok := s.store.Get(
		item.Artist,
		item.Title,
	); ok {

		return match,
			metadata.MatchAccepted,
			nil

	}

	match, status, err := s.provider.Enrich(
		ctx,
		item,
	)
	if err != nil {

		return metadata.EnrichmentMatch{},
			metadata.MatchNone,
			err

	}

	if status == metadata.MatchAccepted {

		s.store.Put(
			item.Artist,
			item.Title,
			match,
		)

	}

	return match, status, nil

}
