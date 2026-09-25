package enrichment

import (
	"context"
	"errors"
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

type fakeProvider struct {
	match  metadata.EnrichmentMatch
	status metadata.MatchStatus
	err    error

	calls int
}

func (p *fakeProvider) Enrich(
	ctx context.Context,
	item metadata.PlaybackItem,
) (
	metadata.EnrichmentMatch,
	metadata.MatchStatus,
	error,
) {

	p.calls++

	return p.match,
		p.status,
		p.err

}

func TestServiceCachesAcceptedMatch(
	t *testing.T,
) {

	store := metadata.NewEnrichmentStore()

	provider := &fakeProvider{
		match: metadata.EnrichmentMatch{
			Track: metadata.CanonicalTrack{
				Artist: "Hozier",
				Title:  "Too Sweet",
				Identifiers: []metadata.Identifier{
					{
						Scheme: metadata.IdentifierMusicBrainz,
						Value:  "recording-id",
					},
				},
			},
			Provider:   "musicbrainz",
			Confidence: 1,
		},
		status: metadata.MatchAccepted,
	}

	service := NewService(
		store,
		provider,
	)

	item := metadata.PlaybackItem{
		Type:   metadata.PlaybackTrack,
		Artist: "HOZIER",
		Title:  "Too Sweet",
	}

	first, status, err := service.Enrich(
		context.Background(),
		item,
	)
	if err != nil {

		t.Fatalf(
			"first enrichment returned error: %v",
			err,
		)

	}

	if status != metadata.MatchAccepted {

		t.Fatalf(
			"expected accepted match, got %v",
			status,
		)

	}

	second, status, err := service.Enrich(
		context.Background(),
		item,
	)
	if err != nil {

		t.Fatalf(
			"second enrichment returned error: %v",
			err,
		)

	}

	if status != metadata.MatchAccepted {

		t.Fatalf(
			"expected cached accepted match, got %v",
			status,
		)

	}

	if provider.calls != 1 {

		t.Errorf(
			"expected provider called once, got %d",
			provider.calls,
		)

	}

	if first.Track.Title != second.Track.Title {

		t.Errorf(
			"expected cached match title %q, got %q",
			first.Track.Title,
			second.Track.Title,
		)

	}

}
func TestServiceSkipsNonTrack(
	t *testing.T,
) {

	store := metadata.NewEnrichmentStore()
	provider := &fakeProvider{}

	service := NewService(
		store,
		provider,
	)

	item := metadata.PlaybackItem{
		Type:  metadata.PlaybackStationID,
		Title: "Z100",
	}

	_, status, err := service.Enrich(
		context.Background(),
		item,
	)
	if err != nil {

		t.Fatalf(
			"enrichment returned error: %v",
			err,
		)

	}

	if status != metadata.MatchNone {

		t.Fatalf(
			"expected no match, got %v",
			status,
		)

	}

	if provider.calls != 0 {

		t.Errorf(
			"expected provider not called, got %d calls",
			provider.calls,
		)

	}

}
func TestServiceDoesNotCacheAmbiguousMatch(
	t *testing.T,
) {

	store := metadata.NewEnrichmentStore()

	provider := &fakeProvider{
		status: metadata.MatchAmbiguous,
	}

	service := NewService(
		store,
		provider,
	)

	item := metadata.PlaybackItem{
		Type:   metadata.PlaybackTrack,
		Artist: "Hozier",
		Title:  "Too Sweet",
	}

	for i := 0; i < 2; i++ {

		_, status, err := service.Enrich(
			context.Background(),
			item,
		)
		if err != nil {

			t.Fatalf(
				"enrichment returned error: %v",
				err,
			)

		}

		if status != metadata.MatchAmbiguous {

			t.Fatalf(
				"expected ambiguous match, got %v",
				status,
			)

		}

	}

	if provider.calls != 2 {

		t.Errorf(
			"expected ambiguous result not cached; got %d provider calls",
			provider.calls,
		)

	}

}

func TestServiceDoesNotCacheProviderError(
	t *testing.T,
) {

	store := metadata.NewEnrichmentStore()

	provider := &fakeProvider{
		err: errors.New(
			"provider exploded",
		),
	}

	service := NewService(
		store,
		provider,
	)

	item := metadata.PlaybackItem{
		Type:   metadata.PlaybackTrack,
		Artist: "Hozier",
		Title:  "Too Sweet",
	}

	for i := 0; i < 2; i++ {

		_, _, err := service.Enrich(
			context.Background(),
			item,
		)

		if err == nil {

			t.Fatal(
				"expected provider error",
			)

		}

	}

	if provider.calls != 2 {

		t.Errorf(
			"expected errors not cached; got %d provider calls",
			provider.calls,
		)

	}

}
