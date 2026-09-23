package radio

import "math/rand/v2"

type SelectionOptions struct {
	RecentLimit int
}

type CandidateChooser func(
	candidates []Station,
) (Station, bool)

func RandomCandidate(
	candidates []Station,
) (Station, bool) {

	if len(candidates) == 0 {
		return Station{}, false
	}

	index := rand.IntN(
		len(candidates),
	)

	return candidates[index], true

}

func Select(
	candidates []Station,
	history History,
	options SelectionOptions,
) (Station, bool) {

	return SelectWithChooser(
		candidates,
		history,
		options,
		nil,
	)

}

func SelectWithChooser(
	candidates []Station,
	history History,
	options SelectionOptions,
	chooser CandidateChooser,
) (Station, bool) {

	if len(candidates) == 0 {
		return Station{}, false
	}

	var eligible []Station

	for _, station := range candidates {

		if history.ContainsRecent(
			station.ID,
			options.RecentLimit,
		) {
			continue
		}

		eligible = append(
			eligible,
			station,
		)

	}

	if len(eligible) == 0 {
		return Station{}, false
	}

	if chooser == nil {
		return eligible[0], true
	}

	return chooser(
		eligible,
	)

}
