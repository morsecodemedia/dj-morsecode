package lastfm

import (
	"context"

	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

const providerName = "lastfm"

func TagEvidence(
	tags []Tag,
	scope metadata.TagScope,
) []metadata.TagEvidence {

	evidence := make(
		[]metadata.TagEvidence,
		0,
		len(tags),
	)

	for _, tag := range tags {

		if tag.Name == "" {
			continue
		}

		evidence = append(
			evidence,
			metadata.TagEvidence{
				Name:     tag.Name,
				Count:    tag.Count,
				Scope:    scope,
				Provider: providerName,
			},
		)

	}

	return evidence

}

func (c *Client) TrackTagEvidence(
	ctx context.Context,
	artist string,
	title string,
) TagEvidenceResult {

	var result TagEvidenceResult

	trackTags, err := c.TrackTopTags(
		ctx,
		artist,
		title,
	)

	if err != nil {

		result.TrackError = err

	} else {

		result.Evidence = append(
			result.Evidence,
			TagEvidence(
				trackTags,
				metadata.TagScopeTrack,
			)...,
		)

	}

	artistTags, err := c.ArtistTopTags(
		ctx,
		artist,
	)

	if err != nil {

		result.ArtistError = err

	} else {

		result.Evidence = append(
			result.Evidence,
			TagEvidence(
				artistTags,
				metadata.TagScopeArtist,
			)...,
		)

	}

	return result

}

type TagEvidenceResult struct {
	Evidence []metadata.TagEvidence

	TrackError  error
	ArtistError error
}

func (r TagEvidenceResult) HasEvidence() bool {
	return len(r.Evidence) > 0
}

func (r TagEvidenceResult) Failed() bool {
	return r.TrackError != nil &&
		r.ArtistError != nil
}
