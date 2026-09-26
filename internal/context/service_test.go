package trackcontext

import (
	"context"
	"errors"
	"testing"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

type fakeProvider struct {
	context metadata.TrackContext
	ok      bool
	err     error
}

func (p fakeProvider) Contextualize(
	ctx context.Context,
	track metadata.CanonicalTrack,
) (
	metadata.TrackContext,
	bool,
	error,
) {

	return p.context,
		p.ok,
		p.err

}

func TestServiceMergesProviderContext(
	t *testing.T,
) {

	releaseProvider := fakeProvider{
		ok: true,
		context: metadata.TrackContext{
			Release: metadata.ReleaseContext{
				Title: "Unheard",
				Date:  "2024-03-22",
			},
		},
	}

	tagProvider := fakeProvider{
		ok: true,
		context: metadata.TrackContext{
			Tags: []metadata.ContextTag{
				{
					Name:        "soul",
					ArtistCount: 56,
					Providers: []string{
						"lastfm",
					},
				},
			},
		},
	}

	service := NewService(
		releaseProvider,
		tagProvider,
	)

	result := service.Contextualize(
		context.Background(),
		metadata.CanonicalTrack{
			Artist: "Hozier",
			Title:  "Too Sweet",
		},
	)

	if !result.HasContext() {
		t.Fatal("expected merged context")
	}

	if result.Context.Release.Title !=
		"Unheard" {

		t.Errorf(
			"expected release %q, got %q",
			"Unheard",
			result.Context.Release.Title,
		)

	}

	if len(result.Context.Tags) != 1 {
		t.Fatalf(
			"expected one tag, got %d",
			len(result.Context.Tags),
		)
	}

}

func TestServicePreservesPartialContext(
	t *testing.T,
) {

	service := NewService(
		fakeProvider{
			ok: true,
			context: metadata.TrackContext{
				Release: metadata.ReleaseContext{
					Title: "Touchdown",
					Date:  "1978",
				},
			},
		},
		fakeProvider{
			err: errors.New(
				"last.fm unavailable",
			),
		},
	)

	result := service.Contextualize(
		context.Background(),
		metadata.CanonicalTrack{},
	)

	if !result.HasContext() {
		t.Fatal(
			"expected partial context to survive",
		)
	}

	if result.Failed() {
		t.Fatal(
			"expected partial result not to be total failure",
		)
	}

	if len(result.Errors) != 1 {
		t.Fatalf(
			"expected one provider error, got %d",
			len(result.Errors),
		)
	}

}

func TestServiceReportsTotalFailure(
	t *testing.T,
) {

	service := NewService(
		fakeProvider{
			err: errors.New("provider one"),
		},
		fakeProvider{
			err: errors.New("provider two"),
		},
	)

	result := service.Contextualize(
		context.Background(),
		metadata.CanonicalTrack{},
	)

	if !result.Failed() {
		t.Fatal(
			"expected total failure",
		)
	}

	if result.HasContext() {
		t.Fatal(
			"expected no context",
		)
	}

}

func TestServiceAllowsEmptyContext(
	t *testing.T,
) {

	service := NewService(
		fakeProvider{},
		fakeProvider{},
	)

	result := service.Contextualize(
		context.Background(),
		metadata.CanonicalTrack{},
	)

	if result.HasContext() {
		t.Fatal(
			"expected no context",
		)
	}

	if result.Failed() {
		t.Fatal(
			"expected empty context not to be failure",
		)
	}

}
