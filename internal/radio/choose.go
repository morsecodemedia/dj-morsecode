package radio

type ChooseOptions struct {
	RecentLimit int
	Chooser     CandidateChooser
	ExcludeIDs  []string
}

func Choose(
	criteria Criteria,
	history History,
	options ChooseOptions,
) (Station, bool) {

	candidates := Match(
		criteria,
	)

	candidates = excludeStations(
		candidates,
		options.ExcludeIDs,
	)

	if len(candidates) == 0 {
		return Station{}, false
	}

	current, hasCurrent := history.Current()

	for recentLimit := options.RecentLimit; recentLimit >= 0; recentLimit-- {

		filtered := excludeCurrent(
			candidates,
			current,
			hasCurrent,
		)

		station, ok := SelectWithChooser(
			filtered,
			history,
			SelectionOptions{
				RecentLimit: recentLimit,
			},
			options.Chooser,
		)
		if ok {
			return station, true
		}

	}

	return Station{}, false

}

func excludeCurrent(
	stations []Station,
	current Tune,
	hasCurrent bool,
) []Station {

	if !hasCurrent {
		return stations
	}

	var filtered []Station

	for _, station := range stations {

		if station.ID == current.StationID {
			continue
		}

		filtered = append(
			filtered,
			station,
		)

	}

	return filtered

}

func excludeStations(
	stations []Station,
	excludeIDs []string,
) []Station {

	if len(excludeIDs) == 0 {
		return stations
	}

	excluded := make(map[string]bool)

	for _, id := range excludeIDs {
		excluded[id] = true
	}

	var filtered []Station

	for _, station := range stations {

		if excluded[station.ID] {
			continue
		}

		filtered = append(
			filtered,
			station,
		)

	}

	return filtered

}
