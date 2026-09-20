package radio

import "strings"

type Criteria struct {
	Genres   []string
	Tags     []string
	Moods    []string
	Contexts []string

	MinEnergy int
	MaxEnergy int
}

func Match(criteria Criteria) []Station {

	var matches []Station

	for _, station := range Stations {

		if !matchesAny(
			[]string{station.Genre},
			criteria.Genres,
		) {
			continue
		}

		if !matchesAny(
			station.Tags,
			criteria.Tags,
		) {
			continue
		}

		if !matchesAny(
			station.Moods,
			criteria.Moods,
		) {
			continue
		}

		if !matchesAny(
			station.Contexts,
			criteria.Contexts,
		) {
			continue
		}

		if criteria.MinEnergy > 0 &&
			station.Energy < criteria.MinEnergy {

			continue
		}

		if criteria.MaxEnergy > 0 &&
			station.Energy > criteria.MaxEnergy {

			continue
		}

		matches = append(
			matches,
			station,
		)

	}

	return matches

}

func matchesAny(
	values []string,
	wanted []string,
) bool {

	if len(wanted) == 0 {
		return true
	}

	for _, candidate := range values {

		for _, target := range wanted {

			if strings.EqualFold(
				candidate,
				target,
			) {
				return true
			}

		}

	}

	return false

}
