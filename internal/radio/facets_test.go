package radio

import (
	"slices"
	"testing"
)

func TestGenres(t *testing.T) {

	genres := Genres()

	if len(genres) == 0 {
		t.Fatal("expected genres")
	}

	assertFacetContains(
		t,
		genres,
		"Ambient",
	)

	assertFacetContains(
		t,
		genres,
		"Eclectic",
	)

	assertFacetContains(
		t,
		genres,
		"Trance",
	)

	assertUniqueFacets(
		t,
		genres,
	)

	assertSortedFacets(
		t,
		genres,
	)

}

func TestMoods(t *testing.T) {

	moods := Moods()

	if len(moods) == 0 {
		t.Fatal("expected moods")
	}

	assertFacetContains(
		t,
		moods,
		"calm",
	)

	assertFacetContains(
		t,
		moods,
		"energetic",
	)

	assertFacetContains(
		t,
		moods,
		"nostalgic",
	)

	assertUniqueFacets(
		t,
		moods,
	)

	assertSortedFacets(
		t,
		moods,
	)

}

func assertFacetContains(
	t *testing.T,
	values []string,
	wanted string,
) {

	t.Helper()

	if !slices.Contains(
		values,
		wanted,
	) {

		t.Errorf(
			"expected facet %q",
			wanted,
		)

	}

}

func assertUniqueFacets(
	t *testing.T,
	values []string,
) {

	t.Helper()

	seen := make(
		map[string]bool,
	)

	for _, value := range values {

		if seen[value] {

			t.Errorf(
				"duplicate facet %q",
				value,
			)

		}

		seen[value] = true

	}

}

func assertSortedFacets(
	t *testing.T,
	values []string,
) {

	t.Helper()

	if !slices.IsSorted(values) {

		t.Errorf(
			"expected sorted facets, got %v",
			values,
		)

	}

}
