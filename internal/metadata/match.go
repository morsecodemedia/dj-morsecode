package metadata

import (
	"strings"
	"time"
)

const minimumProviderScore = 0.90

const durationTolerance = 5 * time.Second

func MatchEnrichment(
	observed PlaybackItem,
	candidates []EnrichmentCandidate,
) (EnrichmentMatch, bool) {

	if !observed.IsTrack() {
		return EnrichmentMatch{}, false
	}

	var accepted []EnrichmentCandidate

	for _, candidate := range candidates {

		if candidate.ProviderScore < minimumProviderScore {
			continue
		}

		if !sameIdentity(
			observed.Artist,
			candidate.Artist,
		) {

			continue
		}

		if !sameIdentity(
			observed.Title,
			candidate.Title,
		) {

			continue
		}

		if !durationMatches(
			observed.Duration,
			candidate.Duration,
		) {

			continue
		}

		accepted = append(
			accepted,
			candidate,
		)

	}

	if len(accepted) != 1 {
		return EnrichmentMatch{}, false
	}

	candidate := accepted[0]

	return EnrichmentMatch{
		Track: CanonicalTrack{
			Artist: candidate.Artist,
			Title:  candidate.Title,

			Identifiers: candidate.Identifiers,
		},
		Provider:   candidate.Provider,
		Confidence: candidate.ProviderScore,
	}, true

}

func sameIdentity(
	left string,
	right string,
) bool {

	return strings.EqualFold(
		strings.TrimSpace(left),
		strings.TrimSpace(right),
	)

}

func durationMatches(
	observed time.Duration,
	candidate time.Duration,
) bool {

	if observed <= 0 ||
		candidate <= 0 {

		return true
	}

	delta := observed - candidate

	if delta < 0 {
		delta = -delta
	}

	return delta <= durationTolerance

}
