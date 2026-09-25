package metadata

import "strings"

func SelectReleaseContext(
	evidence []ReleaseEvidence,
) (ReleaseContext, bool) {

	groups := uniqueReleaseGroups(
		evidence,
	)

	var candidates []ReleaseEvidence

	for _, candidate := range groups {

		if candidate.GroupID == "" {
			continue
		}

		if candidate.GroupTitle == "" {
			continue
		}

		if isCompilation(
			candidate.SecondaryTypes,
		) {
			continue
		}

		if !sameReleaseDate(
			candidate.ReleaseDate,
			candidate.GroupFirstReleaseDate,
		) {
			continue
		}

		candidates = append(
			candidates,
			candidate,
		)

	}

	if len(candidates) != 1 {
		return ReleaseContext{}, false
	}

	selected := candidates[0]

	return ReleaseContext{
		Title: selected.GroupTitle,
		Date:  selected.GroupFirstReleaseDate,
	}, true

}

func uniqueReleaseGroups(
	evidence []ReleaseEvidence,
) []ReleaseEvidence {

	groups := make(
		map[string]ReleaseEvidence,
	)

	order := make(
		[]string,
		0,
		len(evidence),
	)

	for _, candidate := range evidence {

		id := strings.TrimSpace(
			candidate.GroupID,
		)

		if id == "" {
			continue
		}

		existing, ok := groups[id]
		if !ok {

			groups[id] = candidate

			order = append(
				order,
				id,
			)

			continue
		}

		if !sameReleaseDate(
			existing.ReleaseDate,
			existing.GroupFirstReleaseDate,
		) &&
			sameReleaseDate(
				candidate.ReleaseDate,
				candidate.GroupFirstReleaseDate,
			) {

			groups[id] = candidate

		}

	}

	result := make(
		[]ReleaseEvidence,
		0,
		len(order),
	)

	for _, id := range order {

		result = append(
			result,
			groups[id],
		)

	}

	return result

}

func isCompilation(
	secondaryTypes []string,
) bool {

	for _, secondaryType := range secondaryTypes {

		if strings.EqualFold(
			strings.TrimSpace(secondaryType),
			"Compilation",
		) {

			return true
		}

	}

	return false

}

func sameReleaseDate(
	releaseDate string,
	groupFirstReleaseDate string,
) bool {

	releaseDate = strings.TrimSpace(
		releaseDate,
	)

	groupFirstReleaseDate = strings.TrimSpace(
		groupFirstReleaseDate,
	)

	if releaseDate == "" ||
		groupFirstReleaseDate == "" {

		return false
	}

	return releaseDate ==
		groupFirstReleaseDate

}
