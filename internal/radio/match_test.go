package radio

import "testing"

func TestMatchByContextAndEnergy(t *testing.T) {

	results := Match(
		Criteria{
			Contexts: []string{
				"coding",
			},
			MaxEnergy: 2,
		},
	)

	assertContainsStation(
		t,
		results,
		"lofi247",
	)

	assertContainsStation(
		t,
		results,
		"somafm-groovesalad",
	)

	assertNotContainsStation(
		t,
		results,
		"hardrockradiofm",
	)

}

func TestMatchByMoodAndMinimumEnergy(t *testing.T) {

	results := Match(
		Criteria{
			Moods: []string{
				"energetic",
			},
			MinEnergy: 4,
		},
	)

	assertContainsStation(
		t,
		results,
		"z100",
	)

	assertContainsStation(
		t,
		results,
		"hardrockradiofm",
	)

	assertNotContainsStation(
		t,
		results,
		"somafm-dronezone",
	)

}

func TestMatchUsesOrWithinDimension(t *testing.T) {

	results := Match(
		Criteria{
			Contexts: []string{
				"gaming",
				"coding",
			},
		},
	)

	assertContainsStation(
		t,
		results,
		"slayradio",
	)

	assertContainsStation(
		t,
		results,
		"lofi247",
	)

}

func TestMatchUsesAndAcrossDimensions(t *testing.T) {

	results := Match(
		Criteria{
			Genres: []string{
				"Ambient",
			},
			Contexts: []string{
				"sleep",
			},
			MaxEnergy: 1,
		},
	)

	assertContainsStation(
		t,
		results,
		"somafm-dronezone",
	)

	assertNotContainsStation(
		t,
		results,
		"amambient",
	)

}

func TestMatchEmptyCriteriaReturnsAllStations(t *testing.T) {

	results := Match(
		Criteria{},
	)

	if len(results) != len(Stations) {
		t.Fatalf(
			"expected %d stations, got %d",
			len(Stations),
			len(results),
		)
	}

}

func TestMatchNoResults(t *testing.T) {

	results := Match(
		Criteria{
			Genres: []string{
				"Definitely Not A Genre",
			},
		},
	)

	if len(results) != 0 {
		t.Fatalf(
			"expected no matches, got %d",
			len(results),
		)
	}

}

func assertContainsStation(
	t *testing.T,
	stations []Station,
	id string,
) {

	t.Helper()

	for _, station := range stations {

		if station.ID == id {
			return
		}

	}

	t.Errorf(
		"expected station %q in results",
		id,
	)

}

func assertNotContainsStation(
	t *testing.T,
	stations []Station,
	id string,
) {

	t.Helper()

	for _, station := range stations {

		if station.ID == id {

			t.Errorf(
				"did not expect station %q in results",
				id,
			)

			return

		}

	}

}
