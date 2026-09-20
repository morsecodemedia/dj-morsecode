package radio

import (
	"sort"
	"strings"
)

func Genres() []string {

	values := make(map[string]string)

	for _, station := range Stations {

		value := strings.TrimSpace(
			station.Genre,
		)

		if value == "" {
			continue
		}

		key := strings.ToLower(value)

		if _, exists := values[key]; !exists {
			values[key] = value
		}

	}

	return sortedFacetValues(values)

}

func Moods() []string {

	values := make(map[string]string)

	for _, station := range Stations {

		for _, mood := range station.Moods {

			value := strings.TrimSpace(
				mood,
			)

			if value == "" {
				continue
			}

			key := strings.ToLower(value)

			if _, exists := values[key]; !exists {
				values[key] = value
			}

		}

	}

	return sortedFacetValues(values)

}

func sortedFacetValues(
	values map[string]string,
) []string {

	result := make(
		[]string,
		0,
		len(values),
	)

	for _, value := range values {
		result = append(
			result,
			value,
		)
	}

	sort.Slice(
		result,
		func(i, j int) bool {

			return strings.ToLower(result[i]) <
				strings.ToLower(result[j])

		},
	)

	return result

}
