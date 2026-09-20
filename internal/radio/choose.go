package radio

type ChooseOptions struct {
	RecentLimit int
}

func Choose(
	criteria Criteria,
	history History,
	options ChooseOptions,
) (Station, bool) {

	candidates := Match(
		criteria,
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

		station, ok := Select(
			filtered,
			history,
			SelectionOptions{
				RecentLimit: recentLimit,
			},
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
