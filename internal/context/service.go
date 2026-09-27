package trackcontext

import (
	"context"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

type Task func(
	context.Context,
	metadata.CanonicalTrack,
) (
	metadata.TrackContext,
	bool,
	error,
)

func (s *Service) Tasks() []Task {

	tasks := make(
		[]Task,
		0,
		len(s.providers),
	)

	for _, provider := range s.providers {

		if provider == nil {
			continue
		}

		tasks = append(
			tasks,
			provider.Contextualize,
		)

	}

	return tasks

}

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

		result.Context = metadata.MergeTrackContext(
			result.Context,
			partial,
		)

	}

	return result

}
