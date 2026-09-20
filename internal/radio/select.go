package radio

type SelectionOptions struct {
	RecentLimit int
}

func Select(
	candidates []Station,
	history History,
	options SelectionOptions,
) (Station, bool) {

	if len(candidates) == 0 {
		return Station{}, false
	}

	for _, station := range candidates {

		if history.ContainsRecent(
			station.ID,
			options.RecentLimit,
		) {
			continue
		}

		return station, true

	}

	return Station{}, false

}
