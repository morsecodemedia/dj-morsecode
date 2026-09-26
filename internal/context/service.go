package trackcontext

import (
	"context"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

type Provider interface {
	Contextualize(
		ctx context.Context,
		track metadata.CanonicalTrack,
	) (
		metadata.TrackContext,
		bool,
		error,
	)
}

type Service struct {
	providers []Provider
}

func NewService(
	providers ...Provider,
) *Service {

	return &Service{
		providers: providers,
	}

}

type Result struct {
	Context metadata.TrackContext
	Errors  []error
}

func (r Result) HasContext() bool {

	return r.Context.Release.Title != "" ||
		r.Context.Release.Date != "" ||
		len(r.Context.Tags) > 0 ||
		len(r.Context.Genres) > 0

}

func (r Result) Failed() bool {

	return !r.HasContext() &&
		len(r.Errors) > 0

}

func (s *Service) Contextualize(
	ctx context.Context,
	track metadata.CanonicalTrack,
) Result {

	var result Result

	for _, provider := range s.providers {

		if provider == nil {
			continue
		}

		partial, ok, err := provider.Contextualize(
			ctx,
			track,
		)

		if err != nil {

			result.Errors = append(
				result.Errors,
				err,
			)

			continue
		}

		if !ok {
			continue
		}

		result.Context = mergeContext(
			result.Context,
			partial,
		)

	}

	return result

}

func mergeContext(
	current metadata.TrackContext,
	incoming metadata.TrackContext,
) metadata.TrackContext {

	if current.Release.Title == "" &&
		current.Release.Date == "" {

		current.Release =
			incoming.Release
	}

	current.Tags = mergeTags(
		current.Tags,
		incoming.Tags,
	)

	if len(current.Genres) == 0 &&
		len(incoming.Genres) > 0 {

		current.Genres = append(
			[]string(nil),
			incoming.Genres...,
		)

	}

	return current

}

func mergeTags(
	current []metadata.ContextTag,
	incoming []metadata.ContextTag,
) []metadata.ContextTag {

	evidence := make(
		[]metadata.TagEvidence,
		0,
		len(current)+len(incoming),
	)

	for _, tag := range current {

		for _, provider := range tag.Providers {

			if tag.TrackCount > 0 {

				evidence = append(
					evidence,
					metadata.TagEvidence{
						Name:     tag.Name,
						Count:    tag.TrackCount,
						Scope:    metadata.TagScopeTrack,
						Provider: provider,
					},
				)

			}

			if tag.ArtistCount > 0 {

				evidence = append(
					evidence,
					metadata.TagEvidence{
						Name:     tag.Name,
						Count:    tag.ArtistCount,
						Scope:    metadata.TagScopeArtist,
						Provider: provider,
					},
				)

			}

		}

	}

	for _, tag := range incoming {

		for _, provider := range tag.Providers {

			if tag.TrackCount > 0 {

				evidence = append(
					evidence,
					metadata.TagEvidence{
						Name:     tag.Name,
						Count:    tag.TrackCount,
						Scope:    metadata.TagScopeTrack,
						Provider: provider,
					},
				)

			}

			if tag.ArtistCount > 0 {

				evidence = append(
					evidence,
					metadata.TagEvidence{
						Name:     tag.Name,
						Count:    tag.ArtistCount,
						Scope:    metadata.TagScopeArtist,
						Provider: provider,
					},
				)

			}

		}

	}

	return metadata.ResolveContextTags(
		evidence,
	)

}
