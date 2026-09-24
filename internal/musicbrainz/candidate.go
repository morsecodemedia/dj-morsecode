package musicbrainz

import (
	"github.com/morsecodemedia/dj-morsecode/internal/metadata"
)

const providerName = "musicbrainz"

func Candidate(
	recording Recording,
) metadata.EnrichmentCandidate {

	identifiers := []metadata.Identifier{
		{
			Scheme: metadata.IdentifierMusicBrainz,
			Value:  recording.ID,
		},
	}

	for _, isrc := range recording.ISRCs {

		if isrc == "" {
			continue
		}

		identifiers = append(
			identifiers,
			metadata.Identifier{
				Scheme: metadata.IdentifierISRC,
				Value:  isrc,
			},
		)

	}

	return metadata.EnrichmentCandidate{
		Artist: recording.Artist,
		Title:  recording.Title,

		Duration: recording.Duration,
		Variant:  recording.Disambiguation,

		Identifiers: identifiers,

		Provider: providerName,
		ProviderScore: float64(
			recording.Score,
		) / 100,
	}

}
