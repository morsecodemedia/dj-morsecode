package metadata

import (
	"strings"
	"time"
)

const minimumProviderScore = 0.90

const durationTolerance = 5 * time.Second

type MatchStatus int

const (
	MatchNone MatchStatus = iota
	MatchAccepted
	MatchAmbiguous
)

func MatchEnrichment(
	observed PlaybackItem,
	candidates []EnrichmentCandidate,
) (EnrichmentMatch, MatchStatus) {

	if !observed.IsTrack() {
		return EnrichmentMatch{}, MatchNone
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

	if len(accepted) == 0 {
		return EnrichmentMatch{}, MatchNone
	}

	candidate, ok := resolveCandidate(
		accepted,
	)
	if !ok {
		return EnrichmentMatch{}, MatchAmbiguous
	}

	return EnrichmentMatch{
		Track: CanonicalTrack{
			Artist: candidate.Artist,
			Title:  candidate.Title,

			Identifiers: candidate.Identifiers,
		},
		Provider:   candidate.Provider,
		Confidence: candidate.ProviderScore,
	}, MatchAccepted

}

func resolveCandidate(
	candidates []EnrichmentCandidate,
) (EnrichmentCandidate, bool) {

	if len(candidates) == 1 {
		return candidates[0], true
	}

	var unqualified []EnrichmentCandidate

	for _, candidate := range candidates {

		if strings.TrimSpace(
			candidate.Variant,
		) != "" {

			continue
		}

		unqualified = append(
			unqualified,
			candidate,
		)

	}

	if len(unqualified) != 1 {
		return EnrichmentCandidate{}, false
	}

	return unqualified[0], true

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
